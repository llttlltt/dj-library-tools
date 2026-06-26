package cli

import (
	"testing"

	"github.com/llttlltt/dj-library-tools/internal/provider"
	"github.com/llttlltt/dj-library-tools/pkg/rekordbox"
)

func TestSortTracks(t *testing.T) {
	tracks := []rekordbox.Track{
		{Name: "C", AverageBpm: "128", Artist: "Z"},
		{Name: "A", AverageBpm: "120", Artist: "X"},
		{Name: "B", AverageBpm: "124", Artist: "Y"},
	}

	t.Run("sort by bpm ascending", func(t *testing.T) {
		sortTracks(tracks, "bpm")
		if tracks[0].Name != "A" {
			t.Errorf("expected A, got %s", tracks[0].Name)
		}
	})

	t.Run("sort by bpm descending", func(t *testing.T) {
		sortTracks(tracks, "-bpm")
		if tracks[0].Name != "C" {
			t.Errorf("expected C, got %s", tracks[0].Name)
		}
	})

	t.Run("sort by artist ascending", func(t *testing.T) {
		sortTracks(tracks, "artist")
		if tracks[0].Name != "A" {
			t.Errorf("expected A (Artist X), got %s", tracks[0].Name)
		}
	})
}

func TestSortNodes(t *testing.T) {
	nodes := []provider.NodeResult{
		{Name: "Z", Entries: 10},
		{Name: "A", Entries: 50},
		{Name: "M", Entries: 5},
	}

	t.Run("sort by name ascending", func(t *testing.T) {
		sortNodes(nodes, "name")
		if nodes[0].Name != "A" {
			t.Errorf("expected A, got %s", nodes[0].Name)
		}
	})

	t.Run("sort by entries descending", func(t *testing.T) {
		sortNodes(nodes, "-entries")
		if nodes[0].Name != "A" {
			t.Errorf("expected A (50 entries), got %s", nodes[0].Name)
		}
	})
}
