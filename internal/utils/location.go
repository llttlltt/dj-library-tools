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

// ParseLocation parses a provider/resource string and an optional query.
func ParseLocation(locStr string, query string) Location {
	loc := Location{
		Query: query,
	}

	// For file-based providers (m3u/m3u8), the path might contain slashes.
	// We want to detect the provider prefix but keep the rest as the resource path.
	if strings.HasPrefix(locStr, "m3u/") || strings.HasPrefix(locStr, "m3u8/") {
		parts := strings.SplitN(locStr, "/", 2)
		loc.Provider = parts[0]
		loc.Resource = parts[1]
	} else {
		// If query is empty, check if locStr contains a space-separated query
		if loc.Query == "" && strings.Contains(locStr, " ") {
			parts := strings.SplitN(locStr, " ", 2)
			locStr = parts[0]
			loc.Query = parts[1]
		}

		// Split provider/resource
		parts := strings.SplitN(locStr, "/", 2)
		loc.Provider = parts[0]
		if len(parts) > 1 {
			loc.Resource = parts[1]
		}
	}

	return loc
}
