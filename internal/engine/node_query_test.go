package engine

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
}
