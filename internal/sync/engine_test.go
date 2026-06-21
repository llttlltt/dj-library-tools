package sync

import (
	"fmt"
	"testing"

	"github.com/llttlltt/dj-library-tools/pkg/rekordbox"
)

func makeLibrary() *rekordbox.RekordboxLibraryXML {
	tracks := make([]rekordbox.Track, 5)
	for i := range tracks {
		tracks[i] = rekordbox.Track{
			TrackID: i + 1,
			Name:    fmt.Sprintf("Track %d", i+1),
			Artist:  "Test Artist",
		}
	}
	return &rekordbox.RekordboxLibraryXML{
		Collection: rekordbox.Collection{TRACK: tracks},
		Playlists:  rekordbox.Playlists{},
	}
}

func TestInjectPlaylist_Create(t *testing.T) {
	lib := makeLibrary()
	eng := NewEngine(nil, lib)

	result := eng.InjectPlaylist("Summer", []string{"1", "2", "3"})

	if result.Updated {
		t.Error("expected Updated=false for a new playlist")
	}
	if result.TracksInjected != 3 {
		t.Errorf("TracksInjected: got %d, want 3", result.TracksInjected)
	}

	folder := findFolder(lib, PlexSyncFolder)
	if folder == nil {
		t.Fatal("PlexSyncFolder not created")
	}
	if len(folder.Node) != 1 {
		t.Fatalf("expected 1 playlist in folder, got %d", len(folder.Node))
	}
	if folder.Node[0].Name != "Summer" {
		t.Errorf("playlist name: got %q, want %q", folder.Node[0].Name, "Summer")
	}
	if len(folder.Node[0].TRACK) != 3 {
		t.Errorf("track count: got %d, want 3", len(folder.Node[0].TRACK))
	}
}

func TestInjectPlaylist_Update(t *testing.T) {
	lib := makeLibrary()
	eng := NewEngine(nil, lib)

	eng.InjectPlaylist("Summer", []string{"1", "2", "3"})
	result := eng.InjectPlaylist("Summer", []string{"4", "5"})

	if !result.Updated {
		t.Error("expected Updated=true on second inject of same name")
	}
	if result.TracksInjected != 2 {
		t.Errorf("TracksInjected: got %d, want 2", result.TracksInjected)
	}

	folder := findFolder(lib, PlexSyncFolder)
	if len(folder.Node) != 1 {
		t.Errorf("expected 1 playlist (upsert), got %d", len(folder.Node))
	}
	if len(folder.Node[0].TRACK) != 2 {
		t.Errorf("track count after update: got %d, want 2", len(folder.Node[0].TRACK))
	}
}

func TestInjectPlaylist_MultiplePlaylists(t *testing.T) {
	lib := makeLibrary()
	eng := NewEngine(nil, lib)

	eng.InjectPlaylist("Summer", []string{"1"})
	eng.InjectPlaylist("Winter", []string{"2", "3"})

	folder := findFolder(lib, PlexSyncFolder)
	if len(folder.Node) != 2 {
		t.Errorf("expected 2 playlists in folder, got %d", len(folder.Node))
	}
	if folder.Count != 2 {
		t.Errorf("folder.Count: got %d, want 2", folder.Count)
	}
}

func TestRemovePlaylist(t *testing.T) {
	lib := makeLibrary()
	eng := NewEngine(nil, lib)

	eng.InjectPlaylist("Summer", []string{"1"})
	eng.InjectPlaylist("Winter", []string{"2"})

	removed := eng.RemovePlaylist("Summer")
	if !removed {
		t.Error("expected RemovePlaylist to return true")
	}

	folder := findFolder(lib, PlexSyncFolder)
	if len(folder.Node) != 1 {
		t.Errorf("expected 1 playlist after removal, got %d", len(folder.Node))
	}
	if folder.Node[0].Name != "Winter" {
		t.Errorf("remaining playlist: got %q, want %q", folder.Node[0].Name, "Winter")
	}
	if folder.Count != 1 {
		t.Errorf("folder.Count after removal: got %d, want 1", folder.Count)
	}
}

func TestRemovePlaylist_NotFound(t *testing.T) {
	lib := makeLibrary()
	eng := NewEngine(nil, lib)

	removed := eng.RemovePlaylist("Nonexistent")
	if removed {
		t.Error("expected RemovePlaylist to return false for missing playlist")
	}
}

// findFolder is a test helper that locates a top-level folder by name.
func findFolder(lib *rekordbox.RekordboxLibraryXML, name string) *rekordbox.Node {
	for i := range lib.Playlists.Node.Node {
		if lib.Playlists.Node.Node[i].Name == name {
			return &lib.Playlists.Node.Node[i]
		}
	}
	return nil
}
