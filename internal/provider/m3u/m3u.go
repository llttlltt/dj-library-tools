package m3u

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/llttlltt/dj-library-tools/internal/models"
	"github.com/llttlltt/dj-library-tools/internal/playlist"
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

func (p *M3UProvider) Name() string {
	return "m3u"
}

func (p *M3UProvider) load() error {
	f, err := os.Open(p.path)
	if err != nil {
		return err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	var currentMeta playlist.AudioMetadata
	var tracks []models.Track

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#EXTM3U") {
			continue
		}

		if strings.HasPrefix(line, "#EXTINF:") {
			info := strings.TrimPrefix(line, "#EXTINF:")
			commaIdx := strings.Index(info, ",")
			if commaIdx != -1 {
				metaStr := info[commaIdx+1:]
				if strings.Contains(metaStr, " - ") {
					parts := strings.SplitN(metaStr, " - ", 2)
					currentMeta.Artist = strings.TrimSpace(parts[0])
					currentMeta.Title = strings.TrimSpace(parts[1])
				} else {
					currentMeta.Title = strings.TrimSpace(metaStr)
				}
			}
			continue
		}

		trackPath := line
		if !filepath.IsAbs(trackPath) {
			trackPath = filepath.Join(filepath.Dir(p.path), trackPath)
		}

		title := currentMeta.Title
		if title == "" {
			title = filepath.Base(trackPath)
		}

		tracks = append(tracks, models.Track{
			ID:       trackPath,
			Title:    title,
			Artist:   currentMeta.Artist,
			Location: trackPath,
		})
		currentMeta = playlist.AudioMetadata{}
	}

	p.tracks = tracks
	return scanner.Err()
}

func (p *M3UProvider) Capabilities() provider.ProviderCapabilities {
	return provider.ProviderCapabilities{
		CanWrite:          true,
		CanManageGroups:   false,
		CanUpdateMetadata: false,
		SupportsCues:      false,
		SupportsBeatgrids: false,
		IsFileBased:       true,
	}
}

func (p *M3UProvider) GetContainmentPolicy() provider.ContainmentPolicy {
	return provider.ContainmentPolicy{
		AllowTracksInFolders:    false,
		AllowFoldersInPlaylists: false,
		AllowNestedFolders:      false,
	}
}

func (p *M3UProvider) CustomMatch(track models.Track, field string, op query.Operator, value string) bool {
	return false
}

func (p *M3UProvider) GetTracks(ctx provider.ExecutionContext, queryString string) ([]models.Track, error) {
	parser := query.NewParser()
	q := parser.Parse(queryString)
	eval := query.NewEvaluator(q)

	var results []models.Track
	for _, t := range p.tracks {
		if eval.Matches(t) {
			results = append(results, t)
		}
	}
	return results, nil
}

func (p *M3UProvider) GetPlaylists(ctx provider.ExecutionContext, queryString string) ([]models.ResourceGroup, error) {
	name := filepath.Base(p.path)
	n := models.ResourceGroup{
		ID:    p.path,
		Name:  name,
		Type:  models.GroupTypePlaylist,
		Items: len(p.tracks),
	}

	parser := query.NewParser()
	q := parser.Parse(queryString)
	eval := query.NewEvaluator(q)
	if eval.MatchesGroup(n) {
		return []models.ResourceGroup{n}, nil
	}
	return nil, nil
}

func (p *M3UProvider) GetFolders(_ provider.ExecutionContext, _ string) ([]models.ResourceGroup, error) {
	return nil, nil
}

func (p *M3UProvider) GetResources(ctx provider.ExecutionContext, resource string, query string) ([]models.Resource, error) {
	var items []models.Resource
	switch resource {
	case "tracks":
		tracks, err := p.GetTracks(ctx, query)
		if err != nil { return nil, err }
		for _, t := range tracks { items = append(items, t) }
	case "playlists":
		groups, err := p.GetPlaylists(ctx, query)
		if err != nil { return nil, err }
		for _, g := range groups { items = append(items, g) }
	default:
		return nil, provider.ErrUnsupportedResource
	}
	return items, nil
}

func (p *M3UProvider) SortTracks(_ provider.ExecutionContext, tracks []models.Track, field string) {
}

func (p *M3UProvider) SortGroups(_ provider.ExecutionContext, groups []models.ResourceGroup, field string) {
}

func (p *M3UProvider) CanTranscode() bool {
	return true
}

