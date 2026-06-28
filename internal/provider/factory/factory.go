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

// ProviderOptions encapsulates standard configuration for provider initialization.
type ProviderOptions struct {
	FilePath string
	Config   *config.AppConfig
	
	// Internal for testing
	MockXML *rekordbox.RekordboxLibraryXML
}

// NewProvider returns a Provider instance for the given name.
func NewProvider(name string, opts ProviderOptions) (provider.Provider, error) {
	switch name {
	case "rb", "rekordbox":
		var rbXML *rekordbox.RekordboxLibraryXML
		var err error
		
		if opts.MockXML != nil {
			rbXML = opts.MockXML
		} else {
			if opts.FilePath == "" {
				return nil, fmt.Errorf("rekordbox XML library required via --file flag")
			}
			rbXML, err = rekordbox.ReadRekordboxLibrary(opts.FilePath)
			if err != nil {
				return nil, fmt.Errorf("failed to read rekordbox library: %w", err)
			}
		}
		
		eng := library.NewEngine(rb.NewRekordboxLibrary(rbXML))
		return rb.NewRekordboxProviderWithXML(eng, rbXML, opts.FilePath), nil
	case "plex":
		token := os.Getenv("PLEX_TOKEN")
		if token == "" && opts.Config != nil {
			token = opts.Config.PlexToken
		}
		if token == "" {
			return nil, fmt.Errorf("plex token not found; run 'djlt auth plex' or set PLEX_TOKEN")
		}
		return plex.NewPlexProvider(token, opts.Config.PlexHost, opts.Config.PlexPort), nil
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
