package provider

import (
	"errors"

	"github.com/llttlltt/dj-library-tools/internal/models"
	"github.com/llttlltt/dj-library-tools/internal/query"
)

var (
	ErrReadOnly            = errors.New("provider is read-only")
	ErrInvalidParent       = errors.New("invalid parent for this group type")
	ErrUnsupportedResource = errors.New("resource type not supported by this provider")
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
	SupportedResources() []string
	// MetadataCapabilities returns a list of fields this provider can serve/update.
	MetadataCapabilities() []string
}

// ReadableProvider extends BaseProvider with read operations.
type ReadableProvider interface {
	BaseProvider
	// GetResources is the primary discovery method for all resource types (tracks, playlists, folders).
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
	UpdateTracks(ctx ExecutionContext, query string, changes map[string]string) (int, error)
	MoveTracks(ctx ExecutionContext, source models.ResourceGroup, target models.ResourceGroup, tracks []models.Track) (int, error)
	CreateGroup(ctx ExecutionContext, parent models.ResourceGroup, name string, groupType models.GroupType, position int) (models.ResourceGroup, error)
	DeleteGroup(ctx ExecutionContext, node models.ResourceGroup) error
	RenameGroup(ctx ExecutionContext, node models.ResourceGroup, newName string, groupType models.GroupType) error
	MoveGroup(ctx ExecutionContext, node models.ResourceGroup, targetParent models.ResourceGroup) error
	
	Sync(ctx ExecutionContext, tracks []models.Track, sourceQuery string, targetQuery string, options SyncOptions) error
	
	// UpdateMetadata applies specific field updates to matched tracks.
	UpdateMetadata(ctx ExecutionContext, matches []models.MetadataMatch, fields []string) error

	// Fix performs provider-specific health/formatting repairs (e.g. M3U tag enrichment).
	Fix(ctx ExecutionContext, resource string, query string) error


	// Validation methods for pre-flight checks
	ValidateAddTracks(target models.ResourceGroup) error
	ValidateMoveGroup(src models.ResourceGroup, target models.ResourceGroup) error
	ValidateCreateGroup(parent models.ResourceGroup, groupType models.GroupType) error

	// IdentifyGroup returns the provider-specific ID for a group name and type.
	IdentifyGroup(name string, groupType models.GroupType) string

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

func ToTrackSlice(res []models.Resource) []models.Track {
	var tracks []models.Track
	for _, r := range res {
		if t, ok := r.(models.Track); ok {
			tracks = append(tracks, t)
		}
	}
	return tracks
}
