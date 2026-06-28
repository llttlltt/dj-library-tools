package sync

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

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
	Apply     bool
	Verbose    bool
	Listener   ProgressListener
	Matcher    *Matcher
}

func NewOrchestrator(lib library.WritableLibrary, apply, verbose bool) *Orchestrator {
	var tracks []models.Track
	if lib != nil {
		resources := lib.GetResources("track")
		for _, r := range resources {
			tracks = append(tracks, r.(models.Track))
		}
	}
	
	return &Orchestrator{
		Library:    lib,
		Apply:     apply,
		Verbose:    verbose,
		Matcher:    NewMatcher(tracks),
	}
}

func (o *Orchestrator) getExistingTracks() []models.Track {
	resources := o.Library.GetResources("track")
	var tracks []models.Track
	for _, r := range resources {
		tracks = append(tracks, r.(models.Track))
	}
	return tracks
}

type SyncOptions struct {
	ExportDest     string
	ExportFormat   string
	PathMaps       map[string]string
	MetadataFields []string
	MatchFields    []string
}

// SyncToLibrary is a high-level helper that coordinates a full sync from source tracks to a target library.
func SyncToLibrary(lib library.WritableLibrary, tracks []models.Track, targetQuery string, options SyncOptions, apply, verbose bool, appendOnly bool) error {
	orch := NewOrchestrator(lib, apply, verbose)
	
	// Perform metadata reconciliation if requested
	if len(options.MetadataFields) > 0 {
		orch.Matcher = NewMatcher(orch.getExistingTracks()).WithKeys(options.MatchFields)
		matches := orch.Join(tracks, options.MatchFields)
		
		if err := orch.Library.(interface {
			UpdateMetadata(matches []models.MetadataMatch, fields []string) error
		}).UpdateMetadata(matches, options.MetadataFields); err != nil {
			return err
		}
	}

	return orch.SyncToLibrary(tracks, targetQuery, options, appendOnly)
}

// Join matches source tracks against the target library using the specified keys.
func (o *Orchestrator) Join(sourceTracks []models.Track, matchFields []string) []models.MetadataMatch {
	var matches []models.MetadataMatch
	
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

// Relocate searches for physical files for the given tracks in the searchDir.
func (o *Orchestrator) Relocate(tracks []models.Track, searchDir string, matchFields []string) map[string]string {
	relocated := make(map[string]string)
	
	fileMap := make(map[string][]string)
	filepath.Walk(searchDir, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() { return nil }
		name := strings.ToLower(info.Name())
		fileMap[name] = append(fileMap[name], path)
		return nil
	})

	for _, t := range tracks {
		filename := strings.ToLower(filepath.Base(t.Location))
		candidates, ok := fileMap[filename]
		if !ok { continue }

		for _, candidate := range candidates {
			relocated[t.ID] = candidate
			if o.Verbose {
				fmt.Printf("Relocated %s -> %s\n", t.Title, candidate)
			}
			break
		}
	}

	return relocated
}

func (o *Orchestrator) SyncToLibrary(tracks []models.Track, targetID string, opts SyncOptions, appendOnly bool) error {
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

				if !o.Apply {
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

	if !o.Apply {
		if appendOnly {
			o.Library.AddTracks(targetID, trackIDs)
		} else {
			err := o.Library.UpdateGroup(targetID, trackIDs)
			if err != nil {
				o.Library.CreateGroup("", targetID, models.GroupKindPlaylist, -1)
				o.Library.UpdateGroup(targetID, trackIDs)
			}
		}
	}

	if o.Listener != nil {
		l := o.Listener
		l.OnComplete()
	}

	return nil
}
