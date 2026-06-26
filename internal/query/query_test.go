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
		TrackID:      12345,
		Name:         "Sixteen Oceans",
		Artist:       "Four Tet",
		Composer:     "Kieran Hebden",
		Album:        "Sixteen Oceans",
		Grouping:     "Electronic",
		Genre:        "House",
		Kind:         "MP3 File",
		Size:         10485760,
		TotalTime:    300,
		DiscNumber:   1,
		TrackNumber:  5,
		Year:         2020,
		AverageBpm:   "124.0",
		DateAdded:    "2020-03-13",
		DateModified: "2020-03-14",
		LastPlayed:   "2020-06-25",
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
	}

	tests := []struct {
		name    string
		query   string
		matches bool
	}{
		// Standard Spec Fields
		{"ID match", "id:12345", true},
		{"Title match", "title:Oceans", true},
		{"Artist match", "artist:Four", true},
		{"Composer match", "composer:Kieran", true},
		{"Album match", "album:Sixteen", true},
		{"Grouping match", "grouping:Electronic", true},
		{"Genre match", "genre:House", true},
		{"Kind match", "kind:MP3", true},
		{"Size match", "size:10485760", true},
		{"Time match", "time:300", true},
		{"Disc match", "disc:1", true},
		{"Track match", "track:5", true},
		{"Year match", "year:2020", true},
		{"BPM match", "bpm:124", true},
		{"Added match", "added:2020-03-13", true},
		{"Modified match", "modified:2020-03-14", true},
		{"Played match", "played:2020-06-25", true},
		{"Bitrate match", "bitrate:320", true},
		{"Samplerate match", "samplerate:44100", true},
		{"Comment match", "comment:Great", true},
		{"Playcount match", "playcount:10", true},
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
		{"Default field (title)", "Oceans", true},

		// Cues & Tempos
		{"Tempo count", "beatgrids:2", true},
		{"Hot cue count", "hotcues:2", true},
		{"Memory cue count", "memorycues:2", true},
		{"Hot Cue Slot A", "hotcue:a", true},
		{"Hot Cue A BrightGreen", "hotcue:a:brightgreen", true},
		{"Hot Cue B Purple", "hotcue:b:purple", true},
		{"Memory Cue 1", "memorycue:1", true},
		{"Memory Cue 1 Comment", "memorycue:1:comment:GROOVE", true},
		{"Memory Cue 2 No Comment", "memorycue:2:comment:none", true},
		{"Memory Cue Loop", "memorycues:loop", true},
		{"Tempo 1 Bpm", "tempo:1:bpm:124", true},
		{"Tempo 1 Meter", "tempo:1:meter:4/4", true},
		{"Tempo 1 Inizio", "tempo:1:inizio:0.0", true},
		{"Tempo 1 Battito", "tempo:1:battito:1", true},
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
		Name:    "Summer Vibes",
		Type:    1,
		Entries: rekordbox.PtrInt32(12),
	}
	wintersFolder := rekordbox.Node{
		Name: "Winter Sets",
		Type: 0,
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
		{"parent match", "parent:My Sets", summerVibes, "My Sets", true},
		{"parent no match", "parent:Other", summerVibes, "My Sets", false},
		{"entries range", "entries:10..15", summerVibes, "My Sets", true},
		{"entries exact", "entries:12", summerVibes, "My Sets", true},
		{"entries no match", "entries:5", summerVibes, "My Sets", false},
		{"type playlist", "type:1", summerVibes, "My Sets", true},
		{"type folder match", "type:0", wintersFolder, "", true},
		{"type folder no match", "type:1", wintersFolder, "", false},
		{"boolean AND", "name:Summer && parent:My Sets", summerVibes, "My Sets", true},
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
