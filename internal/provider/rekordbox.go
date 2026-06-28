package provider

import (
	"github.com/llttlltt/dj-library-tools/internal/library"
	"github.com/llttlltt/dj-library-tools/internal/models"
)

type RekordboxProvider struct {
	Engine *library.Engine
	path   string
}

func (p *RekordboxProvider) Name() string {
	return "rb"
}

func (p *RekordboxProvider) GetTracks(query string) ([]models.Track, error) {
	return p.Engine.Ls(query)
}

func (p *RekordboxProvider) GetPlaylists(query string) ([]models.ResourceGroup, error) {
	return p.Engine.LsPlaylists(query)
}

func (p *RekordboxProvider) GetFolders(query string) ([]models.ResourceGroup, error) {
	return p.Engine.LsFolders(query)
}

func (p *RekordboxProvider) CanTranscode() bool {
	return true
}

func (p *RekordboxProvider) AddTracks(target models.ResourceGroup, tracks []models.Track) (int, error) {
	var ids []string
	for _, t := range tracks {
		ids = append(ids, t.ID)
	}
	_, added := p.Engine.Library.(library.WritableLibrary).AddTracksToPlaylist(target.Name, ids)
	return added, nil
}

func (p *RekordboxProvider) RemoveTracks(target models.ResourceGroup, tracks []models.Track) (int, error) {
	var ids []string
	for _, t := range tracks {
		ids = append(ids, t.ID)
	}
	_, removed := p.Engine.Library.(library.WritableLibrary).RemoveTracksFromPlaylist(target.Name, ids)
	return removed, nil
}

func (p *RekordboxProvider) CreateNode(parent models.ResourceGroup, name string, nodeType int) (models.ResourceGroup, error) {
	if nodeType == 0 {
		p.Engine.Library.(library.WritableLibrary).CreateFolder(parent.Name, name, -1)
	} else {
		p.Engine.Library.(library.WritableLibrary).AddPlaylist(parent.Name, name, nil, -1)
	}
	return models.ResourceGroup{Name: name, Type: models.GroupType(nodeType)}, nil
}

func (p *RekordboxProvider) DeleteNode(node models.ResourceGroup) error {
	p.Engine.Library.(library.WritableLibrary).RemoveNode(node.Name, int32(node.Type))
	return nil
}

func (p *RekordboxProvider) RenameNode(node models.ResourceGroup, newName string) error {
	p.Engine.Library.(library.WritableLibrary).RenameNode(node.Name, newName, int32(node.Type))
	return nil
}

func (p *RekordboxProvider) MoveNode(node models.ResourceGroup, targetParent models.ResourceGroup) error {
	p.Engine.Library.(library.WritableLibrary).MoveNode(node.Name, int32(node.Type), targetParent.Name)
	return nil
}

func (p *RekordboxProvider) Save(path string) error {
	if path == "" {
		path = p.path
	}
	return p.Engine.Library.(library.WritableLibrary).Save(path)
}

func NewRekordboxProvider(eng *library.Engine, path string) *RekordboxProvider {
	return &RekordboxProvider{Engine: eng, path: path}
}
