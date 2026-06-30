package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/google/uuid"
	"github.com/llttlltt/dj-library-tools/internal/providers/factory"
)

// Source is a user-named, configured instance of a provider. It decouples the
// user-facing name (e.g. "Main Library") from the underlying provider type and
// connection details. Stored as ~/.config/djlt/sources/<uuid>.json.
//
// Minimal JSON format for a rekordbox source:
//
//	{
//	  "id":       "550e8400-e29b-41d4-a716-446655440000",
//	  "name":     "Main Library",
//	  "provider": "rb",
//	  "config":   { "file_path": "/Users/you/Library/rekordbox.xml" }
//	}
//
// For a Plex source the config block uses "host", "port" (string), and "token".
// For an m3u source the config block uses "file_path".
type Source struct {
	ID       string            `json:"id"`
	Name     string            `json:"name"`
	Provider string            `json:"provider"`
	Config   map[string]string `json:"config"`
}

// NewSourceID returns a new UUID v4 string.
func NewSourceID() string { return uuid.New().String() }

// LoadSources reads all *.json files from ~/.config/djlt/sources/ and returns
// them as a slice ordered lexicographically by filename.
func LoadSources() ([]Source, error) {
	dir, err := GetSourcesDir()
	if err != nil {
		return nil, err
	}
	entries, err := filepath.Glob(filepath.Join(dir, "*.json"))
	if err != nil {
		return nil, err
	}
	var out []Source
	for _, p := range entries {
		data, err := os.ReadFile(p)
		if err != nil {
			return nil, fmt.Errorf("reading source %s: %w", p, err)
		}
		var s Source
		if err := json.Unmarshal(data, &s); err != nil {
			return nil, fmt.Errorf("parsing source %s: %w", p, err)
		}
		out = append(out, s)
	}
	return out, nil
}

// SaveSource writes s to ~/.config/djlt/sources/<id>.json.
func SaveSource(s Source) error {
	dir, err := GetSourcesDir()
	if err != nil {
		return err
	}
	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(dir, s.ID+".json"), data, 0644)
}

// DeleteSource removes ~/.config/djlt/sources/<id>.json.
func DeleteSource(id string) error {
	dir, err := GetSourcesDir()
	if err != nil {
		return err
	}
	p := filepath.Join(dir, id+".json")
	if err := os.Remove(p); err != nil && !os.IsNotExist(err) {
		return err
	}
	return nil
}

// FindSourceByID loads all Sources and returns the one with the matching ID.
func FindSourceByID(id string) (*Source, error) {
	sources, err := LoadSources()
	if err != nil {
		return nil, err
	}
	for i := range sources {
		if sources[i].ID == id {
			return &sources[i], nil
		}
	}
	return nil, fmt.Errorf("no Source found with ID %q", id)
}

// FindFirstSource loads all Sources and returns the first whose Provider field
// matches the requested provider string (lexicographic file order). Returns an
// error if none is found.
func FindFirstSource(provider string) (*Source, error) {
	sources, err := LoadSources()
	if err != nil {
		return nil, err
	}
	for i := range sources {
		if sources[i].Provider == provider {
			return &sources[i], nil
		}
	}
	return nil, fmt.Errorf("no Source configured for provider %q — add one via the GUI or use -f to specify a file", provider)
}

// ResolveProviderOptions maps a Source's Config keys to the appropriate
// factory.ProviderOptions fields.
func ResolveProviderOptions(s Source) factory.ProviderOptions {
	opts := factory.ProviderOptions{}
	opts.FilePath = s.Config["file_path"]
	opts.Host = s.Config["host"]
	opts.Token = s.Config["token"]
	if port := s.Config["port"]; port != "" {
		fmt.Sscanf(port, "%d", &opts.Port) //nolint:errcheck
	}
	return opts
}
