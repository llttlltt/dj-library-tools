package models


// Resource is the interface for any item in a music library (Track, Playlist, Folder).
type Resource interface {
	GetID() string
	GetName() string
	GetKind() string // "track", "playlist", "folder"
}

// Track is the provider-neutral representation of a music track.
type Track struct {
	ID            string
	Title         string
	Artist        string
	Album         string
	Display       string // Provider-specific display name (e.g. M3U display string)
	BPM           float64
	Key           string
	Genre         string
	Comment       string
	Label         string
	Year          int
	Location      string
	Duration      int
	Rating        int // 0-255
	Plays         int
	DateAdded     string
	DateModified  string
	Color         string
	Bitrate       int
	SampleRate    int
	Size          int64
	Remixer       string
	Mix           string

	// Advanced DJ Metadata
	CuePoints    []CuePoint
	TempoMarkers []TempoMarker

	Raw interface{}
}

// CuePoint represents a specific marker in a track (HotCue or Memory Cue).
type CuePoint struct {
	Name     string
	Position float64 // In seconds
	Color    string
	Num      int // -1 for memory cues, 0+ for hotcues
}

// TempoMarker represents a BPM marker for a beatgrid.
type TempoMarker struct {
	Position float64 // In seconds
	BPM      float64
}

func (t Track) GetID() string   { return t.ID }
func (t Track) GetName() string { return t.Title }
func (t Track) GetKind() string { return "track" }

// ResourceGroup represents a container like a playlist or folder.
type ResourceGroup struct {
	ID           string
	Name         string
	Items        int
	ParentFolder string
	Type         GroupType
	Raw          interface{}
}

type GroupType int

const (
	GroupTypeFolder   GroupType = 0
	GroupTypePlaylist GroupType = 1
)

func (g GroupType) String() string {
	if g == GroupTypeFolder {
		return "folder"
	}
	return "playlist"
}

func (n ResourceGroup) GetID() string   { return n.ID }
func (n ResourceGroup) GetName() string { return n.Name }
func (n ResourceGroup) GetKind() string {
	return n.Type.String()
}

// MetadataMatch pairs a source track with a target track for reconciliation.
type MetadataMatch struct {
	Source Track
	Target Track
}

// NormalizeRating scales a rating from a source range (e.g. 0-5) to our 0-255 standard.
func NormalizeRating(val float64, max float64) int {
	if max == 0 { return 0 }
	return int((val / max) * 255)
}

// ScaleRating scales our 0-255 rating back to a provider-specific range.
func ScaleRating(val int, max float64) float64 {
	return (float64(val) / 255.0) * max
}
