package sync

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/llttlltt/dj-library-tools/internal/library"
	"github.com/llttlltt/dj-library-tools/internal/media"
	"github.com/llttlltt/dj-library-tools/internal/models"
)

type ProgressListener interface {
	OnStart(total int64)
	OnTrackStart(trackTitle string)
	OnTrackEnd()
	OnComplete()
}

type Orchestrator struct {
	Library    library.WritableLibrary
	DryRun     bool
	Verbose    bool
	Listener   ProgressListener
	Matcher    *Matcher
}

func NewOrchestrator(lib library.WritableLibrary, dryRun, verbose bool) *Orchestrator {
	var tracks []models.Track
	resources := lib.GetResources("track")
	for _, r := range resources {
		tracks = append(tracks, r.(models.Track))
	}
	
	return &Orchestrator{
		Library:    lib,
		DryRun:     dryRun,
		Verbose:    verbose,
		Matcher:    NewMatcher(tracks),
	}
}

type SyncOptions struct {
	ExportDest   string
	ExportFormat string
	PathMaps     map[string]string
}

// Join matches source tracks against the target library using the specified keys.
func (o *Orchestrator) Join(sourceTracks []models.Track, matchFields []string) []models.MetadataMatch {
	var matches []models.MetadataMatch
	
	// Create a new matcher specifically for these fields if needed, 
	// for now we use the default (Artist/Title).
	for _, st := range sourceTracks {
		match := o.Matcher.Match(st)
		if match.TargetTrack != nil && match.Confidence >= 0.8 {
			matches = append(matches, models.MetadataMatch{
				Source: st,
				Target: *match.TargetTrack,
			})
		}
	}
	
	return matches
}

func (o *Orchestrator) SyncToLibrary(tracks []models.Track, sourceQuery string, playlistName string, opts SyncOptions, appendOnly bool) error {
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

	if o.Listener != nil {
		o.Listener.OnStart(int64(len(tracks)))
	}

	type transcodeJob struct {
		track models.Track
		target *models.Track
	}
	jobs := make(chan transcodeJob, len(tracks))
	results := make(chan string, len(tracks))
	errors := make(chan error, len(tracks))

	numWorkers := 4
	for w := 0; w < numWorkers; w++ {
		go func() {
			for job := range jobs {
				track := job.track
				targetTrack := job.target

				if o.Listener != nil {
					o.Listener.OnTrackStart(track.Title)
				}

				if track.Location == "" {
					errors <- fmt.Errorf("no media file for: %s - %s", track.Artist, track.Title)
					results <- ""
					if o.Listener != nil { o.Listener.OnTrackEnd() }
					continue
				}

				if transcoder == nil {
					if targetTrack != nil {
						results <- targetTrack.ID
					} else {
						results <- ""
					}
					if o.Listener != nil { o.Listener.OnTrackEnd() }
					continue
				}

				destPath, err := transcoder.GetDestinationPath(media.PathMetadata{
					Artist: track.Artist,
					Album:  track.Album,
					Title:  track.Title,
				})
				if err != nil {
					errors <- fmt.Errorf("path error for %s: %v", track.Title, err)
					results <- ""
					if o.Listener != nil { o.Listener.OnTrackEnd() }
					continue
				}

				sourceFile := track.Location
				if _, err := os.Stat(transcoder.ApplyPathMap(sourceFile)); err != nil {
					errors <- fmt.Errorf("source not found for %s: %s", track.Title, sourceFile)
					results <- ""
					if o.Listener != nil { o.Listener.OnTrackEnd() }
					continue
				}

				if !o.DryRun {
					if err := os.MkdirAll(filepath.Dir(destPath), 0755); err != nil {
						errors <- fmt.Errorf("mkdir error for %s: %v", track.Title, err)
						results <- ""
						if o.Listener != nil { o.Listener.OnTrackEnd() }
						continue
					}
					if err := transcoder.Transcode(sourceFile, destPath); err != nil {
						errors <- fmt.Errorf("transcode error for %s: %v", track.Title, err)
						results <- ""
						if o.Listener != nil { o.Listener.OnTrackEnd() }
						continue
					}
				}

				if targetTrack != nil {
					results <- targetTrack.ID
				} else {
					results <- ""
				}
				if o.Listener != nil { o.Listener.OnTrackEnd() }
			}
		}()
	}

	for _, track := range tracks {
		match := o.Matcher.Match(track)
		var targetTrack *models.Track
		if match.TargetTrack != nil && match.Confidence >= 0.8 {
			targetTrack = match.TargetTrack
		}
		jobs <- transcodeJob{track: track, target: targetTrack}
	}
	close(jobs)
	
	var trackIDs []string
	for i := 0; i < len(tracks); i++ {
		if res := <-results; res != "" {
			trackIDs = append(trackIDs, res)
		}
	}
	close(errors)
	for err := range errors {
		if o.Verbose { fmt.Printf("  Error: %v\n", err) }
	}

	if !o.DryRun {
		if appendOnly {
			o.Library.LinkTracks(playlistName, trackIDs)
		} else {
			err := o.Library.UpdateGroup(playlistName, trackIDs)
			if err != nil {
				o.Library.CreateGroup("", playlistName, models.GroupTypePlaylist, -1)
				o.Library.UpdateGroup(playlistName, trackIDs)
			}
		}
	}

	if o.Listener != nil {
		o.Listener.OnComplete()
	}

	return nil
}
