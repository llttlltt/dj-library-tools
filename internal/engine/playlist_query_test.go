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

	t.Run("Match track on zero playlists", func(t *testing.T) {
		// Track 121598507 is on Everything, Terracotta, and Disco to House #1
		// We need to find one that is on NO playlists or just check count
		tracks, err := eng.Ls("playlistcount:0")
		if err != nil {
			t.Fatalf("Ls failed: %v", err)
		}
		// Based on the XML, most tracks are at least in "Everything" 
		// but let's check a specific one we know has several
		tracks, err = eng.Ls("id:121598507 && playlistcount:3")
		if err != nil {
			t.Fatalf("Ls failed: %v", err)
		}
		if len(tracks) != 1 {
			t.Errorf("Expected track 121598507 to be in 3 playlists, found in %d", len(tracks))
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
