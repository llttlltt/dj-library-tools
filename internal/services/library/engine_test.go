package library

import (
	"testing"

	"github.com/llttlltt/dj-library-tools/internal/core/models"
)

type MockLibrary struct {
	Tracks    []models.Track
	Playlists []models.ResourceGroup
}

func (m *MockLibrary) GetResources(kind string) []models.Resource {
	var items []models.Resource
	switch kind {
	case "track":
		for _, t := range m.Tracks {
			items = append(items, t)
		}
	case "group":
		for _, p := range m.Playlists {
			items = append(items, p)
		}
	}
	return items
}

func (m *MockLibrary) GetMembershipMap() map[string][]models.PlaylistMembership { return nil }

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
			{Name: "House", Kind: models.GroupKindPlaylist},
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
