package provider

import (
	"github.com/llttlltt/dj-library-tools/internal/djerr"
	"github.com/llttlltt/dj-library-tools/internal/models"
)

var (
	ErrReadOnly            = djerr.ErrReadOnly
	ErrInvalidParent       = djerr.ErrInvalidParent
	ErrUnsupportedResource = djerr.ErrUnsupportedResource
)

// ProviderCapabilities defines what a provider is able to do.
type ProviderCapabilities struct {
	CanWrite          bool
	CanManageGroups   bool // Create/Move/Rename Folders and Playlists
	CanUpdateMetadata bool // Update track properties (bpm, comment, etc.)
	SupportsCues      bool // Custom matching for hotcues/memorycues
	SupportsBeatgrids bool // Custom matching for beatgrids
	IsFileBased       bool // Requires --file flag
}

// ContainmentPolicy defines the structural rules of the library.
type ContainmentPolicy struct {
	AllowTracksInFolders   bool
	AllowFoldersInPlaylists bool
	AllowNestedFolders      bool
}

// ExecutionContext holds runtime state for provider operations.
type ExecutionContext struct {
	DryRun  bool
	Verbose bool
}

// Provider is the entry point for all provider-specific operations.
type Provider interface {
	Name() string
	Tracks() TrackService
	Groups() GroupService
	System() SystemService
}

// TrackService handles operations on individual music files and their metadata.
type TrackService interface {
	// List returns tracks matching a query.
	List(ctx ExecutionContext, query string) ([]models.Track, error)

	// Update modifies metadata for tracks matching a query.
	Update(ctx ExecutionContext, query string, changes map[string]string) (int, error)

	// UpdateBatch applies specific field updates to matched tracks (used for syncing).
	UpdateBatch(ctx ExecutionContext, matches []models.MetadataMatch, fields []string) error

	// Delete removes tracks from the provider's collection.
	Delete(ctx ExecutionContext, query string) (int, error)

	// Groups returns a service for managing track-to-group relationships.
	Groups() TrackGroupService

	// Sort orders a slice of tracks in-place.
	Sort(ctx ExecutionContext, tracks []models.Track, field string)
}

// TrackGroupService handles track memberships within groups (playlists).
type TrackGroupService interface {
	// Add adds tracks to a specific group.
	Add(ctx ExecutionContext, tracks []models.Track, target models.ResourceGroup) (int, error)

	// Remove removes tracks from a specific group.
	Remove(ctx ExecutionContext, tracks []models.Track, group models.ResourceGroup) (int, error)

	// Move transfers tracks from one group to another.
	Move(ctx ExecutionContext, tracks []models.Track, from models.ResourceGroup, to models.ResourceGroup) (int, error)
}

// GroupService handles structural items (Playlists, Folders) themselves.
type GroupService interface {
	// List returns groups matching a query.
	List(ctx ExecutionContext, query string) ([]models.ResourceGroup, error)

	// Create creates a new group container.
	Create(ctx ExecutionContext, parent models.ResourceGroup, name string, groupType models.GroupType, position int) (models.ResourceGroup, error)

	// Update modifies group properties (rename, move group in tree).
	Update(ctx ExecutionContext, group models.ResourceGroup, newName string, newParent *models.ResourceGroup) error

	// Delete removes a group container.
	Delete(ctx ExecutionContext, group models.ResourceGroup) error

	// Sort orders a slice of groups in-place.
	Sort(ctx ExecutionContext, groups []models.ResourceGroup, field string)
}

// SystemService handles provider-wide configuration, maintenance, and orchestration.
type SystemService interface {
	Capabilities() ProviderCapabilities
	Containment() ContainmentPolicy
	MetadataCapabilities() []string
	SupportedResources() []string

	// Save writes changes to persistent storage.
	Save(ctx ExecutionContext, path string) error

	// Fix performs health/formatting repairs.
	Fix(ctx ExecutionContext, resource string, query string) error

	// Sync orchestrates a full library sync.
	Sync(ctx ExecutionContext, tracks []models.Track, targetQuery string, options SyncOptions) error

	// Identify returns a provider-specific ID for a name.
	Identify(name string, groupType models.GroupType) string
}

type SyncOptions struct {
	ExportDest   string
	ExportFormat string
	PathMaps     map[string]string
	AppendOnly   bool
}

func ToTrackSlice(res []models.Resource) []models.Track {
	var tracks []models.Track
	for _, r := range res {
		if t, ok := r.(models.Track); ok {
			tracks = append(tracks, t)
		}
	}
	return tracks
}
