package media

import (
	"testing"
)

func TestApplyPathMap(t *testing.T) {
	cfg := &Config{
		PathMaps: map[string]string{
			"/remote/path": "/local/path",
		},
	}
	tr := NewTranscoder(cfg)

	tests := []struct {
		input    string
		expected string
	}{
		{"/remote/path/track.mp3", "/local/path/track.mp3"},
		{"/other/path/track.mp3", "/other/path/track.mp3"},
	}

	for _, tt := range tests {
		result := tr.ApplyPathMap(tt.input)
		if result != tt.expected {
			t.Errorf("ApplyPathMap(%q) = %q; want %q", tt.input, result, tt.expected)
		}
	}
}

func TestSanitizePathComponent(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"AC/DC", "AC-DC"},
		{"artist: name", "artist- name"},
		{"title?", "title"},
		{`say "hello"`, "say hello"},
		{"pipe|sep", "pipe-sep"},
		{"  padded  ", "padded"},
		{"normal", "normal"},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := sanitizePathComponent(tt.input)
			if got != tt.want {
				t.Errorf("sanitizePathComponent(%q) = %q; want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestFormatPath(t *testing.T) {
	cfg := DefaultConfig()
	tr := NewTranscoder(cfg)

	meta := PathMetadata{
		Artist: "Four Tet",
		Album:  "Sixteen Oceans",
		Title:  "Lush",
	}

	expected := "Four Tet - Sixteen Oceans - Lush.mp3"
	result, err := tr.FormatPath(meta)
	if err != nil {
		t.Fatalf("FormatPath failed: %v", err)
	}
	if result != expected {
		t.Errorf("FormatPath() = %q; want %q", result, expected)
	}
}
