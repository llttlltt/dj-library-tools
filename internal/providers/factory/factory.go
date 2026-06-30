package factory

import (
	"fmt"
	"sync"

	provider "github.com/llttlltt/dj-library-tools/internal/providers"
)

var (
	providersMu sync.RWMutex
	providers   = make(map[string]ProviderFactory)
)

type ProviderFactory func(opts ProviderOptions) (provider.Provider, error)

// ProviderOptions carries resolved Source connection fields. Each provider
// factory reads from the fields relevant to its type:
//   - rekordbox / m3u: FilePath
//   - plex:            Host, Port, Token
type ProviderOptions struct {
	FilePath string
	Host     string
	Port     int
	Token    string
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
	prov, err := factory(opts)
	if err != nil {
		return nil, err
	}
	return &provider.GatedProvider{Base: prov}, nil
}
