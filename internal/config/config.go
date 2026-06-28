package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type AppConfig struct {
	PlexToken       string            `json:"plex_token"`
	PlexHost        string            `json:"plex_host"`
	PlexPort        int               `json:"plex_port"`
	PrimaryFilePath string            `json:"primary_file_path"`
	PathMaps        map[string]string `json:"path_maps"`
}

func GetConfigDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".config", "djlt"), nil
}

func LoadAppConfig() (*AppConfig, error) {
	dir, err := GetConfigDir()
	if err != nil {
		return nil, err
	}
	path := filepath.Join(dir, "config.json")

	if _, err := os.Stat(path); os.IsNotExist(err) {
		return &AppConfig{}, nil
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg AppConfig
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
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
