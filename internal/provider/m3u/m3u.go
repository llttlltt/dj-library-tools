package m3u

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/llttlltt/dj-library-tools/internal/models"
	"github.com/llttlltt/dj-library-tools/internal/m3u"
	"github.com/llttlltt/dj-library-tools/internal/provider"
	"github.com/llttlltt/dj-library-tools/internal/provider/factory"
	"github.com/llttlltt/dj-library-tools/internal/query"
)

func init() {
	factory.Register("m3u", func(opts factory.ProviderOptions) (provider.Provider, error) {
		return NewM3UProvider(opts.FilePath)
	})
	factory.Register("m3u8", func(opts factory.ProviderOptions) (provider.Provider, error) {
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
			if err := p.load(); err != nil {
				return nil, err
			}
		}
	}
	return p, nil
}

func (p *M3UProvider) Name() string { return "m3u" }

func (p *M3UProvider) Tracks() provider.TrackService { return &m3uTrackService{p} }
func (p *M3UProvider) Groups() provider.GroupService { return &m3uGroupService{p} }
func (p *M3UProvider) System() provider.SystemService { return &m3uSystemService{p} }

type m3uTrackService struct{ *M3UProvider }

func (s *m3uTrackService) List(ctx provider.ExecutionContext, queryString string) ([]models.Track, error) {
	parser := query.NewParser()
	q := parser.Parse(queryString)
	eval := query.NewEvaluator(q)

	var results []models.Track
	for _, t := range s.tracks {
		if eval.Matches(t) {
			results = append(results, t)
		}
	}
	return results, nil
}

func (s *m3uTrackService) Update(ctx provider.ExecutionContext, query string, changes map[string]string) (int, error) {
	return 0, fmt.Errorf("m3u provider does not support metadata modification")
}

func (s *m3uTrackService) UpdateBatch(ctx provider.ExecutionContext, matches []models.MetadataMatch, fields []string) error {
	return fmt.Errorf("m3u does not support batch metadata updates")
}

func (s *m3uTrackService) Delete(ctx provider.ExecutionContext, query string) (int, error) {
	return 0, fmt.Errorf("m3u provider does not support track deletion")
}

func (s *m3uTrackService) Groups() provider.TrackGroupService { return s }

func (s *m3uTrackService) Add(ctx provider.ExecutionContext, tracks []models.Track, target models.ResourceGroup) (int, error) {
	added := 0
	existing := make(map[string]bool)
	for _, t := range s.tracks { existing[t.Location] = true }
	for _, t := range tracks {
		if !existing[t.Location] {
			s.tracks = append(s.tracks, t)
			existing[t.Location] = true
			added++
		}
	}
	return added, nil
}

func (s *m3uTrackService) Remove(ctx provider.ExecutionContext, tracks []models.Track, group models.ResourceGroup) (int, error) {
	toRemove := make(map[string]bool)
	for _, t := range tracks { toRemove[t.Location] = true }
	var kept []models.Track
	removed := 0
	for _, t := range s.tracks {
		if toRemove[t.Location] { removed++ } else { kept = append(kept, t) }
	}
	s.tracks = kept
	return removed, nil
}

func (s *m3uTrackService) Move(ctx provider.ExecutionContext, tracks []models.Track, from models.ResourceGroup, to models.ResourceGroup) (int, error) {
	return 0, fmt.Errorf("m3u provider does not support move")
}

func (s *m3uTrackService) Sort(ctx provider.ExecutionContext, tracks []models.Track, field string) {}

type m3uGroupService struct{ *M3UProvider }

func (s *m3uGroupService) List(ctx provider.ExecutionContext, queryString string) ([]models.ResourceGroup, error) {
	name := filepath.Base(s.path)
	n := models.ResourceGroup{ID: s.path, Name: name, Type: models.GroupTypePlaylist, Items: len(s.tracks)}
	parser := query.NewParser()
	q := parser.Parse(queryString)
	eval := query.NewEvaluator(q)
	if eval.MatchesGroup(n) { return []models.ResourceGroup{n}, nil }
	return nil, nil
}

