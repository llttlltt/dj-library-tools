package query

import (
	"fmt"
	"testing"

	"github.com/llttlltt/dj-library-tools/internal/models"
	"github.com/llttlltt/dj-library-tools/internal/rekordbox"
)

// TestPlaylistMembershipMatching verifies that 'playlists' handles both
// numeric count matching and string name matching.
func TestPlaylistMembershipMatching(t *testing.T) {
	track := rekordbox.Track{Name: "Test Track"}.ToNeutral()
	parser := NewParser()

	tests := []struct {
		name      string
		query     string
		playlists []string
		want      bool
	}{
		// Numeric (Count) logic
		{"exact zero match", "playlists:0", nil, true},
		{"exact zero no false positive from 10", "playlists:0", makePlaylists(10), false},
		{"exact two match", "playlists:2", makePlaylists(2), true},
		{"exact two mismatch", "playlists:2", makePlaylists(3), false},
		{"gt operator", "playlists:>3", makePlaylists(4), true},
		{"gt operator no match", "playlists:>3", makePlaylists(3), false},
		{"gte operator", "playlists:>=3", makePlaylists(3), true},
		{"range operator", "playlists:2..4", makePlaylists(3), true},
		{"range no match", "playlists:2..4", makePlaylists(1), false},

		// String (Name) logic
		{"name substring match", "playlists:Summer", []string{"Summer Vibes"}, true},
		{"name exact match (quoted number)", "playlists:'101'", []string{"101"}, true},
		{"name exact mismatch (quoted number)", "playlists:'101'", []string{"102"}, false},
		{"name no match", "playlists:Winter", []string{"Summer Vibes"}, false},
		{"mixed case match", "playlists:summer", []string{"Summer Vibes"}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			q := parser.Parse(tt.query)
			eval := NewEvaluator(q)
			got := eval.MatchesWithPlaylists(track, tt.playlists)
			if got != tt.want {
				t.Errorf("query %q playlists=%v: got %v, want %v", tt.query, tt.playlists, got, tt.want)
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
		TrackID:      12345,
		Name:         "Sixteen Oceans",
		Artist:       "Four Tet",
		Album:        "Sixteen Oceans",
		Genre:        "House",
		Size:         10485760,
		Year:         2020,
		AverageBpm:   "124.0",
		DateAdded:    "2020-03-13",
		DateModified: "2020-03-14",
		BitRate:      320,
		SampleRate:   44100.0,
		Comments:     "Great track",
		PlayCount:    10,
		Rating:       255, // 5 stars
		Location:     "file://localhost/Users/test/track.mp3",
		Remixer:      "Four Tet",
		Tonality:     "8A",
		Label:        "Text Records",
		Mix:          "Original Mix",
		Colour:       "0xFF007F", // pink
		Tempo: []rekordbox.Tempo{
			{Bpm: "124.0", Metro: "4/4", Inizio: "0.0", Battito: 1},
			{Bpm: "125.0"},
		},
		PositionMark: []rekordbox.PositionMark{
			{Num: 0, Red: 16, Green: 177, Blue: 118},         // brightgreen
			{Num: 1, Red: 180, Green: 50, Blue: 255},         // purple
			{Num: -1, Name: "GROOVE", Type: 4, Start: "60.0"}, // Memory Cue 1
			{Num: -1, Name: "", Start: "30.0"},                // Memory Cue 2
		},
	}.ToNeutral()

	tests := []struct {
		name    string
		query   string
		matches bool
	}{
		// Standard Spec Fields
		{"ID match", "id:12345", true},
		{"Title match", "title:Oceans", true},
		{"Artist match", "artist:Four", true},
		{"Album match", "album:Sixteen", true},
		{"Genre match", "genre:House", true},
		{"Size match", "size:10485760", true},
		{"Year match", "year:2020", true},
		{"BPM match", "bpm:124", true},
		{"Added match", "added:2020-03-13", true},
		{"Modified match", "modified:2020-03-14", true},
		{"Bitrate match", "bitrate:320", true},
		{"Samplerate match", "samplerate:44100", true},
		{"Comment match", "comment:Great", true},
		{"Plays match", "plays:10", true},
		{"Rating match (stars)", "rating:5", true},
		{"Location match", "location:track.mp3", true},
		{"Remixer match", "remixer:Four", true},
		{"Key match", "key:8A", true},
		{"Label match", "label:Text", true},
		{"Mix match", "mix:Original", true},
		{"Track Color", "color:pink", true},

		// Logical & Operators
		{"Substring match", "artist:Four", true},
		{"Exact match", "artist=\"Four Tet\"", true},
		{"Regex match", "artist::^Four", true},
		{"Range match", "bpm:120..130", true},
		{"Negation", "!genre:Techno", true},

		// Cues & Tempos
		{"Tempo count", "beatgrids:2", true},
		{"Hot cue count", "hotcues:2", true},
		{"Memory cue count", "memorycues:2", true},
	}

	parser := NewParser()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			q := parser.Parse(tt.query)
			eval := NewEvaluatorWithMatcher(q, MockMatcher{})
			if eval.Matches(track) != tt.matches {
				t.Errorf("Query '%s' match expected %v, got %v", tt.query, tt.matches, !tt.matches)
			}
		})
	}
}

func TestEvaluatorMatchesNode(t *testing.T) {
	summerVibes := rekordbox.Node{
		Name:    "Summer Vibes",
		Type:    1,
		Entries: rekordbox.PtrInt32(12),
	}.ToNeutral("My Sets")
	wintersFolder := rekordbox.Node{
		Name: "Winter Sets",
		Type: 0,
	}.ToNeutral("")

	parser := NewParser()

	tests := []struct {
		name  string
		query string
		node  models.ResourceGroup
		want  bool
	}{
		{"name substring match", "name:Summer", summerVibes, true},
		{"name no match", "name:Winter", summerVibes, false},
		{"parent match", "parent:'My Sets'", summerVibes, true},
		{"parent no match", "parent:Other", summerVibes, false},
		{"items range", "items:10..15", summerVibes, true},
		{"items exact", "items:12", summerVibes, true},
		{"items no match", "items:5", summerVibes, false},
		{"type playlist", "type:1", summerVibes, true},
		{"type folder match", "type:0", wintersFolder, true},
		{"type folder no match", "type:1", wintersFolder, false},
		{"boolean AND", "name:Summer && parent:'My Sets'", summerVibes, true},
		{"empty query matches all", "", summerVibes, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			q := parser.Parse(tt.query)
			eval := NewEvaluator(q)
			got := eval.MatchesNode(tt.node)
			if got != tt.want {
				t.Errorf("query %q: got %v, want %v", tt.query, got, tt.want)
			}
		})
	}
}

type MockMatcher struct{}
func (m MockMatcher) CustomMatch(track models.Track, field string, op Operator, value string) bool {
	if field == "beatgrids" || field == "hotcues" || field == "memorycues" {
		return true // Mock behavior for test
	}
	return false
}
func (m MockMatcher) GetTrackColorName(hex string) string {
	if hex == "0xFF007F" { return "pink" }
	return hex
}
