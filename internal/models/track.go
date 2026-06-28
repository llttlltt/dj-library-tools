package models

import (
	"fmt"
	"strconv"
)

// Track is the provider-neutral representation of a music track.
type Track struct {
	// Identity
	ID       string
	Location string
	Display  string

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

	// ImplementationState is an opaque bucket for provider-specific state.
	ImplementationState interface{}
}

func (t Track) GetID() string   { return t.ID }
func (t Track) GetName() string { return t.Title }
func (t Track) GetKind() string { return "track" }

// Value returns a string representation of a track property for querying.
func (t Track) Value(key string) string {
	switch key {
	case "id":
		return t.ID
	case "location":
		return t.Location
	case "display":
		return t.Display
	case "title":
		return t.Title
	case "artist":
		return t.Artist
	case "album":
		return t.Album
	case "genre":
		return t.Genre
	case "comment":
		return t.Comment
	case "label":
		return t.Label
	case "year":
		return strconv.Itoa(t.Year)
	case "color":
		return t.Color
	case "bpm":
		return fmt.Sprintf("%.2f", t.BPM)
	case "key":
		return t.Key
	case "rating":
		return strconv.Itoa(t.Rating)
	case "plays":
		return strconv.Itoa(t.Plays)
	case "added":
		return t.DateAdded
	case "modified":
		return t.DateModified
	case "bitrate":
		return strconv.Itoa(t.Bitrate)
	case "samplerate":
		return strconv.Itoa(t.SampleRate)
	case "size":
		return strconv.FormatInt(t.Size, 10)
	case "remixer":
		return t.Remixer
	case "mix":
		return t.Mix
	case "hotcues":
		return strconv.Itoa(t.countCues(CueTypeHot))
	case "memorycues":
		return strconv.Itoa(t.countCues(CueTypeMemory))
	case "beatgrids":
		return strconv.Itoa(len(t.TempoMarkers))
	case "duration":
		return strconv.Itoa(t.Duration)
	}
	return ""
}

func (t Track) countCues(cueType CueType) int {
	count := 0
	for _, cp := range t.CuePoints {
		if cp.Type == cueType {
			count++
		}
	}
	return count
}

// CuePoint represents a specific marker or performance pad in a track.
type CuePoint struct {
	Name     string
	Position float64
	Color    string
	Type     CueType
	Index    int
}

type CueType string

const (
	CueTypeHot    CueType = "hot"
	CueTypeMemory CueType = "memory"
)

// TempoMarker represents a BPM change or beatgrid anchor point.
type TempoMarker struct {
	Position float64
	BPM      float64
}
