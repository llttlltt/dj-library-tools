package rekordbox

import (
	"testing"
)

func TestToNeutralGroup(t *testing.T) {
	n := Node{Name: "Inbox", Type: 1, Entries: PtrInt32(42)}
	got := ToNeutralGroup(n, "")
	if got.Items != 42 {
		t.Errorf("got %d, want 42", got.Items)
	}
}
