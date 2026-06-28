package sync

import (
	"testing"
	"github.com/llttlltt/dj-library-tools/internal/provider/rb"
	"github.com/llttlltt/dj-library-tools/internal/rekordbox"
)

func TestSync(t *testing.T) {
	lib := &rekordbox.RekordboxLibraryXML{}
	_ = rb.NewRekordboxLibrary(lib)
}
