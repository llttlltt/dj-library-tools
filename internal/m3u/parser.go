package m3u

import (
	"bufio"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/llttlltt/dj-library-tools/internal/models"
)

// ReadM3U8 reads an M3U8 file and returns a slice of Tracks.
func ReadM3U8(path string) ([]models.Track, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	return ParseM3U8(f, filepath.Dir(path))
}

// ParseM3U8 parses an M3U8 stream.
func ParseM3U8(r io.Reader, baseDir string) ([]models.Track, error) {
	scanner := bufio.NewScanner(r)
	var lastDuration int
	var lastDisplay string
	var tracks []models.Track

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#EXTM3U") {
			continue
		}

		if strings.HasPrefix(line, "#EXTINF:") {
			info := strings.TrimPrefix(line, "#EXTINF:")
			commaIdx := strings.Index(info, ",")
			if commaIdx != -1 {
				durStr := info[:commaIdx]
				if d, err := strconv.Atoi(durStr); err == nil {
					lastDuration = d
				}
				lastDisplay = strings.TrimSpace(info[commaIdx+1:])
			}
			continue
		}

		trackPath := line
		if !filepath.IsAbs(trackPath) && baseDir != "" {
			trackPath = filepath.Join(baseDir, trackPath)
		}

		tracks = append(tracks, models.Track{
			ID:       trackPath,
			Display:  lastDisplay,
			Duration: lastDuration,
			Location: trackPath,
		})
		
		lastDuration = 0
		lastDisplay = ""
	}

	return tracks, scanner.Err()
}
