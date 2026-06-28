package library

import (
	"github.com/llttlltt/dj-library-tools/internal/models"
)

// ReadableLibrary defines the read-only interface for a music library source.
type ReadableLibrary interface {
	// GetResources returns all resources of a specific kind (track, group).
	GetResources(kind string) []models.Resource
	// GetMembershipMap returns a mapping of track IDs to the names of groups they belong to.
	GetMembershipMap() map[string][]string
}

// WritableLibrary extends ReadableLibrary with operations to modify the library.
type WritableLibrary interface {
	ReadableLibrary
	// CreateGroup creates a new group (Playlist or Folder) under the specified parent.
	CreateGroup(parentID, name string, groupType models.GroupType, position int) (models.ResourceGroup, error)
	// DeleteGroup removes a group from the library.
	DeleteGroup(groupID string, groupType models.GroupType) error
	// LinkTracks adds track memberships to a group.
	LinkTracks(groupID string, trackIDs []string) (int, error)
	// UnlinkTracks removes track memberships from a group.
	UnlinkTracks(groupID string, trackIDs []string) (int, error)
	// UpdateGroup replaces all track memberships in a group.
	UpdateGroup(groupID string, trackIDs []string) error
	// RenameGroup renames a group.
	RenameGroup(groupID, newName string, groupType models.GroupType) error
	// MoveGroup detaches a group and re-attaches it under a new parent.
	MoveGroup(groupID string, groupType models.GroupType, targetParentID string) error
	// Save writes changes back to persistent storage.
	Save(path string) error
}
