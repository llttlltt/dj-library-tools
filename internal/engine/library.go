package engine

import (
	"github.com/llttlltt/dj-library-tools/pkg/rekordbox"
)

// Library defines the interface for a music library source.
// This allows the Engine to operate on any source (Rekordbox XML, SQL, Mock)
// without being coupled to a specific implementation.
type Library interface {
	// GetTracks returns all tracks in the library.
	GetTracks() []rekordbox.Track
	// GetPlaylists returns the root nodes of the playlist tree.
	GetPlaylists() []rekordbox.Node
}

// RekordboxLibrary is an adapter that makes RekordboxLibraryXML satisfy the Library interface.
type RekordboxLibrary struct {
	XML *rekordbox.RekordboxLibraryXML
}

func (r *RekordboxLibrary) GetTracks() []rekordbox.Track {
	return r.XML.Collection.TRACK
}

func (r *RekordboxLibrary) GetPlaylists() []rekordbox.Node {
	return r.XML.Playlists.Node.Node
}

// NewRekordboxLibrary creates a new Library wrapper for a Rekordbox XML.
func NewRekordboxLibrary(xml *rekordbox.RekordboxLibraryXML) *RekordboxLibrary {
	return &RekordboxLibrary{XML: xml}
}
