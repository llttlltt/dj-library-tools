package m3u

import (
	"testing"
)

func TestFormatPath(t *testing.T) {
	tests := []struct {
		path     string
		ext      string
		expected string
	}{
		{"music/track.flac", ".mp3", "music/track.mp3"},
		{"music/track.wav", "mp3", "music/track.mp3"},
		{"music/track", ".m3u8", "music/track.m3u8"},
		{"music/track.m3u", ".m3u8", "music/track.m3u8"},
	}

	for _, tt := range tests {
		result := FormatPath(tt.path, tt.ext)
		if result != tt.expected {
			t.Errorf("FormatPath(%q, %q) = %q; want %q", tt.path, tt.ext, result, tt.expected)
		}
	}
}

func TestIsM3UHeader(t *testing.T) {
	if !IsM3UHeader("#EXTM3U") {
		t.Error("Expected true for #EXTM3U")
	}
	if IsM3UHeader("#EXTINF:0,Artist - Title") {
		t.Error("Expected false for #EXTINF line")
	}
}

func TestIsExtInfLine(t *testing.T) {
	if !IsExtInfLine("#EXTINF:0,Artist - Title") {
		t.Error("Expected true for #EXTINF line")
	}
	if IsExtInfLine("#EXTM3U") {
		t.Error("Expected false for #EXTM3U line")
	}
}
