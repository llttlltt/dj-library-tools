package m3u

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/llttlltt/dj-library-tools/internal/core/models"
	"github.com/llttlltt/dj-library-tools/internal/providers"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestM3UProvider_Load(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "m3u-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	m3uPath := filepath.Join(tmpDir, "test.m3u8")
	content := "#EXTM3U\n#EXTINF:180,Artist One - Title One\ntrack1.mp3\n#EXTINF:240,Artist Two - Title Two\ntrack2.mp3\n"
	err = os.WriteFile(m3uPath, []byte(content), 0644)
	require.NoError(t, err)

	p, err := NewM3UProvider(m3uPath)
	require.NoError(t, err)

	tracks, err := p.Tracks().List(context.Background(), provider.ExecutionContext{}, "")
	require.NoError(t, err)
	assert.Len(t, tracks, 2)

	assert.Equal(t, "Artist One - Title One", tracks[0].Display)
	assert.Equal(t, filepath.Join(tmpDir, "track1.mp3"), tracks[0].Location)

	assert.Equal(t, "Artist Two - Title Two", tracks[1].Display)
	assert.Equal(t, filepath.Join(tmpDir, "track2.mp3"), tracks[1].Location)
}

func TestM3UProvider_AddRemoveSave(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "m3u-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	m3uPath := filepath.Join(tmpDir, "sync.m3u8")
	p, err := NewM3UProvider(m3uPath)
	require.NoError(t, err)

	// Add tracks
	newTracks := []models.Track{
		{Display: "New Track Display", Location: "/tmp/new.mp3"},
	}
	added, err := p.Tracks().Groups().Add(context.Background(), provider.ExecutionContext{}, newTracks, models.ResourceGroup{})
	assert.NoError(t, err)
	assert.Equal(t, 1, added)

	// Save
	err = p.System().Save(context.Background(), provider.ExecutionContext{}, m3uPath)
	assert.NoError(t, err)

	// Reload and verify
	p2, err := NewM3UProvider(m3uPath)
	assert.NoError(t, err)
	tracks, _ := p2.Tracks().List(context.Background(), provider.ExecutionContext{}, "")
	assert.Len(t, tracks, 1)
	assert.Equal(t, "New Track Display", tracks[0].Display)

	// Remove
	removed, err := p2.Tracks().Groups().Remove(context.Background(), provider.ExecutionContext{}, tracks, models.ResourceGroup{})
	assert.NoError(t, err)
	assert.Equal(t, 1, removed)
	assert.Len(t, p2.tracks, 0)
}

func TestM3UProvider_FormatCorrectness(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "m3u-format-*")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	tracks := []models.Track{
		{Display: "Track One", Location: "/path/to/one.mp3", Duration: 120},
	}

	// 1. Test .m3u8 (Extended)
	m3u8Path := filepath.Join(tmpDir, "test.m3u8")
	p8, _ := NewM3UProvider(m3u8Path)
	p8.tracks = tracks
	err = p8.System().Save(context.Background(), provider.ExecutionContext{Apply: true}, m3u8Path)
	require.NoError(t, err)

	content8, _ := os.ReadFile(m3u8Path)
	expected8 := "#EXTM3U\n#EXTINF:120,Track One\n/path/to/one.mp3\n"
	assert.Equal(t, expected8, string(content8))

	// 2. Test .m3u (Basic)
	m3uPath := filepath.Join(tmpDir, "test.m3u")
	p, _ := NewM3UProvider(m3uPath)
	p.tracks = tracks
	err = p.System().Save(context.Background(), provider.ExecutionContext{Apply: true}, m3uPath)
	require.NoError(t, err)

	content, _ := os.ReadFile(m3uPath)
	expected := "#EXTM3U\n/path/to/one.mp3\n"
	assert.Equal(t, expected, string(content))
}
