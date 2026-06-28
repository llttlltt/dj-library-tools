package models

// Track is the provider-neutral representation of a music track.
type Track struct {
	// Identity
	ID       string
	Location string
	Display  string // Formatted display name

	// Basic Metadata
	Title   string
	Artist  string
	Album   string
	Genre   string
	Comment string
	Label   string
	Year    int
	Color   string

	// Musical Properties
	BPM  float64
	Key  string
	Size int64

	// Audio Properties
	Duration     int
	Bitrate      int
	SampleRate   int
	DateAdded    string
	DateModified string

	// Stats
	Plays  int
	Rating int // 0-255

	// DJ-Specific Metadata
	Remixer      string
	Mix          string
	CuePoints    []CuePoint
	TempoMarkers []TempoMarker

	// ImplementationState is an opaque bucket for provider-specific state (e.g. raw XML pointers).
	ImplementationState interface{}
}

func (t Track) GetID() string   { return t.ID }
func (t Track) GetName() string { return t.Title }
func (t Track) GetKind() string { return "track" }

// CuePoint represents a specific marker or performance pad in a track.
type CuePoint struct {
	Name     string
	Position float64 // In seconds
	Color    string
	Type     CueType
	Index    int // 0-indexed position within its type
}

type CueType string

const (
	CueTypeHot    CueType = "hot"
	CueTypeMemory CueType = "memory"
)

// TempoMarker represents a BPM change or beatgrid anchor point.
type TempoMarker struct {
	Position float64 // In seconds
	BPM      float64
}
