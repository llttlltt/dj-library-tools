package sync

import (
	"path/filepath"
	"strings"

	"github.com/llttlltt/dj-library-tools/internal/core/models"
)

// MatchResult holds the outcome of matching a source track against the destination collection.
type MatchResult struct {
	Track       models.Track
	TargetTrack *models.Track
	Confidence  float64
}

type Matcher struct {
	collection []models.Track
	keys       []string
}

func NewMatcher(collection []models.Track) *Matcher {
	return &Matcher{
		collection: collection,
		keys:       []string{"artist", "title"}, // Default match criteria
	}
}

func (m *Matcher) WithKeys(keys []string) *Matcher {
	m.keys = keys
	return m
}

// Match finds the best target collection entry for the given neutral track using the configured keys.
func (m *Matcher) Match(t models.Track) MatchResult {
	for _, rt := range m.collection {
		if m.matchTrack(t, rt) {
			return MatchResult{Track: t, TargetTrack: &rt, Confidence: 1.0}
		}
	}
	return MatchResult{Track: t, TargetTrack: nil, Confidence: 0}
}

func (m *Matcher) matchTrack(s, t models.Track) bool {
	for _, key := range m.keys {
		switch strings.ToLower(key) {
		case "artist":
			if !strings.EqualFold(s.Artist, t.Artist) {
				return false
			}
		case "title":
			if !strings.EqualFold(s.Title, t.Title) {
				return false
			}
		case "album":
			if !strings.EqualFold(s.Album, t.Album) {
				return false
			}
		case "filename":
			if filepath.Base(s.Location) != filepath.Base(t.Location) {
				return false
			}
		case "path":
			if s.Location != t.Location {
				return false
			}
		}
	}
	return true
}
