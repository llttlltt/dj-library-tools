package sync

import (
	"github.com/llttlltt/dj-library-tools/internal/models"
)

// MatchResult holds the outcome of matching a source track against the destination collection.
type MatchResult struct {
	Track          models.Track
	TargetTrack    *models.Track
	Confidence     float64
}

type Matcher struct {
	collection []models.Track
}

func NewMatcher(collection []models.Track) *Matcher {
	return &Matcher{collection: collection}
}

// Match finds the best target collection entry for the given neutral track.
func (m *Matcher) Match(t models.Track) MatchResult {
	for _, rt := range m.collection {
		if rt.Title == t.Title && rt.Artist == t.Artist {
			return MatchResult{Track: t, TargetTrack: &rt, Confidence: 1.0}
		}
	}
	return MatchResult{Track: t, TargetTrack: nil, Confidence: 0}
}
