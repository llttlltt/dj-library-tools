package engine

import (
	"path/filepath"
	"testing"

	"github.com/llttlltt/dj-library-tools/pkg/rekordbox"
)

func TestEnginePrimitives(t *testing.T) {
	fixturePath := filepath.Join("../../tests/fixtures/rekordbox/rekordbox_xml_2026-05-10.xml")
	lib, err := rekordbox.ReadRekordboxLibrary(fixturePath)
	if err != nil {
		t.Fatalf("Failed to read fixture: %v", err)
	}

	eng := NewEngine(lib)

	t.Run("Ls Primitive", func(t *testing.T) {
		tracks, err := eng.Ls("bpm:120..130")
		if err != nil {
			t.Errorf("Ls failed: %v", err)
		}
		if len(tracks) == 0 {
			t.Log("Note: No tracks matched in fixture, but primitive executed")
		}
	})

	t.Run("Stat Primitive", func(t *testing.T) {
		res, err := eng.Stat("") // All tracks
		if err != nil {
			t.Errorf("Stat failed: %v", err)
		}
		if res.Count != int(lib.Collection.Entries) {
			t.Errorf("Stat count mismatch: expected %d, got %d", lib.Collection.Entries, res.Count)
		}
	})

	t.Run("Modify Primitive", func(t *testing.T) {
		query := "genre:Electronic"
		changes := map[string]string{"comment": "EngineVerified"}
		
		count, err := eng.Modify(query, changes)
		if err != nil {
			t.Errorf("Modify failed: %v", err)
		}
		t.Logf("Modified %d tracks", count)
	})
}
