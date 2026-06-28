package library

import (
	"testing"
)

func TestEngine_LsPlaylists(t *testing.T) {
	mock := makeMockLibrary()
	eng := NewEngine(mock)

	t.Run("filter by name", func(t *testing.T) {
		results, _ := eng.LsPlaylists("name:Inbox")
		if len(results) != 1 {
			t.Fatalf("expected 1 playlist, got %d", len(results))
		}
		if results[0].Name != "Inbox" {
			t.Errorf("expected Inbox, got %s", results[0].Name)
		}
	})

	t.Run("does not return folders", func(t *testing.T) {
		results, _ := eng.LsPlaylists("")
		for _, r := range results {
			if r.Type != 1 {
				t.Errorf("LsPlaylists returned a non-playlist node: %s (Type=%d)", r.Name, r.Type)
			}
		}
	})
}

func TestEngine_LsFolders(t *testing.T) {
	mock := makeMockLibrary()
	eng := NewEngine(mock)

	t.Run("filter by name", func(t *testing.T) {
		results, _ := eng.LsFolders("name:Sets")
		if len(results) != 1 {
			t.Fatalf("expected 1 folder, got %d", len(results))
		}
		if results[0].Name != "Sets" {
			t.Errorf("expected Sets, got %s", results[0].Name)
		}
	})

	t.Run("does not return playlists", func(t *testing.T) {
		results, _ := eng.LsFolders("")
		for _, r := range results {
			if r.Type != 0 {
				t.Errorf("LsFolders returned a non-folder node: %s (Type=%d)", r.Name, r.Type)
			}
		}
	})
}
