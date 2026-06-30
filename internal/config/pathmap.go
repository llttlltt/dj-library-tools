package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/google/uuid"
)

// PathRule is a single path-translation rule within a PathMap.
type PathRule struct {
	From string `json:"from"`
	To   string `json:"to"`
}

// PathMap declares a path-translation relationship between two Sources.
// It enables the sync engine to reconcile file paths across providers.
// Stored as ~/.config/djlt/path-maps/<uuid>.json.
type PathMap struct {
	ID        string     `json:"id"`
	SourceAID string     `json:"source_a_id"`
	SourceBID string     `json:"source_b_id"`
	Rules     []PathRule `json:"rules"`
}

// NewPathMapID returns a new UUID v4 string for a PathMap.
func NewPathMapID() string { return uuid.New().String() }

// LoadPathMaps reads all *.json files from ~/.config/djlt/path-maps/.
func LoadPathMaps() ([]PathMap, error) {
	dir, err := GetPathMapsDir()
	if err != nil {
		return nil, err
	}
	entries, err := filepath.Glob(filepath.Join(dir, "*.json"))
	if err != nil {
		return nil, err
	}
	var out []PathMap
	for _, p := range entries {
		data, err := os.ReadFile(p)
		if err != nil {
			return nil, fmt.Errorf("reading path-map %s: %w", p, err)
		}
		var pm PathMap
		if err := json.Unmarshal(data, &pm); err != nil {
			return nil, fmt.Errorf("parsing path-map %s: %w", p, err)
		}
		out = append(out, pm)
	}
	return out, nil
}

// SavePathMap writes pm to ~/.config/djlt/path-maps/<id>.json.
func SavePathMap(pm PathMap) error {
	dir, err := GetPathMapsDir()
	if err != nil {
		return err
	}
	data, err := json.MarshalIndent(pm, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(dir, pm.ID+".json"), data, 0644)
}

// DeletePathMap removes ~/.config/djlt/path-maps/<id>.json.
func DeletePathMap(id string) error {
	dir, err := GetPathMapsDir()
	if err != nil {
		return err
	}
	p := filepath.Join(dir, id+".json")
	if err := os.Remove(p); err != nil && !os.IsNotExist(err) {
		return err
	}
	return nil
}
