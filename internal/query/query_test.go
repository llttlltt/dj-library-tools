package query

import (
	"testing"

	"github.com/llttlltt/dj-library-tools/pkg/rekordbox"
)

func TestQueryEvaluator(t *testing.T) {
	track := rekordbox.Track{
		Name:   "Sixteen Oceans",
		Artist: "Four Tet",
		Tempo: []rekordbox.Tempo{
			{Inizio: 124.0},
		},
		Genre: "Electronic",
	}

	tests := []struct {
		name    string
		query   string
		matches bool
	}{
		{"Substring match", "artist:Four", true},
		{"Exact match", "artist=\"Four Tet\"", true},
		{"Exact mismatch", "artist=Four", false},
		{"Regex match", "artist::^Four", true},
		{"Range match", "bpm:120..130", true},
		{"Range mismatch", "bpm:130..140", false},
		{"Multi-criteria", "artist:Four bpm:120..130", true},
		{"Default field (name)", "Oceans", true},
		{"Negation", "!artist:Skrillex", true},
	}

	parser := NewParser()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			q := parser.Parse(tt.query)
			eval := NewEvaluator(q)
			if eval.Matches(track) != tt.matches {
				t.Errorf("Query '%s' match expected %v, got %v", tt.query, tt.matches, !tt.matches)
			}
		})
	}
}
