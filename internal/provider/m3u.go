package provider

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/llttlltt/dj-library-tools/internal/models"
	"github.com/llttlltt/dj-library-tools/internal/playlist"
	"github.com/llttlltt/dj-library-tools/internal/query"
)

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
			// Try to parse metadata if available
			// #EXTINF:duration,Artist - Title
			info := strings.TrimPrefix(line, "#EXTINF:")
			commaIdx := strings.Index(info, ",")
			if commaIdx != -1 {
				// We don't strictly need duration yet, but could parse it
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

		// It's a path
		trackPath := line
		if !filepath.IsAbs(trackPath) {
			trackPath = filepath.Join(filepath.Dir(p.path), trackPath)
		}

		title := currentMeta.Title
		if title == "" {
			title = filepath.Base(trackPath)
		}

		tracks = append(tracks, models.Track{
			ID:       trackPath, // Use path as ID for M3U
			Title:    title,
			Artist:   currentMeta.Artist,
			Location: trackPath,
		})
		currentMeta = playlist.AudioMetadata{}
	}

	p.tracks = tracks
	return scanner.Err()
}

func (p *M3UProvider) GetTracks(queryString string) ([]models.Track, error) {
	q := query.NewParser().Parse(queryString)
	eval := query.NewEvaluator(q)

	var results []models.Track
	for _, t := range p.tracks {
		if eval.Matches(t) {
			results = append(results, t)
		}
	}
	return results, nil
}

func (p *M3UProvider) GetPlaylists(queryString string) ([]models.Node, error) {
	// An M3U file is itself a single playlist
	name := filepath.Base(p.path)
	n := models.Node{
		ID:    p.path,
		Name:  name,
		Type:  1,
		Items: len(p.tracks),
	}

	q := query.NewParser().Parse(queryString)
	eval := query.NewEvaluator(q)
	if eval.MatchesNode(n) {
		return []models.Node{n}, nil
	}
	return nil, nil
}

func (p *M3UProvider) GetFolders(_ string) ([]models.Node, error) {
	return nil, nil
}

func (p *M3UProvider) CanTranscode() bool {
	return true
}

func (p *M3UProvider) AddTracks(target models.Node, tracks []models.Track) (int, error) {
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

func (p *M3UProvider) RemoveTracks(target models.Node, tracks []models.Track) (int, error) {
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

func (p *M3UProvider) CreateNode(parent models.Node, name string, nodeType int) (models.Node, error) {
	if nodeType == 0 {
		return models.Node{}, fmt.Errorf("m3u provider does not support folders")
	}
	// For M3U, "creating a node" just means setting the path if it wasn't already.
	// But usually the path is provided in the location.
	return models.Node{Name: name, Type: 1}, nil
}

func (p *M3UProvider) DeleteNode(node models.Node) error {
	return os.Remove(p.path)
}

func (p *M3UProvider) RenameNode(node models.Node, newName string) error {
	newPath := filepath.Join(filepath.Dir(p.path), newName)
	if err := os.Rename(p.path, newPath); err != nil {
		return err
	}
	p.path = newPath
	return nil
}

func (p *M3UProvider) MoveNode(node models.Node, targetParent models.Node) error {
	return fmt.Errorf("m3u provider does not support move")
}

func (p *M3UProvider) Save(path string) error {
	// If path is "playlists" or "tracks", it's likely a CLI mask, ignore it
	if path == "playlists" || path == "tracks" {
		path = ""
	}
	if path == "" {
		path = p.path
	}
	if path == "" {
		return fmt.Errorf("no path specified for M3U save")
	}

	// Ensure the directory exists
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
		// Try to preserve relative paths if they were loaded that way?
		// For now, let's use absolute or whatever is in .Location
		if err := playlist.WriteM3U8Entry(f, meta, t.Location, float64(t.Duration)); err != nil {
			return err
		}
	}

	return nil
}
