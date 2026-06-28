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
	// AddGroup creates a new playlist at the given position in a folder.
	AddGroup(folder, name string, trackIDs []string, position int)
	// UpdateGroup replaces an existing playlist's tracks.
	UpdateGroup(name string, trackIDs []string) bool
	// AddTracksToGroup appends tracks to an existing playlist.
	AddTracksToGroup(name string, trackIDs []string) (bool, int)
	// RemoveTracksFromGroup removes tracks from an existing playlist.
	RemoveTracksFromGroup(name string, trackIDs []string) (bool, int)
	// CreateContainer creates a new folder.
	CreateContainer(folder, name string, position int) bool
	// RenameGroup renames a folder or playlist.
	RenameGroup(name, newName string, nodeType int32) bool
	// MoveGroup moves a folder or playlist.
	MoveGroup(name string, nodeType int32, targetFolder string) bool
	// RemoveGroup removes a folder or playlist.
	RemoveGroup(name string, nodeType int32) bool
	// Save writes changes back to persistent storage.
	Save(path string) error
}
