package library

import (
	"github.com/llttlltt/dj-library-tools/internal/models"
)

type MockLibrary struct {
	Tracks    []models.Track
	Playlists []models.ResourceGroup
}

func (m *MockLibrary) GetTracks() []models.Track    { return m.Tracks }
func (m *MockLibrary) GetPlaylists() []models.ResourceGroup { return m.Playlists }

func makeMockLibrary() *MockLibrary {
	return &MockLibrary{
		Tracks: []models.Track{
			{ID: "1", Title: "Track 1", Artist: "Artist A", BPM: 124.0},
			{ID: "2", Title: "Track 2", Artist: "Artist B", BPM: 128.0},
		},
		Playlists: []models.ResourceGroup{
			{Name: "Inbox", Type: models.GroupTypePlaylist, Items: 2},
			{Name: "Sets", Type: models.GroupTypeFolder, Items: 0},
		},
	}
}
