package rb

import (
	"testing"
	"github.com/llttlltt/dj-library-tools/internal/rekordbox"
)

func TestToNeutralGroup(t *testing.T) {
	n := rekordbox.Node{Name: "Inbox", Type: 1, Entries: rekordbox.PtrInt32(42)}
	got := ToNeutralGroup(n, "")
	if got.Items != 42 {
		t.Errorf("got %d, want 42", got.Items)
	}
}
