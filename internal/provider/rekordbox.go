package provider

import (
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
	return p.Engine.Ls(query)
}

func (p *RekordboxProvider) GetPlaylists(query string) ([]models.Node, error) {
	return p.Engine.LsPlaylists(query)
}

func (p *RekordboxProvider) GetRawTracks(query string) (interface{}, error) {
	matched, _ := p.Engine.Ls(query)
	var raw []rekordbox.Track
	for _, m := range matched {
		if rt, ok := m.Raw.(rekordbox.Track); ok {
			raw = append(raw, rt)
		}
	}
	return raw, nil
}

func (p *RekordboxProvider) CanTranscode() bool {
	return true
}

func (p *RekordboxProvider) AddTracks(target models.Node, tracks []models.Track) (int, error) {
	var ids []string
	for _, t := range tracks {
		ids = append(ids, t.ID)
	}
	_, added := p.Engine.Library.(engine.WritableLibrary).AddTracksToPlaylist(target.Name, ids)
	return added, nil
}

func (p *RekordboxProvider) RemoveTracks(target models.Node, tracks []models.Track) (int, error) {
	var ids []string
	for _, t := range tracks {
		ids = append(ids, t.ID)
	}
	_, removed := p.Engine.Library.(engine.WritableLibrary).RemoveTracksFromPlaylist(target.Name, ids)
	return removed, nil
}

func NewRekordboxProvider(eng *engine.Engine) *RekordboxProvider {
	return &RekordboxProvider{Engine: eng}
}
