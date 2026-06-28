package models

// Resource is the interface for any item in a music library (Track, Playlist, Folder).
type Resource interface {
	GetID() string
	GetName() string
	GetKind() string
}

type GroupKind string

const (
	GroupKindFolder   GroupKind = "folder"
	GroupKindPlaylist GroupKind = "playlist"
)

func (g GroupKind) String() string {
	return string(g)
}

// Track is the provider-neutral representation of a music track.
type Track struct {
	ID           string
	Title        string
	Artist       string
	Album        string
	Display      string
	BPM          float64
	Key          string
	Genre        string
	Comment      string
	Label        string
	Year         int
	Location     string
	Duration     int
	Rating       int // Standardized 0-255
	Plays        int
	DateAdded    string
	DateModified string
	Color        string
	Bitrate      int
	SampleRate   int
	Size         int64
	Remixer      string
	Mix          string

	// Structured DJ Metadata
	CuePoints    []CuePoint
	TempoMarkers []TempoMarker

	// ImplementationState is an opaque pointer used by providers to store 
	// state required for surgical operations (e.g., raw XML pointers).
	ImplementationState interface{}
}

func (t Track) GetID() string   { return t.ID }
func (t Track) GetName() string { return t.Title }
func (t Track) GetKind() string { return "track" }

// CuePoint represents a specific marker in a track.
type CuePoint struct {
	Name     string
	Position float64
	Color    string
	Type     CueType
	Index    int // 0-indexed position within its type
}

type CueType string

const (
	CueTypeHot    CueType = "hot"
	CueTypeMemory CueType = "memory"
	CueTypeMarker CueType = "marker"
)

// TempoMarker represents a BPM change or beatgrid anchor.
type TempoMarker struct {
	Position float64
	BPM      float64
}

// ResourceGroup represents a container like a playlist or folder.
type ResourceGroup struct {
	ID                  string
	Name                string
	Items               int
	ParentFolder        string
	Kind                GroupKind
	ImplementationState interface{}
}

func (n ResourceGroup) GetID() string   { return n.ID }
func (n ResourceGroup) GetName() string { return n.Name }
func (n ResourceGroup) GetKind() string { return n.Kind.String() }

// MetadataMatch pairs a source track with a target track for reconciliation.
type MetadataMatch struct {
	Source Track
	Target Track
}
