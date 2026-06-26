package provider

import (
	"strconv"

	"github.com/llttlltt/dj-library-tools/internal/engine"
	"github.com/llttlltt/dj-library-tools/internal/models"
	"github.com/llttlltt/dj-library-tools/pkg/rekordbox"
)

type RekordboxProvider struct {
	Engine *engine.Engine
}

func (p *RekordboxProvider) Name() string {
	return "rb"
}

func (p *RekordboxProvider) GetTracks(query string) ([]models.Track, error) {
	rbTracks, err := p.Engine.Ls(query)
	if err != nil {
		return nil, err
	}

	var tracks []models.Track
	for _, rt := range rbTracks {
		t := models.Track{
			ID:       strconv.Itoa(rt.TrackID),
			Title:    rt.Name,
			Artist:   rt.Artist,
			Album:    rt.Album,
			Key:      rt.Tonality,
			Location: rt.Location,
			Rating:   int(rt.Rating / 51),
			Raw:      rt,
		}
		if rt.AverageBpm != "" {
			t.BPM, _ = strconv.ParseFloat(rt.AverageBpm, 64)
		}
		tracks = append(tracks, t)
	}
	return tracks, nil
}

func (p *RekordboxProvider) GetPlaylists(query string) ([]models.Node, error) {
	results, err := p.Engine.LsPlaylists(query)
	if err != nil {
		return nil, err
	}

	var out []models.Node
	for _, r := range results {
		out = append(out, models.Node{
			Name:         r.Node.Name,
			Entries:      int(rekordbox.DerefInt32(r.Node.Entries)),
			ParentFolder: r.ParentFolder,
			Type:         1,
			Raw:          r.Node,
		})
	}
	return out, nil
}

func (p *RekordboxProvider) GetRawTracks(query string) (interface{}, error) {
	// Re-fetch since we need the specific rekordbox slice for some engines
	return p.Engine.Ls(query)
}

func NewRekordboxProvider(eng *engine.Engine) *RekordboxProvider {
	return &RekordboxProvider{Engine: eng}
}
