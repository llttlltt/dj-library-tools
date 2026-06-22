package engine

import (
	"path/filepath"
	"testing"

	"github.com/llttlltt/dj-library-tools/pkg/rekordbox"
)

func TestPlaylistQuery(t *testing.T) {
	fixturePath := filepath.Join("../../tests/fixtures/rekordbox/rekordbox_xml_2026-06-22.xml")
	lib, err := rekordbox.ReadRekordboxLibrary(fixturePath)
	if err != nil {
		t.Fatalf("Failed to read fixture: %v", err)
	}

	eng := NewEngine(lib)

	t.Run("Match track in playlist", func(t *testing.T) {
		// "Everything" playlist has track 261282774
		tracks, err := eng.Ls("playlist:Everything && id:261282774")
		if err != nil {
			t.Fatalf("Ls failed: %v", err)
		}
		if len(tracks) != 1 {
			t.Errorf("Expected 1 track, got %d", len(tracks))
		}
	})

	t.Run("Match track in nested playlist", func(t *testing.T) {
		// Shows > Mike's BBQ
		tracks, err := eng.Ls("playlist:\"Mike's BBQ\" && id:265849715")
		if err != nil {
			t.Fatalf("Ls failed: %v", err)
		}
		if len(tracks) != 1 {
			t.Errorf("Expected 1 track, got %d", len(tracks))
		}
	})

	t.Run("Match track by exact playlist count", func(t *testing.T) {
		// Track 121598507 appears in exactly 13 unique playlists in the fixture.
		// Previously playlistcount:3 matched it via substring ("13" contains "3");
		// exact numeric equality should not match.
		tracks, err := eng.Ls("id:121598507 && playlistcount:3")
		if err != nil {
			t.Fatalf("Ls failed: %v", err)
		}
		if len(tracks) != 0 {
			t.Errorf("playlistcount:3 should not match a track in 13 playlists, got %d", len(tracks))
		}

		// Now verify the real count matches.
		tracks, err = eng.Ls("id:121598507 && playlistcount:13")
		if err != nil {
			t.Fatalf("Ls failed: %v", err)
		}
		if len(tracks) != 1 {
			t.Errorf("Expected track 121598507 to be in 13 playlists, got %d results", len(tracks))
		}
	})

	t.Run("Match track in two specific playlists", func(t *testing.T) {
		// Track 267482775 is in "Mike's BBQ" and "Terracotta"
		tracks, err := eng.Ls("playlist:\"Mike's BBQ\" && playlist:Terracotta && id:267482775")
		if err != nil {
			t.Fatalf("Ls failed: %v", err)
		}
		if len(tracks) != 1 {
			t.Errorf("Expected 1 track, got %d", len(tracks))
		}
	})
}
