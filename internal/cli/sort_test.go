package cli

import (
	"testing"

	"github.com/llttlltt/dj-library-tools/internal/models"
)

func TestSortTracks(t *testing.T) {
	tracks := []models.Track{
		{Title: "C", BPM: 128, Artist: "Z"},
		{Title: "A", BPM: 120, Artist: "X"},
		{Title: "B", BPM: 124, Artist: "Y"},
	}

	t.Run("sort by bpm ascending", func(t *testing.T) {
		sortTracks(tracks, "bpm")
		if tracks[0].Title != "A" {
			t.Errorf("expected A, got %s", tracks[0].Title)
		}
	})

	t.Run("sort by bpm descending", func(t *testing.T) {
		sortTracks(tracks, "-bpm")
		if tracks[0].Title != "C" {
			t.Errorf("expected C, got %s", tracks[0].Title)
		}
	})

	t.Run("sort by artist ascending", func(t *testing.T) {
		sortTracks(tracks, "artist")
		if tracks[0].Title != "A" {
			t.Errorf("expected A (Artist X), got %s", tracks[0].Title)
		}
	})
}

func TestSortNodes(t *testing.T) {
	nodes := []models.Node{
		{Name: "Z", Items: 10},
		{Name: "A", Items: 50},
		{Name: "M", Items: 5},
	}

	t.Run("sort by name ascending", func(t *testing.T) {
		sortNodes(nodes, "name")
		if nodes[0].Name != "A" {
			t.Errorf("expected A, got %s", nodes[0].Name)
		}
	})

	t.Run("sort by items descending", func(t *testing.T) {
		sortNodes(nodes, "-items")
		if nodes[0].Name != "A" {
			t.Errorf("expected A (50 items), got %s", nodes[0].Name)
		}
	})
}
