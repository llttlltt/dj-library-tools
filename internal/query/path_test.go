package query

import (
	"testing"

	"github.com/llttlltt/dj-library-tools/internal/models"
)

func TestResolvePath(t *testing.T) {
	track := models.Track{
		Duration: 60,
		TempoMarkers: []models.TempoMarker{
			{Position: 0, BPM: 120.0},
			{Position: 0.5, BPM: 120.1},
			{Position: 1.0, BPM: 120.2},
		},
		CuePoints: []models.CuePoint{
			{Name: "Intro", Color: "red", Type: models.CueTypeHot, Position: 0},
			{Name: "Break", Color: "blue", Type: models.CueTypeHot, Position: 30},
		},
	}

	tests := []struct {
		name     string
		path     string
		expected string
	}{
		{"Count Beatgrids", "beatgrids-count", "3"},
		{"Beatgrid Density", "beatgrids-density", "3.00"},
		{"Beatgrid BPM Drift", "beatgrids/bpm-drift", "0.2000"},
		{"First Beatgrid BPM", "beatgrids.1/bpm", "120.0000"},
		{"Second Beatgrid BPM", "beatgrids.2/bpm", "120.1000"},
		{"Hotcue Count", "hotcues-count", "2"},
		{"First Hotcue Color", "hotcues.1/color", "red"},
		{"Second Hotcue Name", "hotcues.2/name", "Break"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			val, ok := ResolvePath(track, tt.path)
			if !ok {
				t.Errorf("ResolvePath() ok = false, want true")
			}
			if val != tt.expected {
				t.Errorf("ResolvePath() = %v, want %v", val, tt.expected)
			}
		})
	}
}
