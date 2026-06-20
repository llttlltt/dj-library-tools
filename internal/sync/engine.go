package sync

import (
	"fmt"

	"github.com/llttlltt/dj-library-tools/internal/plex"
	"github.com/llttlltt/dj-library-tools/pkg/rekordbox"
)

type Engine struct {
	PlexClient *plex.Client
	RBXML      *rekordbox.RekordboxLibraryXML
	Matcher    *Matcher
}

func NewEngine(plexClient *plex.Client, rbXML *rekordbox.RekordboxLibraryXML) *Engine {
	return &Engine{
		PlexClient: plexClient,
		RBXML:      rbXML,
		Matcher:    NewMatcher(rbXML.Collection.TRACK),
	}
}

// SyncPlaylist takes a Plex playlist and adds it to the Rekordbox XML.
func (e *Engine) SyncPlaylist(baseURL, playlistKey string) error {
	_, err := e.PlexClient.GetPlaylistTracks(baseURL, playlistKey)
	if err != nil {
		return fmt.Errorf("failed to get plex tracks: %w", err)
	}

	// For simplicity, we'll find or create a "Plex Sync" folder in the root.
	var plexSyncFolder *rekordbox.Node
	for i := range e.RBXML.Playlists.Node.Node {
		if e.RBXML.Playlists.Node.Node[i].Name == "Plex Sync" && e.RBXML.Playlists.Node.Node[i].Type == 0 {
			plexSyncFolder = &e.RBXML.Playlists.Node.Node[i]
			break
		}
	}

	if plexSyncFolder == nil {
		e.RBXML.Playlists.Node.Node = append(e.RBXML.Playlists.Node.Node, rekordbox.Node{
			BaseNode: rekordbox.BaseNode{
				Type: 0,
				Name: "Plex Sync",
			},
		})
		plexSyncFolder = &e.RBXML.Playlists.Node.Node[len(e.RBXML.Playlists.Node.Node)-1]
	}

	// Create the playlist node
	// Note: We need to find the playlist name from the key if not passed
	// For now, let's assume we'll just name it based on the first track's info or similar
	// But ideally we'd pass the playlist title.

	return nil
}

// SaveXML writes the modified XML back to disk.
func (e *Engine) SaveXML(path string) error {
	return rekordbox.WriteRekordboxLibrary(path, e.RBXML)
}
