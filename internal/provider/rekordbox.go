package provider

import (
	"github.com/llttlltt/dj-library-tools/internal/engine"
	"github.com/llttlltt/dj-library-tools/internal/models"
)

type RekordboxProvider struct {
	Engine *engine.Engine
	path   string
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

func (p *RekordboxProvider) GetFolders(query string) ([]models.Node, error) {
	return p.Engine.LsFolders(query)
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

func (p *RekordboxProvider) CreateNode(parent models.Node, name string, nodeType int) (models.Node, error) {
	p.Engine.Library.(engine.WritableLibrary).CreateFolder(parent.Name, name, -1)
	return models.Node{Name: name, Type: nodeType}, nil
}

func (p *RekordboxProvider) DeleteNode(node models.Node) error {
	p.Engine.Library.(engine.WritableLibrary).RemoveNode(node.Name, int32(node.Type))
	return nil
}

func (p *RekordboxProvider) RenameNode(node models.Node, newName string) error {
	p.Engine.Library.(engine.WritableLibrary).RenameNode(node.Name, newName, int32(node.Type))
	return nil
}

func (p *RekordboxProvider) MoveNode(node models.Node, targetParent models.Node) error {
	p.Engine.Library.(engine.WritableLibrary).MoveNode(node.Name, int32(node.Type), targetParent.Name)
	return nil
}

func (p *RekordboxProvider) Save(path string) error {
	if path == "" {
		path = p.path
	}
	return p.Engine.Library.(engine.WritableLibrary).Save(path)
}

func NewRekordboxProvider(eng *engine.Engine, path string) *RekordboxProvider {
	return &RekordboxProvider{Engine: eng, path: path}
}
