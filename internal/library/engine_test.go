package library

import (
	"testing"
)

func TestEngine_Ls(t *testing.T) {
	mock := makeMockLibrary()
	eng := NewEngine(mock)

	t.Run("filter by artist", func(t *testing.T) {
		tracks, _ := eng.Ls("artist:Artist A")
		if len(tracks) != 1 {
			t.Fatalf("expected 1 track, got %d", len(tracks))
		}
		if tracks[0].Title != "Track 1" {
			t.Errorf("expected Track 1, got %s", tracks[0].Title)
		}
	})

	t.Run("filter by bpm range", func(t *testing.T) {
		tracks, _ := eng.Ls("bpm:125..130")
		if len(tracks) != 1 {
			t.Fatalf("expected 1 track, got %d", len(tracks))
		}
		if tracks[0].Title != "Track 2" {
			t.Errorf("expected Track 2, got %s", tracks[0].Title)
		}
	})
}
