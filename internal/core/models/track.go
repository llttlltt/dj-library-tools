package models

import (
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
	Playlists    []PlaylistMembership

	// ImplementationState is an opaque bucket for provider-specific state.
	ImplementationState interface{}
}

// PlaylistMembership represents a track's presence in a specific playlist.
type PlaylistMembership struct {
	Name   string
	Folder string
}

func (t Track) GetID() string   { return t.ID }
func (t Track) GetName() string { return t.Title }
func (t Track) GetKind() string { return "track" }

// Value returns a string representation of a track property for querying.
func (t Track) Value(key string) string {
	if def, ok := TrackFields[key]; ok {
		return def.Accessor(t)
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

func (t Track) Hotcues() int    { return t.countCues(CueTypeHot) }
func (t Track) Memorycues() int { return t.countCues(CueTypeMemory) }
func (t Track) Beatgrids() int  { return len(t.TempoMarkers) }

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
