package query

import (
	"fmt"
	"testing"

	"github.com/llttlltt/dj-library-tools/pkg/rekordbox"
)

// TestPlaylistCountMatching verifies that playlistcount uses exact numeric equality,
// not substring matching, to prevent false positives.
func TestPlaylistCountMatching(t *testing.T) {
	track := rekordbox.Track{Name: "Test Track"}
	parser := NewParser()

	tests := []struct {
		name      string
		query     string
		playlists []string
		want      bool
	}{
		{"exact zero match", "playlistcount:0", nil, true},
		{"exact zero no false positive from 10", "playlistcount:0", makePlaylists(10), false},
		{"exact two match", "playlistcount:2", makePlaylists(2), true},
		{"exact two mismatch", "playlistcount:2", makePlaylists(3), false},
		{"gt operator", "playlistcount:>3", makePlaylists(4), true},
		{"gt operator no match", "playlistcount:>3", makePlaylists(3), false},
		{"gte operator", "playlistcount:>=3", makePlaylists(3), true},
		{"range operator", "playlistcount:2..4", makePlaylists(3), true},
		{"range no match", "playlistcount:2..4", makePlaylists(1), false},
		{"bpm exact no substring", "bpm:12", nil, false}, // bpm is 0 for empty track
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			q := parser.Parse(tt.query)
			eval := NewEvaluator(q)
			got := eval.MatchesWithPlaylists(track, tt.playlists)
			if got != tt.want {
				t.Errorf("query %q playlists=%d: got %v, want %v", tt.query, len(tt.playlists), got, tt.want)
			}
		})
	}
}

func makePlaylists(n int) []string {
	playlists := make([]string, n)
	for i := range playlists {
		playlists[i] = fmt.Sprintf("Playlist%d", i+1)
	}
	return playlists
}

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

func TestEvaluatorMatchesNode(t *testing.T) {
	summerVibes := rekordbox.Node{
		BaseNode: rekordbox.BaseNode{Type: 1, Name: "Summer Vibes"},
		Entries:  12,
	}
	wintersFolder := rekordbox.Node{
		BaseNode: rekordbox.BaseNode{Type: 0, Name: "Winter Sets"},
	}

	parser := NewParser()

	tests := []struct {
		name         string
		query        string
		node         rekordbox.Node
		parentFolder string
		want         bool
	}{
		{"name substring match", "name:Summer", summerVibes, "My Sets", true},
		{"name no match", "name:Winter", summerVibes, "My Sets", false},
		{"folder match", "folder:My Sets", summerVibes, "My Sets", true},
		{"folder no match", "folder:Other", summerVibes, "My Sets", false},
		{"root level folder empty string", "folder:", summerVibes, "", true},
		{"entries range", "entries:10..15", summerVibes, "My Sets", true},
		{"entries exact", "entries:12", summerVibes, "My Sets", true},
		{"entries no match", "entries:5", summerVibes, "My Sets", false},
		{"type playlist", "type:1", summerVibes, "My Sets", true},
		{"type folder match", "type:0", wintersFolder, "", true},
		{"type folder no match", "type:1", wintersFolder, "", false},
		{"boolean AND", "name:Summer && folder:My Sets", summerVibes, "My Sets", true},
		{"boolean AND fail", "name:Summer && folder:Other", summerVibes, "My Sets", false},
		{"empty query matches all", "", summerVibes, "My Sets", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			q := parser.Parse(tt.query)
			eval := NewEvaluator(q)
			got := eval.MatchesNode(tt.node, tt.parentFolder)
			if got != tt.want {
				t.Errorf("query %q: got %v, want %v", tt.query, got, tt.want)
			}
		})
	}
}
