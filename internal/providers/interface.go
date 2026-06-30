package provider

import (
	"context"

	djerrors "github.com/llttlltt/dj-library-tools/internal/core/errors"
	"github.com/llttlltt/dj-library-tools/internal/core/location"
	"github.com/llttlltt/dj-library-tools/internal/core/models"
)

var (
	ErrReadOnly            = djerrors.ErrReadOnly
	ErrInvalidParent       = djerrors.ErrInvalidParent
	ErrUnsupportedResource = djerrors.ErrUnsupportedResource
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

// ResolveAvailableFields returns a list of queryable fields based on provider capabilities.
func ResolveAvailableFields(caps ProviderCapabilities) []string {
	var fields []string
	// Map flags to model capability set
	enabled := make(map[models.Capability]bool)
	enabled[models.CapNone] = true
	if caps.CanUpdateMetadata { enabled[models.CapMetadata] = true }
	if caps.SupportsCues { enabled[models.CapCues] = true }
	if caps.SupportsBeatgrids { enabled[models.CapBeatgrids] = true }

	for name, def := range models.TrackFields {
		if enabled[def.RequiredCap] {
			fields = append(fields, name)
		}
	}
	return fields
}

// ContainmentPolicy defines the structural rules of the library.
type ContainmentPolicy struct {
	AllowTracksInFolders   bool
	AllowFoldersInPlaylists bool
	AllowNestedFolders      bool
}

// ExecutionContext holds runtime state and feedback channels for provider operations.
type ExecutionContext struct {
	Apply    bool
	Verbose  bool
	Feedback Feedback
}

// Feedback defines an interface for providing user feedback during operations.
type Feedback interface {
	OnPreview(message string)
	OnSuccess(message string)
	OnWarning(message string)
	OnStatus(message string)
	OnProgress(done, total int)
	OnTable(headers []string, rows [][]string)
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
	List(ctx context.Context, ectx ExecutionContext, query string) ([]models.Track, error)

	// Update modifies metadata for tracks matching a query.
	Update(ctx context.Context, ectx ExecutionContext, query string, changes map[string]string) (int, error)

	// UpdateBatch applies specific field updates to matched tracks (used for syncing).
	UpdateBatch(ctx context.Context, ectx ExecutionContext, matches []models.MetadataMatch, fields []string) error

	// Delete removes tracks from the provider's collection.
	Delete(ctx context.Context, ectx ExecutionContext, query string) (int, error)

	// Groups returns a service for managing track-to-group relationships.
	Groups() TrackGroupService

	// Sort orders a slice of tracks in-place.
	Sort(ctx context.Context, ectx ExecutionContext, tracks []models.Track, field string)
}

// TrackGroupService handles track memberships within groups (playlists).
type TrackGroupService interface {
	// Add adds tracks to a specific group.
	Add(ctx context.Context, ectx ExecutionContext, tracks []models.Track, target models.ResourceGroup) (int, error)

	// Remove removes tracks from a specific group.
	Remove(ctx context.Context, ectx ExecutionContext, tracks []models.Track, group models.ResourceGroup) (int, error)

	// Move transfers tracks from one group to another.
	Move(ctx context.Context, ectx ExecutionContext, tracks []models.Track, from models.ResourceGroup, to models.ResourceGroup) (int, error)
}

// GroupService handles structural items (Playlists, Folders) themselves.
type GroupService interface {
	// List returns groups matching a query.
	List(ctx context.Context, ectx ExecutionContext, query string) ([]models.ResourceGroup, error)

	// Create creates a new group container.
	Create(ctx context.Context, ectx ExecutionContext, parent models.ResourceGroup, name string, groupType models.GroupKind, position int) (models.ResourceGroup, error)

	// Update modifies group properties (rename, move group in tree).
	Update(ctx context.Context, ectx ExecutionContext, group models.ResourceGroup, newName string, newParent *models.ResourceGroup) error

	// Delete removes a group container.
	Delete(ctx context.Context, ectx ExecutionContext, group models.ResourceGroup) error

	// Sort orders a slice of groups in-place.
	Sort(ctx context.Context, ectx ExecutionContext, groups []models.ResourceGroup, field string)
}

// FixType defines the kind of repair operation to perform.
type FixType string

const (
	FixDuplicates FixType = "duplicates" // Remove duplicate tracks/memberships
	FixMetadata   FixType = "metadata"   // Normalize or fix metadata fields
	FixPaths      FixType = "paths"      // Repair broken file paths
	FixOrphans    FixType = "orphans"    // Remove orphaned resources
)

// FixOptions defines targets for each fix type.
type FixOptions struct {
	Actions map[FixType][]string
}

// Selection represents a resolved set of resources from a provider.
type Selection struct {
	Items    []models.Resource
	Tracks   []models.Track
	Groups   []models.ResourceGroup
	Location location.Location
}

// SystemService handles provider-wide configuration, maintenance, and orchestration.
type SystemService interface {
	Capabilities() ProviderCapabilities
	Containment() ContainmentPolicy
	MetadataCapabilities() []string
	SupportedResources() []string
	TableHeaders() []string

	// Save writes changes to persistent storage.
	Save(ctx context.Context, ectx ExecutionContext, path string) error

	// Fix performs health/formatting repairs based on options.
	Fix(ctx context.Context, ectx ExecutionContext, selection Selection, options FixOptions) (int, error)

	// Sync orchestrates a full library sync.
	Sync(ctx context.Context, ectx ExecutionContext, tracks []models.Track, targetQuery string, options SyncOptions) error

	// Identify(ctx context.Context, name string, groupType models.GroupKind) string
	Identify(name string, groupType models.GroupKind) string
}

type SyncOptions struct {
	ExportDest     string
	ExportFormat   string
	PathMaps       map[string]string
	AppendOnly     bool
	MetadataFields []string // If set, sync these metadata fields
	MatchFields    []string // Keys to use for matching tracks (e.g. artist, title, filename)
}

