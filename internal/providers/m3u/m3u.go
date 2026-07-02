package m3u

import (
	"context"
	"fmt"
	"github.com/llttlltt/dj-library-tools/internal/core/models"
	"github.com/llttlltt/dj-library-tools/internal/providers"
	"github.com/llttlltt/dj-library-tools/internal/providers/factory"
	"github.com/llttlltt/dj-library-tools/internal/services/library"
	"os"
	"path/filepath"
	"strings"
)

func init() {
	resources := []factory.ResourceInfo{
		{Name: "tracks", CanWrite: true, SupportsQuery: true},
		{Name: "playlists", CanWrite: true, SupportsQuery: true},
	}
	caps := provider.ProviderCapabilities{CanWrite: true, IsFileBased: true}
	factory.Register("m3u", resources, caps, func(opts factory.ProviderOptions) (provider.Provider, error) {
		return NewM3UProvider(opts.FilePath)
	})
	factory.Register("m3u8", resources, caps, func(opts factory.ProviderOptions) (provider.Provider, error) {
		return factory.NewProvider("m3u", opts)
	})
}

type M3UProvider struct {
	path   string
	tracks []models.Track
}

func NewM3UProvider(path string) (*M3UProvider, error) {
	p := &M3UProvider{path: path}
	if path != "" {
		if _, err := os.Stat(path); err == nil {
			tracks, err := ReadM3U8(path)
			if err != nil {
				return nil, err
			}
			p.tracks = tracks
		}
	}
	return p, nil
}

func (p *M3UProvider) Name() string { return "m3u" }

func (p *M3UProvider) Tracks() provider.TrackService  { return &m3uTrackService{p} }
func (p *M3UProvider) Groups() provider.GroupService  { return &m3uGroupService{p} }
func (p *M3UProvider) System() provider.SystemService { return &m3uSystemService{p} }

func ToNeutralTrack(t models.Track) models.Track {
	// If Title is missing, derive it from the location/filename
	if t.Title == "" && t.Location != "" {
		t.Title = filepath.Base(t.Location)
	}
	// If Display was set (from M3U EXTINF), use it for Title/Artist if they are blank
	if t.Display != "" {
		if t.Title == "" || t.Title == filepath.Base(t.Location) {
			// Basic heuristic: Artist - Title
			if parts := strings.SplitN(t.Display, " - ", 2); len(parts) == 2 {
				t.Artist = strings.TrimSpace(parts[0])
				t.Title = strings.TrimSpace(parts[1])
			} else {
				t.Title = t.Display
			}
		}
	}
	return t
}

type m3uTrackService struct{ *M3UProvider }

func (s *m3uTrackService) List(ctx context.Context, ectx provider.ExecutionContext, query string) ([]models.Track, error) {
	var tracks []models.Track
	for _, t := range s.tracks {
		tracks = append(tracks, ToNeutralTrack(t))
	}
	// Use Engine for agnostic querying
	eng := library.NewEngine(NewLibrary(tracks, filepath.Base(s.path)))
	return eng.Ls(query, nil)
}

func (s *m3uTrackService) Update(ctx context.Context, ectx provider.ExecutionContext, query string, changes map[string]string) (int, error) {
	return 0, fmt.Errorf("m3u provider does not support metadata modification")
}

func (s *m3uTrackService) UpdateBatch(ctx context.Context, ectx provider.ExecutionContext, matches []models.MetadataMatch, fields []string) error {
	return fmt.Errorf("m3u does not support batch metadata updates")
}

func (s *m3uTrackService) Delete(ctx context.Context, ectx provider.ExecutionContext, query string) (int, error) {
	return 0, fmt.Errorf("m3u provider does not support track deletion")
}

func (s *m3uTrackService) Groups() provider.TrackGroupService { return s }

func (s *m3uTrackService) Add(ctx context.Context, ectx provider.ExecutionContext, tracks []models.Track, target models.ResourceGroup) (int, error) {
	added := 0
	existing := make(map[string]bool)
	for _, t := range s.tracks {
		existing[t.Location] = true
	}
	for _, t := range tracks {
		if !existing[t.Location] {
			s.tracks = append(s.tracks, t)
			existing[t.Location] = true
			added++
		}
	}
	return added, nil
}

func (s *m3uTrackService) Remove(ctx context.Context, ectx provider.ExecutionContext, tracks []models.Track, group models.ResourceGroup) (int, error) {
	toRemove := make(map[string]bool)
	for _, t := range tracks {
		toRemove[t.Location] = true
	}
	var kept []models.Track
	removed := 0
	for _, t := range s.tracks {
		if toRemove[t.Location] {
			removed++
		} else {
			kept = append(kept, t)
		}
	}
	s.tracks = kept
	return removed, nil
}

