package sync

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/llttlltt/dj-library-tools/internal/library"
	"github.com/llttlltt/dj-library-tools/internal/media"
	"github.com/llttlltt/dj-library-tools/internal/models"
	"github.com/vbauerster/mpb/v8"
	"github.com/vbauerster/mpb/v8/decor"
)

type Orchestrator struct {
	Library    library.WritableLibrary
	SyncEngine *Engine
	DryRun     bool
	Verbose    bool
}

func NewOrchestrator(lib library.WritableLibrary, dryRun, verbose bool) *Orchestrator {
	return &Orchestrator{
		Library:    lib,
		SyncEngine: NewEngine(lib),
		DryRun:     dryRun,
		Verbose:    verbose,
	}
}

type SyncOptions struct {
	ExportDest   string
	ExportFormat string
	PathMaps     map[string]string
}

// SyncToLibrary matches a slice of neutral tracks against the collection,
// optionally transcodes them, then injects or appends to the named playlist.
func (o *Orchestrator) SyncToLibrary(tracks []models.Track, query string, playlistName string, opts SyncOptions, appendOnly bool) error {
	var transcoder *media.Transcoder
	if opts.ExportDest != "" {
		cfgMedia := media.DefaultConfig()
		cfgMedia.Dest = opts.ExportDest
		cfgMedia.PathMaps = opts.PathMaps
		if opts.ExportFormat != "" {
			cfgMedia.Format = opts.ExportFormat
		}
		transcoder = media.NewTranscoder(cfgMedia)
	}

	type transcodeJob struct {
		track models.Track
		rb    *models.Track
	}
	jobs := make(chan transcodeJob, len(tracks))
	results := make(chan string, len(tracks))
	errors := make(chan error, len(tracks))

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

	numWorkers := 4
	for w := 0; w < numWorkers; w++ {
		go func() {
			for job := range jobs {
				track := job.track
				rbTrack := job.rb

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

				if track.Location == "" {
					trackBar.Abort(false)
					errors <- fmt.Errorf("no media file for: %s - %s", track.Artist, track.Title)
					results <- ""
					totalBar.Increment()
					continue
				}

				if transcoder == nil {
					trackBar.Increment()
					if rbTrack != nil {
						results <- rbTrack.ID
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

				sourceFile := track.Location
				if _, err := os.Stat(transcoder.ApplyPathMap(sourceFile)); err != nil {
					trackBar.Abort(false)
					errors <- fmt.Errorf("source not found for %s: %s", track.Title, sourceFile)
					results <- ""
					totalBar.Increment()
					continue
				}

				if o.DryRun {
					trackBar.Increment()
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
				}

				if rbTrack != nil {
					results <- rbTrack.ID
				} else {
					results <- ""
				}
				totalBar.Increment()
			}
		}()
	}

	for _, track := range tracks {
		match := o.SyncEngine.Matcher.Match(track)
		var rbTrack *models.Track
		if match.RBTrack != nil && match.Confidence >= 0.8 {
			rbTrack = match.RBTrack
		}
		jobs <- transcodeJob{track: track, rb: rbTrack}
	}
	close(jobs)
	p.Wait()

	var trackIDs []string
	for i := 0; i < len(tracks); i++ {
		if res := <-results; res != "" {
			trackIDs = append(trackIDs, res)
		}
	}
	close(errors)
	for err := range errors {
		fmt.Printf("  Error: %v\n", err)
	}

	if o.DryRun {
		action := "inject"
		if appendOnly {
			action = "append to"
		}
		fmt.Printf("[Dry Run] Would %s playlist %q with %d tracks into XML\n", action, playlistName, len(trackIDs))
	} else {
		if appendOnly {
			o.SyncEngine.AddTracksToPlaylist(playlistName, trackIDs)
			fmt.Printf("Appended %d tracks to %q\n", len(trackIDs), playlistName)
		} else {
			result := o.SyncEngine.InjectPlaylist(playlistName, trackIDs)
			fmt.Printf("Synced playlist %q (%d tracks).\n", result.PlaylistName, result.TracksInjected)
		}
	}

	return nil
}

// DefaultSyncFolder is the top-level folder name used for sync operations.
const DefaultSyncFolder = "Synced"

// Engine manages sync operations against a music library.
type Engine struct {
	Library library.WritableLibrary
	Matcher *Matcher
}

// NewEngine creates a sync Engine backed by the given library.
func NewEngine(lib library.WritableLibrary) *Engine {
	return &Engine{
		Library: lib,
		Matcher: NewMatcher(lib.GetTracks()),
	}
}

// SyncResult holds the outcome of a single playlist injection.
type SyncResult struct {
	PlaylistName   string
	TracksInjected int
	// Updated is true when an existing playlist was replaced; false when newly created.
	Updated bool
}

// UpsertPlaylist creates or replaces a named playlist inside folder.
// When folder is empty the playlist is placed at the root level.
// position is the 0-indexed position in the folder. -1 appends to the end.
func (e *Engine) UpsertPlaylist(folder, name string, trackIDs []string, position int) *SyncResult {
	updated := e.Library.UpdatePlaylist(name, trackIDs)
	if !updated {
		e.Library.AddPlaylist(folder, name, trackIDs, position)
	}

	return &SyncResult{
		PlaylistName:   name,
		TracksInjected: len(trackIDs),
		Updated:        updated,
	}
}

// InjectPlaylist upserts a named playlist under DefaultSyncFolder.
func (e *Engine) InjectPlaylist(name string, trackIDs []string) *SyncResult {
	return e.UpsertPlaylist(DefaultSyncFolder, name, trackIDs, -1)
}

// AddTracksToPlaylist appends trackIDs to a named playlist anywhere in the tree.
// Returns (true, addedCount) if the playlist was found, (false, 0) otherwise.
func (e *Engine) AddTracksToPlaylist(name string, trackIDs []string) (bool, int) {
	return e.Library.AddTracksToPlaylist(name, trackIDs)
}

// RemoveTracksFromPlaylist removes all trackIDs present in the given slice from a named playlist.
// Returns (true, removedCount) if the playlist was found, (false, 0) otherwise.
func (e *Engine) RemoveTracksFromPlaylist(name string, trackIDs []string) (bool, int) {
	return e.Library.RemoveTracksFromPlaylist(name, trackIDs)
}

// CreateFolder creates a new folder node at the specified position.
func (e *Engine) CreateFolder(folder, name string, position int) bool {
	return e.Library.CreateFolder(folder, name, position)
}

// RenameGroup renames the first node matching name and nodeType anywhere in the tree.
func (e *Engine) RenameGroup(name, newName string, nodeType int32) bool {
	return e.Library.RenameGroup(name, newName, nodeType)
}

// MoveGroup detaches the first node matching name and nodeType from its current location.
func (e *Engine) MoveGroup(name string, nodeType int32, targetFolder string) bool {
	return e.Library.MoveGroup(name, nodeType, targetFolder)
}

// RemoveGroup removes the first node matching name and nodeType from anywhere in the tree.
func (e *Engine) RemoveGroup(name string, nodeType int32) bool {
	return e.Library.RemoveGroup(name, nodeType)
}

// MatchTracks matches a slice of neutral tracks against the collection.
func (e *Engine) MatchTracks(tracks []models.Track, minConfidence float64) []MatchResult {
	out := make([]MatchResult, 0, len(tracks))
	for _, t := range tracks {
		m := e.Matcher.Match(t)
		if m.RBTrack != nil && m.Confidence >= minConfidence {
			out = append(out, m)
		}
	}
	return out
}

// Save writes the modified library back to disk.
func (e *Engine) Save(path string) error {
	return e.Library.Save(path)
}
