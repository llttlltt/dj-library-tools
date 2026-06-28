package factory

import (
	"fmt"
	"sync"

	"github.com/llttlltt/dj-library-tools/internal/config"
	"github.com/llttlltt/dj-library-tools/internal/provider"
	"github.com/llttlltt/dj-library-tools/internal/rekordbox"
)

var (
	providersMu sync.RWMutex
	providers   = make(map[string]ProviderFactory)
)

type ProviderFactory func(opts ProviderOptions) (provider.Provider, error)

type ProviderOptions struct {
	FilePath string
	Config   *config.AppConfig
	
	// MockXML is for internal testing hooks
	MockXML *rekordbox.RekordboxLibraryXML
}

// Register makes a provider factory available by the provided name.
func Register(name string, factory ProviderFactory) {
	providersMu.Lock()
	defer providersMu.Unlock()
	if factory == nil {
		panic("provider: Register factory is nil")
	}
	if _, dup := providers[name]; dup {
		panic("provider: Register called twice for factory " + name)
	}
	providers[name] = factory
}

// NewProvider returns a Provider instance for the given name.
func NewProvider(name string, opts ProviderOptions) (provider.Provider, error) {
	providersMu.RLock()
	factory, ok := providers[name]
	providersMu.RUnlock()
	if !ok {
		return nil, fmt.Errorf("unknown provider: %s (ensure it is registered)", name)
	}
	return factory(opts)
}
