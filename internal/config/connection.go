package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/google/uuid"
	"github.com/llttlltt/dj-library-tools/internal/providers/factory"
)

// Connection is a user-named, configured instance of a provider. It decouples the
// user-facing name (e.g. "Main Library") from the underlying provider type and
// connection details. Stored as ~/.config/djlt/connections/<uuid>.json.
//
// Minimal JSON format for a rekordbox connection:
//
//	{
//	  "id":       "550e8400-e29b-41d4-a716-446655440000",
//	  "name":     "Main Library",
//	  "provider": "rb",
//	  "config":   { "file_path": "/Users/you/Library/rekordbox.xml" }
//	}
//
// For a Plex connection the config block uses "host", "port" (string), and "token".
// For an m3u connection the config block uses "file_path".
type Connection struct {
	ID       string            `json:"id"`
	Name     string            `json:"name"`
	Provider string            `json:"provider"`
	Config   map[string]string `json:"config"`
}

// NewConnectionID returns a new UUID v4 string.
func NewConnectionID() string { return uuid.New().String() }

// LoadConnections reads all *.json files from ~/.config/djlt/connections/ and returns
// them as a slice ordered lexicographically by filename.
func LoadConnections() ([]Connection, error) {
	dir, err := GetConnectionsDir()
	if err != nil {
		return nil, err
	}
	entries, err := filepath.Glob(filepath.Join(dir, "*.json"))
	if err != nil {
		return nil, err
	}
	var out []Connection
	for _, p := range entries {
		data, err := os.ReadFile(p)
		if err != nil {
			return nil, fmt.Errorf("reading connection %s: %w", p, err)
		}
		var s Connection
		if err := json.Unmarshal(data, &s); err != nil {
			return nil, fmt.Errorf("parsing connection %s: %w", p, err)
		}
		out = append(out, s)
	}
	return out, nil
}

// SaveConnection writes s to ~/.config/djlt/connections/<id>.json.
func SaveConnection(c Connection) error {
	dir, err := GetConnectionsDir()
	if err != nil {
		return err
	}
	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(dir, c.ID+".json"), data, 0644)
}

// DeleteConnection removes ~/.config/djlt/connections/<id>.json.
func DeleteConnection(id string) error {
	dir, err := GetConnectionsDir()
	if err != nil {
		return err
	}
	p := filepath.Join(dir, id+".json")
	if err := os.Remove(p); err != nil && !os.IsNotExist(err) {
		return err
	}
	return nil
}

// FindConnectionByID loads all Connections and returns the one with the matching ID.
func FindConnectionByID(id string) (*Connection, error) {
	connections, err := LoadConnections()
	if err != nil {
		return nil, err
	}
	for i := range connections {
		if connections[i].ID == id {
			return &connections[i], nil
		}
	}
	return nil, fmt.Errorf("no Connection found with ID %q", id)
}

// FindFirstConnection loads all Connections and returns the first whose Provider field
// matches the requested provider string (lexicographic file order). Returns an
// error if none is found.
func FindFirstConnection(provider string) (*Connection, error) {
	connections, err := LoadConnections()
	if err != nil {
		return nil, err
	}
	for i := range connections {
		if connections[i].Provider == provider {
			return &connections[i], nil
		}
	}
	return nil, fmt.Errorf("no Connection configured for provider %q — add one via the GUI or use -f to specify a file", provider)
}

// ResolveProviderOptions maps a Connection's Config keys to the appropriate
// factory.ProviderOptions fields.
func ResolveProviderOptions(c Connection) factory.ProviderOptions {
	opts := factory.ProviderOptions{}
	opts.FilePath = c.Config["file_path"]
	opts.Host = c.Config["host"]
	opts.Token = c.Config["token"]
	if port := c.Config["port"]; port != "" {
		fmt.Sscanf(port, "%d", &opts.Port) //nolint:errcheck
	}
	return opts
}