func (s *m3uTrackService) Move(ctx context.Context, ectx provider.ExecutionContext, tracks []models.Track, from models.ResourceGroup, to models.ResourceGroup) (int, error) {
	return 0, fmt.Errorf("m3u provider does not support move")
}

func (s *m3uTrackService) Sort(ctx context.Context, ectx provider.ExecutionContext, tracks []models.Track, field string) {
}

type m3uGroupService struct{ *M3UProvider }

func (s *m3uGroupService) List(ctx context.Context, ectx provider.ExecutionContext, query string) ([]models.ResourceGroup, error) {
	// Use Engine for agnostic group querying
	eng := library.NewEngine(NewLibrary(s.tracks, filepath.Base(s.path)))
	return eng.LsGroups(query)
}

func (s *m3uGroupService) Create(ctx context.Context, ectx provider.ExecutionContext, parent models.ResourceGroup, name string, gt models.GroupKind, pos int) (models.ResourceGroup, error) {
	if gt == models.GroupKindFolder {
		return models.ResourceGroup{}, fmt.Errorf("m3u does not support folders")
	}
	return models.ResourceGroup{Name: name, Kind: models.GroupKindPlaylist}, nil
}

func (s *m3uGroupService) Update(ctx context.Context, ectx provider.ExecutionContext, group models.ResourceGroup, newName string, newParent *models.ResourceGroup) error {
	if newName != "" {
		newPath := filepath.Join(filepath.Dir(s.path), newName)
		if err := os.Rename(s.path, newPath); err != nil {
			return err
		}
		s.path = newPath
	}
	return nil
}

func (s *m3uGroupService) Delete(ctx context.Context, ectx provider.ExecutionContext, group models.ResourceGroup) error {
	return os.Remove(s.path)
}

func (s *m3uGroupService) Sort(ctx context.Context, ectx provider.ExecutionContext, groups []models.ResourceGroup, field string) {
}

type m3uSystemService struct{ *M3UProvider }

func (s *m3uSystemService) Capabilities() provider.ProviderCapabilities {
	return provider.ProviderCapabilities{CanWrite: true, IsFileBased: true}
}
func (s *m3uSystemService) Containment() provider.ContainmentPolicy {
	return provider.ContainmentPolicy{}
}
func (s *m3uSystemService) MetadataCapabilities() []string {
	return provider.ResolveAvailableFields(s.Capabilities())
}

func (s *m3uSystemService) SupportedResources() []string { return []string{"tracks", "playlists"} }
func (s *m3uSystemService) TableHeaders() []string {
	return []string{"duration", "display", "location"}
}

func (s *m3uSystemService) Save(ctx context.Context, ectx provider.ExecutionContext, path string) error {
	if path == "" {
		path = s.path
	}
	if path == "" {
		return fmt.Errorf("no path for M3U save")
	}

	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	WriteM3U8Header(f)

	isM3U8 := strings.ToLower(filepath.Ext(path)) == ".m3u8"

	for _, t := range s.tracks {
		if isM3U8 {
			disp := t.Display
			if disp == "" {
				disp = filepath.Base(t.Location)
			}
			WriteM3U8EntryRaw(f, disp, t.Location, float64(t.Duration))
		} else {
			WriteM3UEntryBasic(f, t.Location)
		}
	}
	return nil
}

func (s *m3uSystemService) Fix(ctx context.Context, ectx provider.ExecutionContext, selection provider.Selection, options provider.FixOptions) (int, error) {
	totalAffected := 0

	for fixType, targets := range options.Actions {
		switch fixType {
		case provider.FixPaths:
			for _, target := range targets {
				if target == "normalize" {
					res, err := FixPlaylist(s.path, FixOptions{
						M3U8:    strings.HasSuffix(s.path, ".m3u8"),
						Verbose: ectx.Verbose,
					})
					if err != nil {
						return totalAffected, err
					}
					totalAffected += res.TotalTracks
				}
			}
		}
	}

	return totalAffected, nil
}

func (s *m3uSystemService) Sync(ctx context.Context, ectx provider.ExecutionContext, tracks []models.Track, targetQuery string, opts provider.SyncOptions) error {
	if !ectx.Apply {
		return nil
	}
	if opts.AppendOnly {
		// Add unique tracks
		existing := make(map[string]bool)
		for _, t := range s.tracks {
			existing[t.Location] = true
		}
		for _, t := range tracks {
			if !existing[t.Location] {
				s.tracks = append(s.tracks, t)
				existing[t.Location] = true
			}
		}
	} else {
		s.tracks = tracks
	}
	return nil
}

func (s *m3uSystemService) Identify(name string, gt models.GroupKind) string { return s.path }
