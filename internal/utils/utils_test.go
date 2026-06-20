package utils

import (
	"testing"
)

func TestParseLocation(t *testing.T) {
	tests := []struct {
		input    string
		expected Location
	}{
		{
			input: "plex/playlists:Summer",
			expected: Location{
				Provider: "plex",
				Resource: "playlists",
				Query:    "Summer",
			},
		},
		{
			input: "rb:bpm:120..130",
			expected: Location{
				Provider: "rb",
				Resource: "tracks",
				Query:    "bpm:120..130",
			},
		},
		{
			input: "plex",
			expected: Location{
				Provider: "plex",
				Resource: "playlists",
				Query:    "",
			},
		},
		{
			input: "m3u8:my_playlist.m3u8",
			expected: Location{
				Provider: "m3u8",
				Resource: "",
				Query:    "my_playlist.m3u8",
			},
		},
	}

	for _, tt := range tests {
		result := ParseLocation(tt.input)
		if result.Provider != tt.expected.Provider || result.Resource != tt.expected.Resource || result.Query != tt.expected.Query {
			t.Errorf("ParseLocation(%q) = %+v; want %+v", tt.input, result, tt.expected)
		}
	}
}

func TestExpandPath(t *testing.T) {
	// Note: ~ expansion depends on home dir, so we just check prefix logic
	path := "~/test.xml"
	expanded := ExpandPath(path)
	if expanded == path {
		t.Errorf("ExpandPath(%q) did not expand tilde", path)
	}

	normal := "/tmp/test.xml"
	if ExpandPath(normal) != normal {
		t.Errorf("ExpandPath(%q) modified absolute path", normal)
	}
}
