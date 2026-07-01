package factory

import (
	"fmt"
	"sort"
	"sync"

	provider "github.com/llttlltt/dj-library-tools/internal/providers"
)

type ProviderFactory func(opts ProviderOptions) (provider.Provider, error)

// ProviderInfo carries static metadata and capabilities for a provider type.
type ProviderInfo struct {
	Name         string                        `json:"name"`
	Resources    []ResourceInfo                `json:"resources"`
	Capabilities provider.ProviderCapabilities `json:"capabilities"`
}

// ResourceInfo defines the capabilities of a specific resource type.
type ResourceInfo struct {
	Name          string `json:"name"`
	CanWrite      bool   `json:"can_write"`
	SupportsQuery bool   `json:"supports_query"`
}

type registeredProvider struct {
	factory func(opts ProviderOptions) (provider.Provider, error)
	info    ProviderInfo
}

var (
	providersMu sync.RWMutex
	providers   = make(map[string]registeredProvider)
)

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
func Register(name string, resources []ResourceInfo, caps provider.ProviderCapabilities, factory ProviderFactory) {
	providersMu.Lock()
	defer providersMu.Unlock()
	if factory == nil {
		panic("provider: Register factory is nil")
	}
	if _, dup := providers[name]; dup {
		panic("provider: Register called twice for factory " + name)
	}
	providers[name] = registeredProvider{
		factory: factory,
		info: ProviderInfo{
			Name:         name,
			Resources:    resources,
			Capabilities: caps,
		},
	}
}

// ValidateResource checks if a provider supports a given resource type and optionally if it's writable.
func ValidateResource(providerName, resource string, mustBeWritable bool) bool {
	providersMu.RLock()
	defer providersMu.RUnlock()
	p, ok := providers[providerName]
	if !ok {
		return false
	}
	for _, r := range p.info.Resources {
		if r.Name == resource {
			if mustBeWritable && !r.CanWrite {
				return false
			}
			return true
		}
	}
	return false
}

// GetProviderInfo returns static metadata for a provider by name.
func GetProviderInfo(name string) (ProviderInfo, error) {
	providersMu.RLock()
	defer providersMu.RUnlock()
	p, ok := providers[name]
	if !ok {
		return ProviderInfo{}, fmt.Errorf("unknown provider: %s", name)
	}
	return p.info, nil
}

// ListProviders returns all registered provider metadata, sorted by name.
func ListProviders() []ProviderInfo {
	providersMu.RLock()
	defer providersMu.RUnlock()
	var out []ProviderInfo
	for _, p := range providers {
		out = append(out, p.info)
	}
	sort.Slice(out, func(i, j int) bool {
		return out[i].Name < out[j].Name
	})
	return out
}

// NewProvider returns a Provider instance for the given name.
func NewProvider(name string, opts ProviderOptions) (provider.Provider, error) {
	providersMu.RLock()
	p, ok := providers[name]
	providersMu.RUnlock()
	if !ok {
		return nil, fmt.Errorf("unknown provider: %s (ensure it is registered)", name)
	}
	prov, err := p.factory(opts)
	if err != nil {
		return nil, err
	}
	return &provider.GatedProvider{Base: prov}, nil
}
