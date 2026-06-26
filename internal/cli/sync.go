package cli

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/llttlltt/dj-library-tools/internal/config"
	"github.com/llttlltt/dj-library-tools/internal/media"
	"github.com/llttlltt/dj-library-tools/internal/plex"
	"github.com/llttlltt/dj-library-tools/internal/sync"
	"github.com/llttlltt/dj-library-tools/internal/utils"
	"github.com/llttlltt/dj-library-tools/internal/playlist"
	"github.com/llttlltt/dj-library-tools/pkg/rekordbox"
	"github.com/spf13/cobra"
	"github.com/vbauerster/mpb/v8"
	"github.com/vbauerster/mpb/v8/decor"
)

var (
	exportDest   string
	exportFormat string
	dryRun       bool
)

var syncCmd = &cobra.Command{
	Use:   "sync [source-location] [target-location]",
	Short: "Sync items between a source and target",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		src := utils.ParseLocation(args[0], "")
		tgt := utils.ParseLocation(args[1], "")

		if src.Provider == "plex" && tgt.Provider == "rb" {
			return syncPlexToRekordbox(src, tgt)
		}
		if src.Provider == "plex" && tgt.Provider == "m3u8" {
			return syncPlexToM3U8(src, tgt)
		}

		return fmt.Errorf("unsupported sync direction: %s to %s", src.Provider, tgt.Provider)
	},
}

