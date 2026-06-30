package query

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/llttlltt/dj-library-tools/internal/core/models"
)

// ResolvePath takes a track and a path-based field (e.g., beatgrids/bpm-drift)
// and returns the calculated value.
func ResolvePath(track models.Track, path string) (string, bool) {
	collectionPart, index, property, stat := ParsePath(path)
	if collectionPart == "" {
		return "", false
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
		if fn, ok := GetStat(stat); ok {
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

// ParsePath decomposes a path into its constituent parts.
func ParsePath(path string) (collection string, index int, property string, stat string) {
	index = -1
	cleanPath := path
	if idx := strings.LastIndex(path, "-"); idx != -1 {
		stat = path[idx+1:]
		cleanPath = path[:idx]
	}

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

	if idx := strings.Index(cleanPath, "/"); idx != -1 {
		property = cleanPath[idx+1:]
		if index == -1 {
			collectionPart = cleanPath[:idx]
		}
	}

	return strings.ToLower(collectionPart), index, property, stat
}

// ValidatePath checks if a path string is semantically valid according to models.
func ValidatePath(path string) error {
	collection, index, property, stat := ParsePath(path)

	// 1. Validate Collection
	fields, ok := models.CollectionFields[collection]
	if !ok {
		return fmt.Errorf("unrecognized collection %q", collection)
	}

	// 2. Validate Property (if present)
	if property != "" {
		if _, ok := fields[property]; !ok {
			return fmt.Errorf("collection %q does not have property %q", collection, property)
		}
	}

	// 3. Validate Stat (if present)
	if stat != "" {
		if stat == "density" {
			return nil
		}
		if _, ok := GetStat(stat); !ok {
			return fmt.Errorf("unrecognized stat %q", stat)
		}
	}

	// 4. Validate Indexing
	if index != -1 && (stat != "" || property == "") {
		// Indexing is allowed with property (hotcues.1/color)
		// but not with stats (beatgrids.1-count is nonsensical)
		// and we need a property if indexed (hotcues.1 is ambiguous but we allow it as "name")
	}

	return nil
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