func (s *m3uGroupService) Create(ctx provider.ExecutionContext, parent models.ResourceGroup, name string, gt models.GroupType, pos int) (models.ResourceGroup, error) {
	if gt == models.GroupTypeFolder { return models.ResourceGroup{}, fmt.Errorf("m3u does not support folders") }
	return models.ResourceGroup{Name: name, Type: models.GroupTypePlaylist}, nil
}

func (s *m3uGroupService) Update(ctx provider.ExecutionContext, group models.ResourceGroup, newName string, newParent *models.ResourceGroup) error {
	if newName != "" {
		newPath := filepath.Join(filepath.Dir(s.path), newName)
		if err := os.Rename(s.path, newPath); err != nil { return err }
		s.path = newPath
	}
	return nil
}

func (s *m3uGroupService) Delete(ctx provider.ExecutionContext, group models.ResourceGroup) error {
	return os.Remove(s.path)
}

func (s *m3uGroupService) Sort(ctx provider.ExecutionContext, groups []models.ResourceGroup, field string) {}

type m3uSystemService struct{ *M3UProvider }

func (s *m3uSystemService) Capabilities() provider.ProviderCapabilities {
	return provider.ProviderCapabilities{CanWrite: true, IsFileBased: true}
}
func (s *m3uSystemService) Containment() provider.ContainmentPolicy { return provider.ContainmentPolicy{} }
func (s *m3uSystemService) MetadataCapabilities() []string { return []string{"display", "location"} }
func (s *m3uSystemService) SupportedResources() []string { return []string{"tracks", "playlists"} }
func (s *m3uSystemService) Save(ctx provider.ExecutionContext, path string) error {
	if path == "" { path = s.path }
	if path == "" { return fmt.Errorf("no path for M3U save") }
	f, err := os.Create(path)
	if err != nil { return err }
	defer f.Close()
	m3u.WriteM3U8Header(f)
	for _, t := range s.tracks {
		disp := t.Display
		if disp == "" { disp = filepath.Base(t.Location) }
		m3u.WriteM3U8EntryRaw(f, disp, t.Location, float64(t.Duration))
	}
	return nil
}
func (s *m3uSystemService) Fix(ctx provider.ExecutionContext, resource string, query string) error { return s.Save(ctx, "") }
func (s *m3uSystemService) Sync(ctx provider.ExecutionContext, tracks []models.Track, srcQ, tgtQ string, opts provider.SyncOptions) error {
	if opts.AppendOnly { s.Tracks().Groups().Add(ctx, tracks, models.ResourceGroup{}) } else { s.tracks = tracks }
	return s.Save(ctx, "")
}
func (s *m3uSystemService) Identify(name string, gt models.GroupType) string { return s.path }

func (p *M3UProvider) load() error {
	f, err := os.Open(p.path)
	if err != nil { return err }
	defer f.Close()
	scanner := bufio.NewScanner(f)
	var lastDur int
	var lastDisp string
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#EXTM3U") { continue }
		if strings.HasPrefix(line, "#EXTINF:") {
			info := strings.TrimPrefix(line, "#EXTINF:")
			commaIdx := strings.Index(info, ",")
			if commaIdx != -1 {
				if d, err := strconv.Atoi(info[:commaIdx]); err == nil { lastDur = d }
				lastDisp = strings.TrimSpace(info[commaIdx+1:])
			}
			continue
		}
		trackPath := line
		if !filepath.IsAbs(trackPath) { trackPath = filepath.Join(filepath.Dir(p.path), trackPath) }
		p.tracks = append(p.tracks, models.Track{ID: trackPath, Display: lastDisp, Duration: lastDur, Location: trackPath})
		lastDur = 0; lastDisp = ""
	}
	return scanner.Err()
}