func syncPlexToRekordbox(src, tgt utils.Location) error {
	cfg, _ := config.LoadAppConfig()
	token := os.Getenv("PLEX_TOKEN")
	if token == "" {
		token = cfg.PlexToken
	}
	if token == "" {
		return fmt.Errorf("Plex token not found. Run 'djlt auth plex' or set PLEX_TOKEN env var")
	}

	path := utils.ExpandPath(xmlPath)
	if path == "" {
		path = utils.ExpandPath(cfg.RekordboxXMLPath)
	}
	if path == "" {
		return fmt.Errorf("Rekordbox XML path not found. Use --xml or run 'djlt config rekordbox --xml PATH'")
	}

	rbXML, err := rekordbox.ReadRekordboxLibrary(path)
	if err != nil {
		return fmt.Errorf("failed to read rekordbox library: %w", err)
	}

	client := plex.NewClient(token)
	resources, err := client.GetResources(context.Background())
	if err != nil {
		return fmt.Errorf("failed to get resources: %w", err)
	}

	var targetServer plex.Resource
	found := false
	for _, res := range resources {
		if res.Provides != "server" {
			continue
		}
		targetServer = res
		found = true
		break
	}

	if !found {
		return fmt.Errorf("no plex servers found")
	}

	serverClient := plex.NewClient(targetServer.AccessToken)
	probe, err := serverClient.ProbeBestConnection(targetServer)
	if err != nil {
		return fmt.Errorf("failed to connect to server: %w", err)
	}

	fmt.Printf("Connected via: %s\n", probe.BaseURL)

	var targetPlaylist plex.Playlist
	for _, pl := range probe.Playlists {
		if pl.Title == src.Query || pl.RatingKey == src.Query {
			targetPlaylist = pl
			break
		}
	}

	if targetPlaylist.RatingKey == "" {
		return fmt.Errorf("plex playlist matching '%s' not found", src.Query)
	}

	fmt.Printf("Syncing Plex playlist: %s (%d tracks)...\n", targetPlaylist.Title, targetPlaylist.LeafCount)

	tracks, err := serverClient.GetPlaylistTracks(context.Background(), probe.BaseURL, targetPlaylist.Key)
	if err != nil {
		return fmt.Errorf("failed to get tracks: %w", err)
	}

	// Setup Media Engine if export flag is set
	var transcoder *media.Transcoder
	if exportDest != "" {
		cfgMedia := media.DefaultConfig()
		cfgMedia.Dest = exportDest
		cfgMedia.PathMaps = cfg.PathMaps
		if exportFormat != "" {
			cfgMedia.Format = exportFormat
		}
		transcoder = media.NewTranscoder(cfgMedia)
		fmt.Printf("Exporting files to: %s (format: %s)\n", exportDest, cfgMedia.Format)
	}

	syncEng := sync.NewEngine(serverClient, rbXML)
	var rbTrackIDs []string

	type transcodeJob struct {
		track plex.Track
		rb    *rekordbox.Track
	}
	jobs := make(chan transcodeJob, len(tracks))
	results := make(chan string, len(tracks))
	errors := make(chan error, len(tracks))

	// Progress bar setup
	p := mpb.New(mpb.WithWidth(64))
	totalBar := p.AddBar(int64(len(tracks)),
		mpb.PrependDecorators(
			decor.Name("Overall Sync", decor.WCSyncSpaceR),
			decor.CountersNoUnit("%d / %d", decor.WCSyncSpace),
		),
		mpb.AppendDecorators(
			decor.Percentage(decor.WCSyncSpace),
			decor.OnComplete(decor.EwmaETA(decor.ET_STYLE_GO, 60), "done!"),
		),
	)

	// Start workers
	numWorkers := 4 
	for w := 1; w <= numWorkers; w++ {
		go func() {
			for job := range jobs {
				track := job.track
				rbTrack := job.rb
				
				// Multi-bar support: create a temporary bar for this track
				// We limit track name length for display
				displayName := track.Title
				if len(displayName) > 20 {
					displayName = displayName[:17] + "..."
				}
				trackBar := p.AddBar(1, 
					mpb.BarRemoveOnComplete(),
					mpb.PrependDecorators(
						decor.Name(fmt.Sprintf("  -> %s", displayName), decor.WCSyncSpaceR),
					),
				)

				if len(track.Media) == 0 || len(track.Media[0].Part) == 0 {
					trackBar.Abort(false)
					errors <- fmt.Errorf("no media file for: %s - %s", track.Artist, track.Title)
					results <- ""
					totalBar.Increment()
					continue
				}

				if transcoder == nil {
					trackBar.Increment()
					if rbTrack != nil {
						results <- fmt.Sprintf("%d", rbTrack.TrackID)
					} else {
						results <- ""
					}
					totalBar.Increment()
					continue
				}

				destPath, err := transcoder.GetDestinationPath(media.PathMetadata{
					Artist: track.Artist,
					Album:  track.Album,
					Title:  track.Title,
				})
				if err != nil {
					trackBar.Abort(false)
					errors <- fmt.Errorf("path error for %s: %v", track.Title, err)
					results <- ""
					totalBar.Increment()
					continue
				}

				sourceFile := track.Media[0].Part[0].File
				if _, err := os.Stat(transcoder.ApplyPathMap(sourceFile)); err != nil {
					trackBar.Abort(false)
					errors <- fmt.Errorf("source not found for %s: %s", track.Title, sourceFile)
					results <- ""
					totalBar.Increment()
					continue
				}

				if dryRun {
					trackBar.Increment()
					if rbTrack != nil {
						rbTrack.Location = "file://localhost" + destPath
					}
				} else {
					if err := os.MkdirAll(filepath.Dir(destPath), 0755); err != nil {
						trackBar.Abort(false)
						errors <- fmt.Errorf("mkdir error for %s: %v", track.Title, err)
						results <- ""
						totalBar.Increment()
						continue
					}
					if err := transcoder.Transcode(sourceFile, destPath); err != nil {
						trackBar.Abort(false)
						errors <- fmt.Errorf("transcode error for %s: %v", track.Title, err)
						results <- ""
						totalBar.Increment()
						continue
					}
					trackBar.Increment()
					if rbTrack != nil {
						rbTrack.Location = "file://localhost" + destPath
					}
				}

				if rbTrack != nil {
					results <- fmt.Sprintf("%d", rbTrack.TrackID)
				} else {
					results <- ""
				}
				totalBar.Increment()
			}
		}()
	}

	// Feed jobs
	for _, track := range tracks {
		match := syncEng.Matcher.Match(track)
		var rbTrack *rekordbox.Track

		if match.RBTrack != nil && match.Confidence >= 0.8 {
			rbTrack = match.RBTrack
		}
		
		jobs <- transcodeJob{track: track, rb: rbTrack}
	}
	close(jobs)

	// Wait for progress bars to complete
	p.Wait()

	// Collect results
	for i := 0; i < len(tracks); i++ {
		res := <-results
		if res != "" {
			rbTrackIDs = append(rbTrackIDs, res)
		}
	}

	// Report errors
	close(errors)
	for err := range errors {
		fmt.Printf("  Error: %v\n", err)
	}

	if dryRun {
		fmt.Printf("[Dry Run] Would inject playlist '%s' with %d tracks into XML\n", targetPlaylist.Title, len(rbTrackIDs))
		fmt.Printf("[Dry Run] Would save updated library to: %s\n", path)
	} else {
		result := syncEng.InjectPlaylist(targetPlaylist.Title, rbTrackIDs)
		action := "Created"
		if result.Updated {
			action = "Updated"
		}
		fmt.Printf("%s playlist '%s' (%d tracks injected).\n", action, result.PlaylistName, result.TracksInjected)
		if err := syncEng.SaveXML(path); err != nil {
			return fmt.Errorf("failed to save XML: %w", err)
		}
	}

	fmt.Printf("Sync complete.\n")
	return nil
}

