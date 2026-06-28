package library

import (
	"github.com/llttlltt/dj-library-tools/internal/models"
)

// Library defines the interface for a music library source.
type Library interface {
	// GetTracks returns all tracks in the library in a neutral format.
	GetTracks() []models.Track
	// GetPlaylists returns the nodes of the playlist tree in a neutral format.
	GetPlaylists() []models.ResourceGroup
	// GetMembershipMap returns a mapping of track IDs to the names of playlists they belong to.
	GetMembershipMap() map[string][]string
}

// WritableLibrary extends Library with operations to modify the playlist tree.
type WritableLibrary interface {
	Library
	// AddPlaylist adds a new playlist at the given position in a folder.
	AddPlaylist(folder, name string, trackIDs []string, position int)
	// UpdatePlaylist replaces an existing playlist's tracks.
	UpdatePlaylist(name string, trackIDs []string) bool
	// AddTracksToPlaylist appends tracks to an existing playlist.
	AddTracksToPlaylist(name string, trackIDs []string) (bool, int)
	// RemoveTracksFromPlaylist removes tracks from an existing playlist.
	RemoveTracksFromPlaylist(name string, trackIDs []string) (bool, int)
	// CreateFolder creates a new folder.
	CreateFolder(folder, name string, position int) bool
	// RenameNode renames a folder or playlist.
	RenameNode(name, newName string, nodeType int32) bool
	// MoveNode moves a folder or playlist.
	MoveNode(name string, nodeType int32, targetFolder string) bool
	// RemoveNode removes a folder or playlist.
	RemoveNode(name string, nodeType int32) bool
	// Save writes changes back to persistent storage.
	Save(path string) error
}
