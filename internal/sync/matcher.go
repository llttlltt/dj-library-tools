package sync

import (
	"strings"

	"github.com/llttlltt/dj-library-tools/internal/plex"
	"github.com/llttlltt/dj-library-tools/pkg/rekordbox"
)

// MatchResult represents the outcome of a matching operation.
type MatchResult struct {
	PlexTrack  plex.Track
	RBTrack    *rekordbox.Track
	Confidence float64 // 0.0 to 1.0
}

// Matcher handles matching Plex tracks to Rekordbox tracks.
type Matcher struct {
	Collection []rekordbox.Track
}

// NewMatcher creates a new Matcher with a Rekordbox collection.
func NewMatcher(collection []rekordbox.Track) *Matcher {
	return &Matcher{Collection: collection}
}

// Match attempts to find the best Rekordbox track for a given Plex track.
func (m *Matcher) Match(plexTrack plex.Track) MatchResult {
	var bestMatch *rekordbox.Track
	maxConfidence := 0.0

	plexTitle := strings.ToLower(normalizeString(plexTrack.Title))
	plexArtist := strings.ToLower(normalizeString(plexTrack.Artist))

	for i := range m.Collection {
		rbTrack := &m.Collection[i]
		rbTitle := strings.ToLower(normalizeString(rbTrack.Name))
		rbArtist := strings.ToLower(normalizeString(rbTrack.Artist))

		confidence := 0.0

		// Exact match
		if plexTitle == rbTitle && plexArtist == rbArtist {
			confidence = 1.0
		} else if plexTitle == rbTitle {
			// Title match, check artist
			if strings.Contains(rbArtist, plexArtist) || strings.Contains(plexArtist, rbArtist) {
				confidence = 0.9
			} else {
				confidence = 0.5
			}
		} else if strings.Contains(rbTitle, plexTitle) && (plexArtist == rbArtist) {
			confidence = 0.8
		}

		if confidence > maxConfidence {
			maxConfidence = confidence
			bestMatch = rbTrack
		}

		if maxConfidence == 1.0 {
			break
		}
	}

	return MatchResult{
		PlexTrack:  plexTrack,
		RBTrack:    bestMatch,
		Confidence: maxConfidence,
	}
}

func normalizeString(s string) string {
	s = strings.ReplaceAll(s, " (", " ")
	s = strings.ReplaceAll(s, "(", " ")
	s = strings.ReplaceAll(s, ")", " ")
	s = strings.ReplaceAll(s, " [", " ")
	s = strings.ReplaceAll(s, "[", " ")
	s = strings.ReplaceAll(s, "]", " ")
	s = strings.ReplaceAll(s, " - ", " ")
	return strings.TrimSpace(s)
}