func syncPlexToM3U8(src, tgt utils.Location) error {
	cfg, _ := config.LoadAppConfig()
	token := os.Getenv("PLEX_TOKEN")
	if token == "" {
		token = cfg.PlexToken
	}
	if token == "" {
		return fmt.Errorf("Plex token not found. Run 'djlt auth plex' or set PLEX_TOKEN env var")
	}

	client := plex.NewClient(token)
	resources, err := client.GetResources(context.Background())
	if err != nil {
		return fmt.Errorf("failed to get resources: %w", err)
	}

	var targetServer plex.Resource
	found := false
	for _, res := range resources {
		if res.Provides != "server" {
			continue
		}
		targetServer = res
		found = true
		break
	}

	if !found {
		return fmt.Errorf("no plex servers found")
	}

	serverClient := plex.NewClient(targetServer.AccessToken)
	probe, err := serverClient.ProbeBestConnection(targetServer)
	if err != nil {
		return fmt.Errorf("failed to connect to server: %w", err)
	}

	var targetPlaylist plex.Playlist
	for _, pl := range probe.Playlists {
		if pl.Title == src.Query || pl.RatingKey == src.Query {
			targetPlaylist = pl
			break
		}
	}

	if targetPlaylist.RatingKey == "" {
		return fmt.Errorf("plex playlist matching '%s' not found", src.Query)
	}

	fmt.Printf("Syncing Plex playlist: %s (%d tracks)...\n", targetPlaylist.Title, targetPlaylist.LeafCount)

	tracks, err := serverClient.GetPlaylistTracks(context.Background(), probe.BaseURL, targetPlaylist.Key)
	if err != nil {
		return fmt.Errorf("failed to get tracks: %w", err)
	}

	m3uPath := utils.ExpandPath(tgt.Query)
	if m3uPath == "" {
		m3uPath = targetPlaylist.Title + ".m3u8"
	}
	if filepath.Ext(m3uPath) == "" {
		m3uPath += ".m3u8"
	}

	// Setup Media Engine if export flag is set
	var transcoder *media.Transcoder
	if exportDest != "" {
		cfgMedia := media.DefaultConfig()
		cfgMedia.Dest = exportDest
		cfgMedia.PathMaps = cfg.PathMaps
		if exportFormat != "" {
			cfgMedia.Format = exportFormat
		}
		transcoder = media.NewTranscoder(cfgMedia)
	}

	var m3uBody strings.Builder
	playlist.WriteM3U8Header(&m3uBody)

	for _, track := range tracks {
		trackPath := track.Media[0].Part[0].File
		
		if transcoder != nil {
			destPath, err := transcoder.GetDestinationPath(media.PathMetadata{
				Artist: track.Artist,
				Album:  track.Album,
				Title:  track.Title,
			})
			if err == nil {
				if dryRun {
					fmt.Printf("    [Dry Run] Would transcode: %s -> %s\n", trackPath, destPath)
					trackPath = destPath
				} else {
					os.MkdirAll(filepath.Dir(destPath), 0755)
					if err := transcoder.Transcode(trackPath, destPath); err == nil {
						trackPath = destPath
					}
				}
			}
		}

		// Make path relative to m3u8 if possible
		if rel, err := filepath.Rel(filepath.Dir(m3uPath), trackPath); err == nil {
			trackPath = rel
		}

		playlist.WriteM3U8Entry(&m3uBody, playlist.AudioMetadata{
			Artist: track.Artist,
			Title:  track.Title,
		}, trackPath, 0)
	}

	if dryRun {
		fmt.Printf("[Dry Run] Would create M3U8: %s with contents:\n%s\n", m3uPath, m3uBody.String())
	} else {
		f, err := os.Create(m3uPath)
		if err != nil {
			return fmt.Errorf("failed to create m3u8 file: %w", err)
		}
		defer f.Close()
		f.WriteString(m3uBody.String())
	}

	fmt.Printf("Sync complete.\n")
	return nil
}



func init() {
	syncCmd.Flags().StringVar(&exportDest, "dest", "", "Destination directory for exported files")
	syncCmd.Flags().StringVar(&exportFormat, "format", "mp3", "Target format for exported files")
	syncCmd.Flags().BoolVar(&dryRun, "dry-run", false, "Preview changes without writing files or XML")
	RootCmd.AddCommand(syncCmd)
}
