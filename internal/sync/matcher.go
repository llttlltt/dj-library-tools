package sync

import (
	"github.com/llttlltt/dj-library-tools/internal/models"
	"github.com/llttlltt/dj-library-tools/internal/plex"
)

type MatchResult struct {
	PlexTrack  plex.Track
	RBTrack    *models.Track
	Confidence float64
}

type Matcher struct {
	collection []models.Track
}

func NewMatcher(collection []models.Track) *Matcher {
	return &Matcher{collection: collection}
}

func (m *Matcher) Match(t plex.Track) MatchResult {
	// Simple matching logic for now
	for _, rt := range m.collection {
		if rt.Title == t.Title && rt.Artist == t.Artist {
			return MatchResult{PlexTrack: t, RBTrack: &rt, Confidence: 1.0}
		}
	}
	return MatchResult{PlexTrack: t, RBTrack: nil, Confidence: 0}
}
