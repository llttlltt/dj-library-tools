package media

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type Config struct {
	Dest                  string            `json:"dest"`
	Format                string            `json:"format"`
	ID3v23                bool              `json:"id3v23"`
	WriteMetadata         bool              `json:"write_metadata"`
	Formats               map[string]string `json:"formats"`
	MaxBitrate            int               `json:"max_bitrate"`
	Embed                 bool              `json:"embed"`
	Paths                 map[string]string `json:"paths"`
	NeverConvertLossy     bool              `json:"never_convert_lossy"`
	CopyAlbumArt          bool              `json:"copy_album_art"`
	AlbumArtMaxWidth      int               `json:"album_art_maxwidth"`
}

func DefaultConfig() *Config {
	return &Config{
		Format:            "mp3",
		ID3v23:            true,
		WriteMetadata:     true,
		MaxBitrate:        320,
		Embed:             true,
		NeverConvertLossy: false,
		CopyAlbumArt:      true,
		AlbumArtMaxWidth:  1000,
		Formats: map[string]string{
			"mp3": "ffmpeg -i \"$source\" -y -vn -b:a 320k -ar 44100 -sample_fmt s16p -map_metadata 0 -id3v2_version 3 \"$dest\"",
		},
		Paths: map[string]string{
			"default":   "{{.Artist}} - {{.Album}} - {{.Title}}",
			"singleton": "{{.Artist}} - {{.Album}} - {{.Title}}",
			"comp":      "{{.Album}} - {{.Artist}} - {{.Title}}",
		},
	}
}

func LoadConfig(path string) (*Config, error) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return DefaultConfig(), nil
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

func (c *Config) Save(path string) error {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}
