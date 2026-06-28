package provider

import (
	"github.com/llttlltt/dj-library-tools/internal/models"
	"github.com/llttlltt/dj-library-tools/internal/query"
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

// BaseProvider defines methods common to all providers.
type BaseProvider interface {
	Name() string
	Capabilities() ProviderCapabilities
	GetContainmentPolicy() ContainmentPolicy
	CustomMatch(track models.Track, field string, op query.Operator, value string) bool
	CanTranscode() bool
}

// ReadableProvider extends BaseProvider with read operations.
type ReadableProvider interface {
	BaseProvider
	GetTracks(ctx ExecutionContext, query string) ([]models.Track, error)
	GetPlaylists(ctx ExecutionContext, query string) ([]models.ResourceGroup, error)
	GetFolders(ctx ExecutionContext, query string) ([]models.ResourceGroup, error)

	// GetResources is the unified entry point for all resource types.
	GetResources(ctx ExecutionContext, resource string, query string) ([]models.Resource, error)

	// Sort operations
	SortTracks(ctx ExecutionContext, tracks []models.Track, field string)
	SortGroups(ctx ExecutionContext, groups []models.ResourceGroup, field string)
}

// SearchableProvider is an optional interface for providers with server-side search.
type SearchableProvider interface {
	ReadableProvider
	Search(ctx ExecutionContext, query string) ([]models.Resource, error)
}

// WritableProvider extends ReadableProvider with modification capabilities.
type WritableProvider interface {
	ReadableProvider
	AddTracks(ctx ExecutionContext, target models.ResourceGroup, tracks []models.Track) (int, error)
	RemoveTracks(ctx ExecutionContext, target models.ResourceGroup, tracks []models.Track) (int, error)
	CreateGroup(ctx ExecutionContext, parent models.ResourceGroup, name string, nodeType int, position int) (models.ResourceGroup, error)
	DeleteGroup(ctx ExecutionContext, node models.ResourceGroup) error
	RenameGroup(ctx ExecutionContext, node models.ResourceGroup, newName string) error
	MoveGroup(ctx ExecutionContext, node models.ResourceGroup, targetParent models.ResourceGroup) error
	
	Sync(ctx ExecutionContext, tracks []models.Track, sourceQuery string, targetQuery string, options SyncOptions) error
	ModifyTracks(ctx ExecutionContext, query string, changes map[string]string) (int, error)

	// Validation methods for pre-flight checks
	ValidateAddTracks(target models.ResourceGroup) error
	ValidateMoveGroup(src models.ResourceGroup, target models.ResourceGroup) error
	ValidateCreateGroup(parent models.ResourceGroup, groupType models.GroupType) error

	Save(ctx ExecutionContext, path string) error
}

// Provider is an alias for ReadableProvider as the standard return type.
type Provider = ReadableProvider

type SyncOptions struct {
	ExportDest   string
	ExportFormat string
	PathMaps     map[string]string
	AppendOnly   bool
}
