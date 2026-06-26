package provider

import (
	"github.com/llttlltt/dj-library-tools/internal/models"
)

// Provider defines the interface for a music library provider.
type Provider interface {
	Name() string
	GetTracks(query string) ([]models.Track, error)
	GetPlaylists(query string) ([]models.Node, error)
	GetRawTracks(query string) (interface{}, error)
	
	// Capabilities
	CanTranscode() bool // Can this provider provide raw audio for transcoding?
}

// WritableProvider extends Provider with modification capabilities.
type WritableProvider interface {
	Provider
	AddTracks(target models.Node, tracks []models.Track) (int, error)
	RemoveTracks(target models.Node, tracks []models.Track) (int, error)
}
