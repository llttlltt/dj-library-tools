package library

import (
	"testing"

	"github.com/llttlltt/dj-library-tools/internal/models"
)

type MockLibrary struct {
	Tracks    []models.Track
	Playlists []models.ResourceGroup
}

func (m *MockLibrary) GetTracks() []models.Track            { return m.Tracks }
func (m *MockLibrary) GetPlaylists() []models.ResourceGroup { return m.Playlists }
func (m *MockLibrary) GetMembershipMap() map[string][]string { return nil }

func TestEngine_Ls(t *testing.T) {
	mock := &MockLibrary{
		Tracks: []models.Track{
			{ID: "1", Title: "Oceans", Artist: "Four Tet"},
		},
	}
	eng := NewEngine(mock)

	matched, err := eng.Ls("title:Oceans", nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(matched) != 1 {
		t.Errorf("got %d, want 1", len(matched))
	}
}

func TestEngine_LsGroups(t *testing.T) {
	mock := &MockLibrary{
		Playlists: []models.ResourceGroup{
			{Name: "House", Type: models.GroupTypePlaylist},
		},
	}
	eng := NewEngine(mock)

	matched, err := eng.LsGroups("name:House")
	if err != nil {
		t.Fatal(err)
	}
	if len(matched) != 1 {
		t.Errorf("got %d, want 1", len(matched))
	}
}
