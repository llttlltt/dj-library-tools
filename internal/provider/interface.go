package provider

import (
	"github.com/llttlltt/dj-library-tools/internal/models"
)

// Provider defines the interface for a music library provider.
type Provider interface {
	Name() string
	GetTracks(query string) ([]models.Track, error)
	GetPlaylists(query string) ([]models.ResourceGroup, error)
	GetFolders(query string) ([]models.ResourceGroup, error)
	// CanTranscode reports whether this provider can supply raw audio for transcoding.
	CanTranscode() bool
}

// WritableProvider extends Provider with modification capabilities.
type WritableProvider interface {
	Provider
	AddTracks(target models.ResourceGroup, tracks []models.Track) (int, error)
	RemoveTracks(target models.ResourceGroup, tracks []models.Track) (int, error)
	CreateNode(parent models.ResourceGroup, name string, nodeType int) (models.ResourceGroup, error)
	DeleteNode(node models.ResourceGroup) error
	RenameNode(node models.ResourceGroup, newName string) error
	MoveNode(node models.ResourceGroup, targetParent models.ResourceGroup) error
	// Save persists any in-memory mutations to the given path.
	Save(path string) error
}
