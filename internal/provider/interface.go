package provider

import (
	"github.com/llttlltt/dj-library-tools/pkg/rekordbox"
)

// Provider defines the interface for a music library provider (Rekordbox, Plex, etc.)
type Provider interface {
	// Name returns the provider's name (e.g. "rb", "plex")
	Name() string
	// GetTracks resolves tracks matching the query.
	GetTracks(query string) ([]rekordbox.Track, error)
	// GetPlaylists resolves playlists matching the query.
	GetPlaylists(query string) ([]NodeResult, error)
}

// NodeResult is a matched playlist or folder node along with its parent info.
type NodeResult struct {
	Name         string
	Entries      int
	ParentFolder string
	Raw          interface{} // Provider-specific node data
}
