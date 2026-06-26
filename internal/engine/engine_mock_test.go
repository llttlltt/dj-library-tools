package engine

import (
	"github.com/llttlltt/dj-library-tools/internal/models"
)

type MockLibrary struct {
	Tracks    []models.Track
	Playlists []models.Node
}

func (m *MockLibrary) GetTracks() []models.Track    { return m.Tracks }
func (m *MockLibrary) GetPlaylists() []models.Node { return m.Playlists }

func makeMockLibrary() *MockLibrary {
	return &MockLibrary{
		Tracks: []models.Track{
			{ID: "1", Title: "Track 1", Artist: "Artist A", BPM: 124.0},
			{ID: "2", Title: "Track 2", Artist: "Artist B", BPM: 128.0},
		},
		Playlists: []models.Node{
			{Name: "Inbox", Type: 1, Entries: 2},
			{Name: "Sets", Type: 0, Entries: 0},
		},
	}
}
