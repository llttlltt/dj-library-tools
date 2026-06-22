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

	t.Run("Match track by partial playlist name", func(t *testing.T) {
		tracks, err := eng.Ls("playlist:BBQ && id:265849715")
		if err != nil {
			t.Fatalf("Ls failed: %v", err)
		}
		if len(tracks) != 1 {
			t.Errorf("Expected 1 track, got %d", len(tracks))
		}
	})
}
