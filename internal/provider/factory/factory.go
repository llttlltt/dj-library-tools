package factory

import (
	"fmt"
	"os"

	"github.com/llttlltt/dj-library-tools/internal/config"
	"github.com/llttlltt/dj-library-tools/internal/library"
	"github.com/llttlltt/dj-library-tools/internal/provider"
	"github.com/llttlltt/dj-library-tools/internal/provider/m3u"
	"github.com/llttlltt/dj-library-tools/internal/provider/plex"
	"github.com/llttlltt/dj-library-tools/internal/provider/rb"
	"github.com/llttlltt/dj-library-tools/internal/rekordbox"
)

// NewProvider returns a Provider instance for the given name.
func NewProvider(name string, rbXML *rekordbox.RekordboxLibraryXML, filePath string, cfg *config.AppConfig) (provider.Provider, error) {
	switch name {
	case "rb", "rekordbox":
		if rbXML == nil {
			return nil, fmt.Errorf("rekordbox XML library required")
		}
		eng := library.NewEngine(rb.NewRekordboxLibrary(rbXML))
		return rb.NewRekordboxProvider(eng, filePath), nil
	case "plex":
		token := os.Getenv("PLEX_TOKEN")
		if token == "" {
			token = cfg.PlexToken
		}
		if token == "" {
			return nil, fmt.Errorf("plex token not found; run 'djlt auth plex' or set PLEX_TOKEN")
		}
		return plex.NewPlexProvider(token, cfg.PlexHost, cfg.PlexPort), nil
	case "m3u", "m3u8":
		p, err := m3u.NewM3UProvider("")
		if err != nil {
			return nil, err
		}
		return p, nil
	default:
		return nil, fmt.Errorf("unknown provider: %s", name)
	}
}
