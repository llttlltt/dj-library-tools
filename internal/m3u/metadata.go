package m3u

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/dhowden/tag"
)

// AudioMetadata holds common metadata extracted from audio files.
type AudioMetadata struct {
	Artist   string
	Title    string
	Album    string
	Duration float64
}

// ExtractMetadata reads audio metadata from the given file path.
func ExtractMetadata(path string) (AudioMetadata, error) {
	f, err := os.Open(path)
	if err != nil {
		return AudioMetadata{}, err
	}
	defer f.Close()

	meta := AudioMetadata{
		Duration: -1, // Standard M3U "unknown" duration
	}

	m, err := tag.ReadFrom(f)
	if err == nil {
		meta.Artist = m.Artist()
		meta.Title = m.Title()
		meta.Album = m.Album()
	}

	// Fallback logic: If Artist or Title is missing, parse from filename
	// Matches legacy Bash script behavior
	if meta.Artist == "" || meta.Title == "" {
		filename := filepath.Base(path)
		filenameNoExt := strings.TrimSuffix(filename, filepath.Ext(filename))

		if strings.Contains(filenameNoExt, " - ") {
			parts := strings.SplitN(filenameNoExt, " - ", 2)
			if meta.Artist == "" {
				meta.Artist = strings.TrimSpace(parts[0])
			}
			if meta.Title == "" {
				meta.Title = strings.TrimSpace(parts[1])
			}
		} else if meta.Title == "" {
			meta.Title = filenameNoExt
		}
	}

	return meta, nil
}
