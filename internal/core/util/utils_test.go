package util

import (
	"github.com/llttlltt/dj-library-tools/internal/core/location"
	"testing"
)

func TestParseLocation(t *testing.T) {
	tests := []struct {
		locStr   string
		query    string
		expected location.Location
	}{
		{
			locStr: "plex/playlists",
			query:  "Summer",
			expected: location.Location{
				Provider: "plex",
				Resource: "playlists",
				Query:    "Summer",
			},
		},
		{
			locStr: "rb/tracks",
			query:  "bpm:120..130",
			expected: location.Location{
				Provider: "rb",
				Resource: "tracks",
				Query:    "bpm:120..130",
			},
		},
		{
			locStr: "plex/playlists",
			query:  "",
			expected: location.Location{
				Provider: "plex",
				Resource: "playlists",
				Query:    "",
			},
		},
		{
			locStr: "m3u8/file",
			query:  "my_playlist.m3u8",
			expected: location.Location{
				Provider: "m3u8",
				Resource: "file",
				Query:    "my_playlist.m3u8",
			},
		},
	}

	for _, tt := range tests {
		result := location.ParseLocation(tt.locStr, tt.query)
		if result.Provider != tt.expected.Provider || result.Resource != tt.expected.Resource || result.Query != tt.expected.Query {
			t.Errorf("ParseLocation(%q, %q) = %+v; want %+v", tt.locStr, tt.query, result, tt.expected)
		}
	}
}
