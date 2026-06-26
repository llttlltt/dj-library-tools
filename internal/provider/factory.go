package provider

import (
	"fmt"
	"os"

	"github.com/llttlltt/dj-library-tools/internal/config"
	"github.com/llttlltt/dj-library-tools/internal/engine"
	"github.com/llttlltt/dj-library-tools/pkg/rekordbox"
)

// NewProvider returns a Provider instance for the given name.
func NewProvider(name string, rbXML *rekordbox.RekordboxLibraryXML, cfg *config.AppConfig) (Provider, error) {
	switch name {
	case "rb", "rekordbox":
		if rbXML == nil {
			return nil, fmt.Errorf("rekordbox XML library required")
		}
		eng := engine.NewEngine(engine.NewRekordboxLibrary(rbXML))
		return NewRekordboxProvider(eng), nil
	case "plex":
		token := os.Getenv("PLEX_TOKEN")
		if token == "" {
			token = cfg.PlexToken
		}
		if token == "" {
			return nil, fmt.Errorf("plex token not found; run 'djlt auth plex' or set PLEX_TOKEN")
		}
		return NewPlexProvider(token, cfg.PlexHost, cfg.PlexPort), nil
	default:
		return nil, fmt.Errorf("unknown provider: %s", name)
	}
}
