package provider

import (
	"github.com/llttlltt/dj-library-tools/internal/models"
)

// Provider defines the interface for a music library provider (Rekordbox, Plex, etc.)
type Provider interface {
	// Name returns the provider's name (e.g. "rb", "plex")
	Name() string
	// GetTracks resolves tracks matching the query.
	GetTracks(query string) ([]models.Track, error)
	// GetPlaylists resolves playlists matching the query.
	GetPlaylists(query string) ([]models.Node, error)
	// GetRawTracks returns provider-specific track models (e.g. []plex.Track)
	GetRawTracks(query string) (interface{}, error)
}
