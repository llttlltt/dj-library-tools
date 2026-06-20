package utils

import (
	"strings"
)

// Location represents a parsed provider/resource:query string.
type Location struct {
	Provider string
	Resource string
	Query    string
}

// ParseLocation parses a string in the format provider/resource:query.
// Defaults:
// - If no resource is specified (e.g., "plex"), it defaults based on provider.
// - If no query is specified, it's empty.
func ParseLocation(input string) Location {
	loc := Location{}

	// Split query part
	parts := strings.SplitN(input, ":", 2)
	base := parts[0]
	if len(parts) > 1 {
		loc.Query = parts[1]
	}

	// Split provider/resource
	baseParts := strings.SplitN(base, "/", 2)
	loc.Provider = baseParts[0]
	if len(baseParts) > 1 {
		loc.Resource = baseParts[1]
	}

	// Apply defaults
	if loc.Resource == "" {
		switch loc.Provider {
		case "plex":
			loc.Resource = "playlists"
		case "rb", "rekordbox":
			loc.Resource = "tracks"
		}
	}

	return loc
}
