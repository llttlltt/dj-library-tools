package query

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/llttlltt/dj-library-tools/internal/models"
)

// ResolvePath takes a track and a path-based field (e.g., beatgrids/bpm-drift)
// and returns the calculated value.
func ResolvePath(track models.Track, path string) (string, bool) {
	// 1. Identify Stat
	stat := ""
	cleanPath := path
	if idx := strings.LastIndex(path, "-"); idx != -1 {
		stat = path[idx+1:]
		cleanPath = path[:idx]
	}

	// 2. Identify Index
	index := -1
	collectionPart := cleanPath
	if idx := strings.Index(cleanPath, "."); idx != -1 {
		indexStr := ""
		if slashIdx := strings.Index(cleanPath, "/"); slashIdx != -1 {
			indexStr = cleanPath[idx+1 : slashIdx]
			collectionPart = cleanPath[:idx]
		} else {
			indexStr = cleanPath[idx+1:]
			collectionPart = cleanPath[:idx]
		}
		
		val, err := strconv.Atoi(indexStr)
		if err == nil {
			index = val - 1 // 1-based to 0-based
		}
	}

	// 3. Identify Property
	property := ""
	if idx := strings.Index(cleanPath, "/"); idx != -1 {
		property = cleanPath[idx+1:]
		if index == -1 {
			collectionPart = cleanPath[:idx]
		}
	}

	// 4. Resolve Collection
	var values []string
	collectionName := strings.ToLower(collectionPart)

	switch collectionName {
	case "hotcues":
		for _, cp := range track.CuePoints {
			if cp.Type == models.CueTypeHot {
				values = append(values, resolveCueProperty(cp, property))
			}
		}
	case "memorycues":
		for _, cp := range track.CuePoints {
			if cp.Type == models.CueTypeMemory {
				values = append(values, resolveCueProperty(cp, property))
			}
		}
	case "beatgrids":
		for _, tm := range track.TempoMarkers {
			values = append(values, resolveMarkerProperty(tm, property))
		}
	case "playlists":
		for _, pm := range track.Playlists {
			values = append(values, resolvePlaylistProperty(pm, property))
		}
	default:
		return "", false
	}

	// 5. Handle Density Special Case
	if stat == "density" {
		return DensityStat(len(values), track.Duration), true
	}

	// 6. Apply Indexing
	if index != -1 {
		if index >= 0 && index < len(values) {
			return values[index], true
		}
		return "", true // Index out of bounds, but valid path
	}

	// 7. Apply Stat
	if stat != "" {
		if fn, ok := Stats[stat]; ok {
			return fn(values), true
		}
	}

	// 8. Default: Return first item or concatenated string?
	// For "any" matching, we usually return a comma-separated list for OpSubstring
	if len(values) > 0 {
		return strings.Join(values, ","), true
	}

	return "", true
}

func resolveCueProperty(cp models.CuePoint, prop string) string {
	switch strings.ToLower(prop) {
	case "color":
		return cp.Color
	case "name":
		return cp.Name
	case "position":
		return fmt.Sprintf("%.4f", cp.Position)
	default:
		return cp.Name // Default to name
	}
}

func resolveMarkerProperty(tm models.TempoMarker, prop string) string {
	switch strings.ToLower(prop) {
	case "bpm":
		return fmt.Sprintf("%.4f", tm.BPM)
	case "position":
		return fmt.Sprintf("%.4f", tm.Position)
	default:
		return fmt.Sprintf("%.4f", tm.BPM)
	}
}

func resolvePlaylistProperty(pm models.PlaylistMembership, prop string) string {
	switch strings.ToLower(prop) {
	case "name":
		return pm.Name
	case "folder":
		return pm.Folder
	default:
		return pm.Name
	}
}
