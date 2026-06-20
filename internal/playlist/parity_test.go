package playlist

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestParityWithLegacyScripts(t *testing.T) {
	// Setup temporary workspace
	tmpDir, err := os.MkdirTemp("", "djlt-parity-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	// 1. Test Extension Fixing Parity
	t.Run("ExtensionFixParity", func(t *testing.T) {
		m3uContent := "#EXTM3U\nmusic/track1.flac\nmusic/track2.wav\n"
		m3uPath := filepath.Join(tmpDir, "test.m3u")
		if err := os.WriteFile(m3uPath, []byte(m3uContent), 0644); err != nil {
			t.Fatal(err)
		}

		// Run Legacy Script
		cwd, _ := os.Getwd()
		scriptPath := filepath.Join(cwd, "../../scripts/legacy/fix_playlist_extensions.sh")
		legacyCmd := exec.Command("bash", scriptPath, m3uPath)
		if out, err := legacyCmd.CombinedOutput(); err != nil {
			t.Fatalf("Legacy script failed: %v\nOutput: %s", err, string(out))
		}
		legacyResult, _ := os.ReadFile(m3uPath)

		// Reset file for Go implementation
		if err := os.WriteFile(m3uPath, []byte(m3uContent), 0644); err != nil {
			t.Fatal(err)
		}

		// Run Go implementation
		opts := FixOptions{Ext: ".mp3"}
		if _, err := FixPlaylist(m3uPath, opts); err != nil {
			t.Fatalf("Go FixPlaylist failed: %v", err)
		}
		goResult, _ := os.ReadFile(m3uPath)

		if strings.TrimSpace(string(legacyResult)) != strings.TrimSpace(string(goResult)) {
			t.Errorf("Parity mismatch!\nLegacy:\n%s\nGo:\n%s", string(legacyResult), string(goResult))
		}
	})
}
