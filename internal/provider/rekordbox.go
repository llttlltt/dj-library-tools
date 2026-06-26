package provider

import (
	"github.com/llttlltt/dj-library-tools/internal/engine"
	"github.com/llttlltt/dj-library-tools/pkg/rekordbox"
)

type RekordboxProvider struct {
	Engine *engine.Engine
}

func (p *RekordboxProvider) Name() string {
	return "rb"
}

func (p *RekordboxProvider) GetTracks(query string) ([]rekordbox.Track, error) {
	return p.Engine.Ls(query)
}

func (p *RekordboxProvider) GetPlaylists(query string) ([]NodeResult, error) {
	results, err := p.Engine.LsPlaylists(query)
	if err != nil {
		return nil, err
	}
	
	var out []NodeResult
	for _, r := range results {
		out = append(out, NodeResult{
			Name:         r.Node.Name,
			Entries:      int(rekordbox.DerefInt32(r.Node.Entries)),
			ParentFolder: r.ParentFolder,
			Raw:          r.Node,
		})
	}
	return out, nil
}

func (p *RekordboxProvider) GetRawTracks(query string) (interface{}, error) {
	return p.GetTracks(query)
}

func NewRekordboxProvider(eng *engine.Engine) *RekordboxProvider {
	return &RekordboxProvider{Engine: eng}
}
