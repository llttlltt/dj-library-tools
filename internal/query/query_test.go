package query

import (
	"testing"

	"github.com/llttlltt/dj-library-tools/pkg/rekordbox"
)

func TestQueryEvaluator(t *testing.T) {
	track := rekordbox.Track{
		Name:       "Sixteen Oceans",
		Artist:     "Four Tet",
		AverageBpm: 124.0,
		Tempo: []rekordbox.Tempo{
			{Bpm: 124.0},
			{Bpm: 125.0},
		},
		Genre: "Electronic",
		PositionMark: []rekordbox.PositionMark{
			{Num: 0, Red: 40, Green: 226, Blue: 20}, // Green Hot Cue A
			{Num: 1, Red: 224, Green: 100, Blue: 27}, // Orange Hot Cue B
			{Num: -1, Name: "GROOVE", Type: 4, Start: 60.0}, // Memory Cue 1 (Highest Start)
			{Num: -1, Name: "", Start: 30.0},                // Memory Cue 2
		},
	}

	tests := []struct {
		name    string
		query   string
		matches bool
	}{
		{"Substring match", "artist:Four", true},
		{"Exact match", "artist=\"Four Tet\"", true},
		{"Exact match no quotes", "artist=Four Tet", true},
		{"Exact mismatch", "artist=Four", false},
		{"Regex match", "artist::^Four", true},
		{"Range match", "bpm:120..130", true},
		{"Range mismatch", "bpm:130..140", false},
		{"Multi-criteria", "artist:Four bpm:120..130", true},
		{"Default field (name)", "Oceans", true},
		{"Negation", "!artist:Skrillex", true},
		{"Tempo count", "tempos:2", true},
		{"Hot cue count", "hotcues:2", true},
		{"Memory cue count", "memorycues:2", true},
		{"Any Green Hot Cue", "hotcues:green", true},
		{"Any Orange Hot Cue", "hotcues:orange", true},
		{"Any Green Memory Cue", "memorycues:green", false},
		{"Hot Cue Slot A", "hotcue:a", true},
		{"Hot Cue Slot B", "hotcue:b", true},
		{"Hot Cue Slot C", "hotcue:c", false},
		{"Hot Cue A Green", "hotcue:a:green", true},
		{"Hot Cue B Green", "hotcue:b:green", false},
		{"Hot Cue B Orange", "hotcue:b:orange", true},
		{"Memory Cue 1", "memorycue:1", true},
		{"Memory Cue 1 Label", "memorycue:1:label:GROOVE", true},
		{"Memory Cue 2 No Label", "memorycue:2:label:\"\"", true},
		{"Memory Cue 2 Empty Alias", "memorycue:2:label:empty", true},
		{"Memory Cue Loop", "memorycues:loop", true},
		{"Any Memory Cue No Label", "memorycues:label:\"\"", true},
		{"Property Chaining", "hotcue:a:green:label:none", true},
		{"Numeric GT", "bpm:>120", true},
		{"Numeric GTE", "bpm:>=124", true},
		{"Numeric LT", "bpm:<125", true},
		{"Numeric LTE", "bpm:<=124", true},
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
