package m3u

import (
	"path/filepath"
	"strings"
)

// FormatPath takes a file path and returns it with the specified new extension.
func FormatPath(path string, newExtension string) string {
	if newExtension == "" {
		return path
	}

	// Ensure extension has a dot prefix
	if !strings.HasPrefix(newExtension, ".") {
		newExtension = "." + newExtension
	}

	// Remove existing extension and append the new one
	base := strings.TrimSuffix(path, filepath.Ext(path))
	return base + newExtension
}

// IsM3UHeader checks if a line is an M3U header.
func IsM3UHeader(line string) bool {
	return strings.HasPrefix(strings.TrimSpace(line), "#EXTM3U")
}

// IsExtInfLine checks if a line is an #EXTINF entry.
func IsExtInfLine(line string) bool {
	return strings.HasPrefix(strings.TrimSpace(line), "#EXTINF")
}
