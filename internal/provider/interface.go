package provider

import (
	"github.com/llttlltt/dj-library-tools/internal/models"
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

// Provider defines the interface for a music library provider.
type Provider interface {
	Name() string
	GetTracks(query string) ([]models.Track, error)
	GetPlaylists(query string) ([]models.ResourceGroup, error)
	GetFolders(query string) ([]models.ResourceGroup, error)

	// Capabilities returns the feature set of this provider.
	Capabilities() ProviderCapabilities

	// GetContainmentPolicy returns the structural rules for this provider.
	GetContainmentPolicy() ContainmentPolicy

	// CustomMatch allows the provider to handle complex query fields.
	CustomMatch(track models.Track, field string, op string, value string) bool

	// CanTranscode reports whether this provider can supply raw audio for transcoding.
	CanTranscode() bool
}

// WritableProvider extends Provider with modification capabilities.
type WritableProvider interface {
	Provider
	AddTracks(target models.ResourceGroup, tracks []models.Track) (int, error)
	RemoveTracks(target models.ResourceGroup, tracks []models.Track) (int, error)
	CreateNode(parent models.ResourceGroup, name string, nodeType int) (models.ResourceGroup, error)
	DeleteNode(node models.ResourceGroup) error
	RenameNode(node models.ResourceGroup, newName string) error
	MoveNode(node models.ResourceGroup, targetParent models.ResourceGroup) error
	// Save persists any in-memory mutations to the given path.
	Save(path string) error
}