func (p *M3UProvider) AddTracks(ctx provider.ExecutionContext, target models.ResourceGroup, tracks []models.Track) (int, error) {
	added := 0
	existing := make(map[string]bool)
	for _, t := range p.tracks {
		existing[t.Location] = true
	}

	for _, t := range tracks {
		if !existing[t.Location] {
			p.tracks = append(p.tracks, t)
			existing[t.Location] = true
			added++
		}
	}
	return added, nil
}

func (p *M3UProvider) RemoveTracks(ctx provider.ExecutionContext, target models.ResourceGroup, tracks []models.Track) (int, error) {
	toRemove := make(map[string]bool)
	for _, t := range tracks {
		toRemove[t.Location] = true
	}

	var kept []models.Track
	removed := 0
	for _, t := range p.tracks {
		if toRemove[t.Location] {
			removed++
		} else {
			kept = append(kept, t)
		}
	}
	p.tracks = kept
	return removed, nil
}

func (p *M3UProvider) CreateGroup(ctx provider.ExecutionContext, parent models.ResourceGroup, name string, groupType models.GroupType, position int) (models.ResourceGroup, error) {
	if groupType == models.GroupTypeFolder {
		return models.ResourceGroup{}, provider.ErrUnsupportedResource
	}
	return models.ResourceGroup{Name: name, Type: models.GroupTypePlaylist}, nil
}

func (p *M3UProvider) DeleteGroup(ctx provider.ExecutionContext, node models.ResourceGroup) error {
	return os.Remove(p.path)
}

func (p *M3UProvider) RenameGroup(ctx provider.ExecutionContext, node models.ResourceGroup, newName string, groupType models.GroupType) error {
	newPath := filepath.Join(filepath.Dir(p.path), newName)
	if err := os.Rename(p.path, newPath); err != nil {
		return err
	}
	p.path = newPath
	return nil
}

func (p *M3UProvider) MoveGroup(ctx provider.ExecutionContext, node models.ResourceGroup, targetParent models.ResourceGroup) error {
	return fmt.Errorf("m3u provider does not support move")
}

func (p *M3UProvider) Save(ctx provider.ExecutionContext, path string) error {
	if path == "playlists" || path == "tracks" {
		path = ""
	}
	if path == "" {
		path = p.path
	}
	if path == "" {
		return fmt.Errorf("no path specified for M3U save")
	}

	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}

	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	if err := playlist.WriteM3U8Header(f); err != nil {
		return err
	}

	for _, t := range p.tracks {
		meta := playlist.AudioMetadata{
			Artist: t.Artist,
			Title:  t.Title,
			Album:  t.Album,
		}
		if err := playlist.WriteM3U8Entry(f, meta, t.Location, float64(t.Duration)); err != nil {
			return err
		}
	}

	return nil
}

func (p *M3UProvider) Sync(ctx provider.ExecutionContext, tracks []models.Track, sourceQuery string, targetQuery string, options provider.SyncOptions) error {
	if options.AppendOnly {
		_, err := p.AddTracks(ctx, models.ResourceGroup{}, tracks)
		if err != nil {
			return err
		}
	} else {
		p.tracks = tracks
	}
	return p.Save(ctx, "")
}

func (p *M3UProvider) ModifyTracks(_ provider.ExecutionContext, _ string, _ map[string]string) (int, error) {
	return 0, fmt.Errorf("m3u provider does not support metadata modification")
}

func (p *M3UProvider) ValidateAddTracks(target models.ResourceGroup) error {
	return nil
}

func (p *M3UProvider) ValidateMoveGroup(src models.ResourceGroup, target models.ResourceGroup) error {
	return fmt.Errorf("m3u provider does not support move")
}

func (p *M3UProvider) ValidateCreateGroup(parent models.ResourceGroup, groupType models.GroupType) error {
	if groupType == models.GroupTypeFolder {
		return fmt.Errorf("m3u provider does not support folders")
	}
	return nil
}

func (p *M3UProvider) IdentifyGroup(_ string, _ models.GroupType) string {
	return p.path
}

func (p *M3UProvider) SupportedResources() []string {
	return []string{"tracks", "playlists"}
}

func (p *M3UProvider) MetadataCapabilities() []string {
	return []string{"title", "artist", "duration"}
}

func (p *M3UProvider) UpdateMetadata(_ provider.ExecutionContext, _ []models.MetadataMatch, _ []string) error {
	return fmt.Errorf("m3u does not support deep metadata updates")
}

func (p *M3UProvider) Fix(ctx provider.ExecutionContext, resource string, query string) error {
	// Transfer logic from fix.go
	return nil
}
