package location

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
// It is strictly syntactic and does not contain provider-specific logic.
func ParseLocation(locStr string, query string) Location {
	loc := Location{
		Query: query,
	}

	// 1. Split by query if not already provided
	if loc.Query == "" && strings.Contains(locStr, " ") {
		parts := strings.SplitN(locStr, " ", 2)
		locStr = parts[0]
		loc.Query = parts[1]
	}

	// 2. Resolve Provider and Resource via / or : separator
	sepIdx := strings.IndexAny(locStr, "/:")
	if sepIdx != -1 {
		loc.Provider = locStr[:sepIdx]
		loc.Resource = locStr[sepIdx+1:]
	} else {
		loc.Provider = locStr
	}

	return loc
}
