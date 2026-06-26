package provider

import (
	"github.com/llttlltt/dj-library-tools/internal/models"
)

// Provider defines the interface for a music library provider.
type Provider interface {
	Name() string
	GetTracks(query string) ([]models.Track, error)
	GetPlaylists(query string) ([]models.Node, error)
	GetFolders(query string) ([]models.Node, error)
	// CanTranscode reports whether this provider can supply raw audio for transcoding.
	CanTranscode() bool
}

// WritableProvider extends Provider with modification capabilities.
type WritableProvider interface {
	Provider
	AddTracks(target models.Node, tracks []models.Track) (int, error)
	RemoveTracks(target models.Node, tracks []models.Track) (int, error)
	CreateNode(parent models.Node, name string, nodeType int) (models.Node, error)
	DeleteNode(node models.Node) error
	RenameNode(node models.Node, newName string) error
	MoveNode(node models.Node, targetParent models.Node) error
	// Save persists any in-memory mutations to the given path.
	Save(path string) error
}
