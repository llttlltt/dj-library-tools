package sync

import (
	"github.com/llttlltt/dj-library-tools/internal/models"
)

// MatchResult holds the outcome of matching a source track against the RB collection.
type MatchResult struct {
	Track      models.Track
	RBTrack    *models.Track
	Confidence float64
}

type Matcher struct {
	collection []models.Track
}

func NewMatcher(collection []models.Track) *Matcher {
	return &Matcher{collection: collection}
}

// Match finds the best RB collection entry for the given neutral track.
func (m *Matcher) Match(t models.Track) MatchResult {
	for _, rt := range m.collection {
		if rt.Title == t.Title && rt.Artist == t.Artist {
			return MatchResult{Track: t, RBTrack: &rt, Confidence: 1.0}
		}
	}
	return MatchResult{Track: t, RBTrack: nil, Confidence: 0}
}
