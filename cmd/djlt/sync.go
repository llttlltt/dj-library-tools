package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/llttlltt/dj-library-tools/internal/media"
	"github.com/llttlltt/dj-library-tools/internal/plex"
	"github.com/llttlltt/dj-library-tools/internal/sync"
	"github.com/llttlltt/dj-library-tools/pkg/rekordbox"
	"github.com/spf13/cobra"
)

var (
	exportDest   string
	exportFormat string
)

var syncCmd = &cobra.Command{
	Use:   "sync [source] [target] [query]",
	Short: "Sync items between a source and target",
	Args:  cobra.MinimumNArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		source := args[0]
		target := args[1]
		queryStr := ""
		if len(args) > 2 {
			queryStr = args[2]
		}

		if source == "plex" && target == "rekordbox" {
			return syncPlexToRekordbox(queryStr)
		}

		return fmt.Errorf("unsupported sync direction: %s to %s", source, target)
	},
}

func syncPlexToRekordbox(queryStr string) error {
	token := os.Getenv("PLEX_TOKEN")
	if token == "" {
		return fmt.Errorf("PLEX_TOKEN environment variable not set")
	}

	if xmlPath == "" {
		return fmt.Errorf("rekordbox XML path not specified (use --xml or -x)")
	}

	rbXML, err := rekordbox.ReadRekordboxLibrary(xmlPath)
	if err != nil {
		return fmt.Errorf("failed to read rekordbox library: %w", err)
	}

	client := plex.NewClient(token)
	resources, err := client.GetResources()
	if err != nil {
		return fmt.Errorf("failed to get resources: %w", err)
	}

	var targetServer plex.Resource
	found := false
	for _, res := range resources {
		if res.Provides == "server" && len(res.Connections) > 0 {
			targetServer = res
			found = true
			break
		}
	}

	if !found {
		return fmt.Errorf("no plex servers found")
	}

	baseURL := targetServer.Connections[0].URI
	serverClient := plex.NewClient(targetServer.AccessToken)

	// Fetch playlists to find the ID
	playlists, err := serverClient.GetPlaylists(baseURL)
	if err != nil {
		return fmt.Errorf("failed to get playlists: %w", err)
	}

	var targetPlaylist plex.Playlist
	for _, pl := range playlists {
		// Matching query as a playlist name for now
		if pl.Title == queryStr || pl.RatingKey == queryStr {
			targetPlaylist = pl
			break
		}
	}

	if targetPlaylist.RatingKey == "" {
		return fmt.Errorf("plex playlist matching '%s' not found", queryStr)
	}

	fmt.Printf("Syncing Plex playlist: %s (%d tracks)...\n", targetPlaylist.Title, targetPlaylist.LeafCount)

	tracks, err := serverClient.GetPlaylistTracks(baseURL, targetPlaylist.Key)
	if err != nil {
		return fmt.Errorf("failed to get tracks: %w", err)
	}

	// Setup Media Engine if export flag is set
	var transcoder *media.Transcoder
	if exportDest != "" {
		cfg := media.DefaultConfig()
		cfg.Dest = exportDest
		if exportFormat != "" {
			cfg.Format = exportFormat
		}
		transcoder = media.NewTranscoder(cfg)
		fmt.Printf("Exporting files to: %s (format: %s)\n", exportDest, cfg.Format)
	}

	matcher := sync.NewMatcher(rbXML.Collection.TRACK)
	var rbTrackIDs []string

	for _, track := range tracks {
		match := matcher.Match(track)
		var rbTrack *rekordbox.Track

		if match.RBTrack != nil && match.Confidence >= 0.8 {
			rbTrack = match.RBTrack
			fmt.Printf("  Matched: %s - %s (Confidence: %.2f)\n", track.Artist, track.Title, match.Confidence)
		} else if transcoder != nil {
			// If not matched but we are exporting, we could potentially "create" it in the XML later.
			// For now, let's focus on the file export part.
			fmt.Printf("  Exporting (No Match): %s - %s\n", track.Artist, track.Title)
		}

		if transcoder != nil {
			destPath, err := transcoder.GetDestinationPath(media.PathMetadata{
				Artist: track.Artist,
				Album:  track.Album,
				Title:  track.Title,
			})
			if err != nil {
				fmt.Printf("    Path error: %v\n", err)
				continue
			}

			// Ensure directory exists
			if err := os.MkdirAll(filepath.Dir(destPath), 0755); err != nil {
				fmt.Printf("    Dir error: %v\n", err)
				continue
			}

			// Transcode (for now we assume track.Media[0].Part[0].File is reachable if running locally or we need a path mapping)
			// This is a risk: Plex paths might not be local.
			sourceFile := track.Media[0].Part[0].File
			if _, err := os.Stat(sourceFile); err != nil {
				fmt.Printf("    Source file not found: %s\n", sourceFile)
				continue
			}

			if err := transcoder.Transcode(sourceFile, destPath); err != nil {
				fmt.Printf("    Transcode error: %v\n", err)
				continue
			}
			
			// If we matched a track, we should update its location in the XML to the exported path
			if rbTrack != nil {
				// Convert to file:// URL or just absolute path depending on RB requirement
				rbTrack.Location = "file://localhost" + destPath
			}
		}

		if rbTrack != nil {
			rbTrackIDs = append(rbTrackIDs, fmt.Sprintf("%d", rbTrack.TrackID))
		}
	}

	// Inject into XML
	if err := injectPlaylist(rbXML, targetPlaylist.Title, rbTrackIDs); err != nil {
		return fmt.Errorf("failed to inject playlist: %w", err)
	}

	if err := rekordbox.WriteRekordboxLibrary(xmlPath, rbXML); err != nil {
		return fmt.Errorf("failed to save XML: %w", err)
	}

	fmt.Printf("Sync complete.\n")
	return nil
}

func injectPlaylist(lib *rekordbox.RekordboxLibraryXML, name string, trackIDs []string) error {
	var syncFolder *rekordbox.Node
	for i := range lib.Playlists.Node.Node {
		if lib.Playlists.Node.Node[i].Name == "Plex Sync" && lib.Playlists.Node.Node[i].Type == 0 {
			syncFolder = &lib.Playlists.Node.Node[i]
			break
		}
	}

	if syncFolder == nil {
		lib.Playlists.Node.Node = append(lib.Playlists.Node.Node, rekordbox.Node{
			BaseNode: rekordbox.BaseNode{
				Type: 0,
				Name: "Plex Sync",
			},
		})
		syncFolder = &lib.Playlists.Node.Node[len(lib.Playlists.Node.Node)-1]
	}

	newPlaylist := rekordbox.Node{
		BaseNode: rekordbox.BaseNode{
			Type: 1,
			Name: name,
		},
		KeyType: 0,
		Entries: int32(len(trackIDs)),
	}

	for _, id := range trackIDs {
		newPlaylist.TRACK = append(newPlaylist.TRACK, struct {
			Key string `xml:"Key,attr"`
		}{Key: id})
	}

	found := false
	for i := range syncFolder.Node {
		if syncFolder.Node[i].Name == name && syncFolder.Node[i].Type == 1 {
			syncFolder.Node[i] = newPlaylist
			found = true
			break
		}
	}

	if !found {
		syncFolder.Node = append(syncFolder.Node, newPlaylist)
		syncFolder.Count++
	}

	return nil
}

func init() {
	syncCmd.Flags().StringVar(&exportDest, "dest", "", "Destination directory for exported files")
	syncCmd.Flags().StringVar(&exportFormat, "format", "mp3", "Target format for exported files")
	rootCmd.AddCommand(syncCmd)
}
