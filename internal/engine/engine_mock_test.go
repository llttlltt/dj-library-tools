package engine

import (
	"testing"

	"github.com/llttlltt/dj-library-tools/pkg/rekordbox"
)

// MockLibrary satisfies the Library interface for unit testing.
type MockLibrary struct {
	Tracks    []rekordbox.Track
	Playlists []rekordbox.Node
}

func (m *MockLibrary) GetTracks() []rekordbox.Track {
	return m.Tracks
}

func (m *MockLibrary) GetPlaylists() []rekordbox.Node {
	return m.Playlists
}

func TestEngineWithMockLibrary(t *testing.T) {
	mock := &MockLibrary{
		Tracks: []rekordbox.Track{
			{TrackID: 1, Name: "Track 1", Artist: "Artist A", AverageBpm: "120.0"},
			{TrackID: 2, Name: "Track 2", Artist: "Artist B", AverageBpm: "130.0"},
		},
		Playlists: []rekordbox.Node{
			{Name: "Inbox", Type: 1, Entries: rekordbox.PtrInt32(2)},
		},
	}

	eng := NewEngine(mock)

	t.Run("Match simple track", func(t *testing.T) {
		tracks, err := eng.Ls("artist:'Artist A'")
		if err != nil {
			t.Fatalf("Ls failed: %v", err)
		}
		if len(tracks) != 1 || tracks[0].Name != "Track 1" {
			t.Errorf("Expected Track 1, got %v", tracks)
		}
	})

	t.Run("Match BPM range", func(t *testing.T) {
		tracks, err := eng.Ls("bpm:>125")
		if err != nil {
			t.Fatalf("Ls failed: %v", err)
		}
		if len(tracks) != 1 || tracks[0].Name != "Track 2" {
			t.Errorf("Expected Track 2, got %v", tracks)
		}
	})
}
