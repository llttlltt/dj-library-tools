package provider

import (
	"fmt"
	"github.com/llttlltt/dj-library-tools/internal/engine"
	"github.com/llttlltt/dj-library-tools/pkg/rekordbox"
)

// NewProvider returns a Provider instance for the given name.
func NewProvider(name string, rbXML *rekordbox.RekordboxLibraryXML) (Provider, error) {
	switch name {
	case "rb", "rekordbox":
		eng := engine.NewEngine(engine.NewRekordboxLibrary(rbXML))
		return NewRekordboxProvider(eng), nil
	case "plex":
		// Plex requires a token and host, which are currently handled in the CLI.
		// For now, we'll return an error or a dummy to satisfy the interface.
		return nil, fmt.Errorf("plex provider initialization not yet moved to provider package")
	default:
		return nil, fmt.Errorf("unknown provider: %s", name)
	}
}
