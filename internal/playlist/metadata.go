package playlist

import (
	"os"

	"github.com/dhowden/tag"
)

// AudioMetadata holds common metadata extracted from audio files.
type AudioMetadata struct {
	Artist string
	Title  string
	Album  string
}

// ExtractMetadata reads audio metadata from the given file path.
func ExtractMetadata(path string) (AudioMetadata, error) {
	f, err := os.Open(path)
	if err != nil {
		return AudioMetadata{}, err
	}
	defer f.Close()

	m, err := tag.ReadFrom(f)
	if err != nil {
		return AudioMetadata{}, err
	}

	return AudioMetadata{
		Artist: m.Artist(),
		Title:  m.Title(),
		Album:  m.Album(),
	}, nil
}
