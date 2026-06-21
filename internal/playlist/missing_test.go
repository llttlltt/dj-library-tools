package playlist

import (
	"os"
	"path/filepath"
	"testing"
)

func TestMissingFileReporting(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "djlt-missing-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	// Create a playlist with relative paths
	// We'll create one track that exists and one that doesn't
	mediaDir := filepath.Join(tmpDir, "Media")
	if err := os.Mkdir(mediaDir, 0755); err != nil {
		t.Fatal(err)
	}

	existingTrack := filepath.Join(mediaDir, "exists.flac")
	if err := os.WriteFile(existingTrack, []byte("dummy content"), 0644); err != nil {
		t.Fatal(err)
	}

	// We also need to create the .mp3 version because the existence check 
	// happens on the TRANSFORMED path
	existingMP3 := filepath.Join(mediaDir, "exists.mp3")
	if err := os.WriteFile(existingMP3, []byte("dummy content"), 0644); err != nil {
		t.Fatal(err)
	}

	playlistContent := "#EXTM3U\n../Media/exists.flac\n../Media/missing.flac\n"
	playlistDir := filepath.Join(tmpDir, "Playlists")
	if err := os.Mkdir(playlistDir, 0755); err != nil {
		t.Fatal(err)
	}
	playlistPath := filepath.Join(playlistDir, "test.m3u")
	if err := os.WriteFile(playlistPath, []byte(playlistContent), 0644); err != nil {
		t.Fatal(err)
	}

	opts := FixOptions{Exts: []string{".mp3"}}
	result, err := FixPlaylist(playlistPath, opts)
	if err != nil {
		t.Fatalf("FixPlaylist failed: %v", err)
	}

	if result.TotalTracks != 2 {
		t.Errorf("Expected 2 total tracks, got %d", result.TotalTracks)
	}

	if len(result.SkippedTracks) != 1 {
		t.Errorf("Expected 1 missing track, got %d", len(result.SkippedTracks))
	}

	expectedMissing := "../Media/missing.flac"
	if result.SkippedTracks[0] != expectedMissing {
		t.Errorf("Expected skipped track %s, got %s", expectedMissing, result.SkippedTracks[0])
	}
}
