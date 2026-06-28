package models

import (
	"fmt"
	"strconv"
)

// Track is the provider-neutral representation of a music track.
type Track struct {
	// Identity
	ID       string `query:"id"`
	Location string `query:"location"`
	Display  string `query:"display"`

	// Basic Metadata
	Title   string `query:"title"`
	Artist  string `query:"artist"`
	Album   string `query:"album"`
	Genre   string `query:"genre"`
	Comment string `query:"comment"`
	Label   string `query:"label"`
	Year    int    `query:"year,numeric"`
	Color   string `query:"color"`

	// Musical Properties
	BPM  float64 `query:"bpm,numeric"`
	Key  string  `query:"key"`
	Size int64   `query:"size,numeric"`

	// Audio Properties
	Duration     int    `query:"duration,numeric"`
	Bitrate      int    `query:"bitrate,numeric"`
	SampleRate   int    `query:"samplerate,numeric"`
	DateAdded    string `query:"added"`
	DateModified string `query:"modified"`

	// Stats
	Plays  int `query:"plays,numeric"`
	Rating int `query:"rating,numeric"` // 0-255

	// DJ-Specific Metadata
	Remixer      string        `query:"remixer"`
	Mix          string        `query:"mix"`
	CuePoints    []CuePoint    `query:"-"`
	TempoMarkers []TempoMarker `query:"-"`

	// ImplementationState is an opaque bucket for provider-specific state.
	ImplementationState interface{} `query:"-"`
}

func (t Track) GetID() string   { return t.ID }
func (t Track) GetName() string { return t.Title }
func (t Track) GetKind() string { return "track" }

// Derived Helpers for Query Engine
func (t Track) Hotcues() int    { return t.countCues(CueTypeHot) }
func (t Track) Memorycues() int { return t.countCues(CueTypeMemory) }
func (t Track) Beatgrids() int  { return len(t.TempoMarkers) }

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

// GetQueryValue provides a fallback for calculated/derived fields.
func (t Track) GetQueryValue(field string) (string, bool) {
	switch field {
	case "hotcues":
		return strconv.Itoa(t.Hotcues()), true
	case "memorycues":
		return strconv.Itoa(t.Memorycues()), true
	case "beatgrids":
		return strconv.Itoa(t.Beatgrids()), true
	case "bpm":
		return fmt.Sprintf("%.2f", t.BPM), true
	}
	return "", false
}
