package rekordbox

import "testing"

func TestNode_ToNeutral_Entries(t *testing.T) {
	t.Run("playlist uses Entries", func(t *testing.T) {
		n := Node{Name: "Inbox", Type: 1, Entries: PtrInt32(42)}
		got := n.ToNeutral("")
		if got.Entries != 42 {
			t.Errorf("playlist Entries: got %d, want 42", got.Entries)
		}
	})

	t.Run("folder uses Count", func(t *testing.T) {
		n := Node{Name: "Shows", Type: 0, Count: PtrInt32(4)}
		got := n.ToNeutral("")
		if got.Entries != 4 {
			t.Errorf("folder Count→Entries: got %d, want 4", got.Entries)
		}
	})

	t.Run("folder with nil Count returns 0", func(t *testing.T) {
		n := Node{Name: "Empty", Type: 0}
		got := n.ToNeutral("")
		if got.Entries != 0 {
			t.Errorf("nil Count: got %d, want 0", got.Entries)
		}
	})

	t.Run("ParentFolder is propagated", func(t *testing.T) {
		n := Node{Name: "Sunset", Type: 0, Count: PtrInt32(3)}
		got := n.ToNeutral("Shows")
		if got.ParentFolder != "Shows" {
			t.Errorf("ParentFolder: got %q, want %q", got.ParentFolder, "Shows")
		}
	})
}
