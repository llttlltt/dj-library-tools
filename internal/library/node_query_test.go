package library

import (
	"testing"

	"github.com/llttlltt/dj-library-tools/internal/models"
)

func TestEngine_LsNodes(t *testing.T) {
	mock := &MockLibrary{
		Playlists: []models.ResourceGroup{
			{Name: "House", Type: models.GroupTypePlaylist},
		},
	}
	eng := NewEngine(mock)

	matched, err := eng.LsPlaylists("name:House")
	if err != nil {
		t.Fatal(err)
	}
	if len(matched) != 1 {
		t.Errorf("got %d, want 1", len(matched))
	}
}
