package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type AppConfig struct {
	Plex      PlexConfig        `json:"plex"`
	Rekordbox RekordboxConfig   `json:"rekordbox"`
	Updates   UpdateConfig      `json:"updates"`
	PathMaps  map[string]string `json:"path_maps"`
}

type UpdateConfig struct {
	LastCheckAt       string `json:"last_check_at"`
	CheckIntervalHour int    `json:"check_interval_hour"`
}

type PlexConfig struct {
	Token string `json:"token"`
	Host  string `json:"host"`
	Port  int    `json:"port"`
}

type RekordboxConfig struct {
	PrimaryFilePath string `json:"primary_file_path"`
}

func GetConnectionsDir() (string, error) {
	dir, err := GetConfigDir()
	if err != nil {
		return "", err
	}
	p := filepath.Join(dir, "connections")
	return p, os.MkdirAll(p, 0755)
}

func GetWorkflowsDir() (string, error) {
	dir, err := GetConfigDir()
	if err != nil {
		return "", err
	}
	p := filepath.Join(dir, "workflows")
	return p, os.MkdirAll(p, 0755)
}

func GetPathMapsDir() (string, error) {
	dir, err := GetConfigDir()
	if err != nil {
		return "", err
	}
	p := filepath.Join(dir, "path-maps")
	return p, os.MkdirAll(p, 0755)
}

func GetConfigDir() (string, error) {
	// Respect XDG_CONFIG_HOME if set, otherwise fallback to ~/.config
	dir := os.Getenv("XDG_CONFIG_HOME")
	if dir == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		dir = filepath.Join(home, ".config")
	}
	return filepath.Join(dir, "djlt"), nil
}

func LoadAppConfig() (*AppConfig, error) {
	// ... (implementation same but ensures maps are initialized)
	dir, err := GetConfigDir()
	if err != nil {
		return nil, err
	}
	path := filepath.Join(dir, "config.json")

	cfg := &AppConfig{
		PathMaps: make(map[string]string),
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		return cfg, nil
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(data, cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}

func SaveAppConfig(cfg *AppConfig) error {
	dir, err := GetConfigDir()
	if err != nil {
		return err
	}
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	path := filepath.Join(dir, "config.json")
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}
