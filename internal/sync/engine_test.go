package sync

import (
	"fmt"
	"testing"

	"github.com/llttlltt/dj-library-tools/internal/library"
	"github.com/llttlltt/dj-library-tools/internal/rekordbox"
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
	eng := NewEngine(nil, library.NewRekordboxLibrary(lib))

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
	eng := NewEngine(nil, library.NewRekordboxLibrary(lib))

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
	eng := NewEngine(nil, library.NewRekordboxLibrary(lib))

	eng.InjectPlaylist("Summer", []string{"1"})
	eng.InjectPlaylist("Winter", []string{"2", "3"})

	folder := findFolder(lib, PlexSyncFolder)
	if len(folder.Node) != 2 {
		t.Errorf("expected 2 playlists in folder, got %d", len(folder.Node))
	}
	if rekordbox.DerefInt32(folder.Count) != 2 {
		t.Errorf("folder.Count: got %d, want 2", rekordbox.DerefInt32(folder.Count))
	}
}

func TestRemovePlaylist(t *testing.T) {
	lib := makeLibrary()
	eng := NewEngine(nil, library.NewRekordboxLibrary(lib))

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
	if rekordbox.DerefInt32(folder.Count) != 1 {
		t.Errorf("folder.Count after removal: got %d, want 1", rekordbox.DerefInt32(folder.Count))
	}
}

func TestRemovePlaylist_NotFound(t *testing.T) {
	lib := makeLibrary()
	eng := NewEngine(nil, library.NewRekordboxLibrary(lib))

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

func TestUpsertPlaylist_RootLevel(t *testing.T) {
	lib := makeLibrary()
	eng := NewEngine(nil, library.NewRekordboxLibrary(lib))

	result := eng.UpsertPlaylist("", "RootPlaylist", []string{"1", "2"}, -1)
	if result.Updated {
		t.Error("expected Updated=false for new playlist")
	}
	if result.TracksInjected != 2 {
		t.Errorf("TracksInjected: got %d, want 2", result.TracksInjected)
	}

	found := false
	for _, n := range lib.Playlists.Node.Node {
		if n.Name == "RootPlaylist" && n.Type == 1 {
			found = true
			if len(n.TRACK) != 2 {
				t.Errorf("track count: got %d, want 2", len(n.TRACK))
			}
		}
	}
	if !found {
		t.Error("root-level playlist not found")
	}
}

func TestUpsertPlaylist_NamedFolder(t *testing.T) {
	lib := makeLibrary()
	eng := NewEngine(nil, library.NewRekordboxLibrary(lib))

	eng.UpsertPlaylist("DJ Sets", "Techno Night", []string{"1"}, -1)
	eng.UpsertPlaylist("DJ Sets", "Techno Night", []string{"2", "3"}, -1)

	folder := findFolder(lib, "DJ Sets")
	if folder == nil {
		t.Fatal("folder DJ Sets not created")
	}
	if len(folder.Node) != 1 {
		t.Errorf("expected 1 playlist (upsert), got %d", len(folder.Node))
	}
	if len(folder.Node[0].TRACK) != 2 {
		t.Errorf("track count after upsert: got %d, want 2", len(folder.Node[0].TRACK))
	}
}

func TestAddTracksToPlaylist(t *testing.T) {
	lib := makeLibrary()
	eng := NewEngine(nil, library.NewRekordboxLibrary(lib))

	eng.UpsertPlaylist("", "Picks", []string{"1", "2"}, -1)

	found, added := eng.AddTracksToPlaylist("Picks", []string{"2", "3", "4"})
	if !found {
		t.Error("expected playlist to be found")
	}
	// "2" is a duplicate — only "3" and "4" should be added
	if added != 2 {
		t.Errorf("added: got %d, want 2", added)
	}

	for _, n := range lib.Playlists.Node.Node {
		if n.Name == "Picks" && n.Type == 1 {
			if len(n.TRACK) != 4 {
				t.Errorf("total tracks: got %d, want 4", len(n.TRACK))
			}
			return
		}
	}
	t.Error("playlist Picks not found after add")
}

func TestAddTracksToPlaylist_NotFound(t *testing.T) {
	lib := makeLibrary()
	eng := NewEngine(nil, library.NewRekordboxLibrary(lib))

	found, added := eng.AddTracksToPlaylist("Nonexistent", []string{"1"})
	if found {
		t.Error("expected found=false for missing playlist")
	}
	if added != 0 {
		t.Errorf("expected added=0, got %d", added)
	}
}

func TestRemoveTracksFromPlaylist(t *testing.T) {
	lib := makeLibrary()
	eng := NewEngine(nil, library.NewRekordboxLibrary(lib))

	eng.UpsertPlaylist("", "Picks", []string{"1", "2", "3", "4"}, -1)

	found, removed := eng.RemoveTracksFromPlaylist("Picks", []string{"2", "4"})
	if !found {
		t.Error("expected playlist to be found")
	}
	if removed != 2 {
		t.Errorf("removed: got %d, want 2", removed)
	}

	for _, n := range lib.Playlists.Node.Node {
		if n.Name == "Picks" && n.Type == 1 {
			if len(n.TRACK) != 2 {
				t.Errorf("remaining tracks: got %d, want 2", len(n.TRACK))
			}
			for _, tr := range n.TRACK {
				if tr.Key == "2" || tr.Key == "4" {
					t.Errorf("track %q should have been removed", tr.Key)
				}
			}
			return
		}
	}
	t.Error("playlist Picks not found")
}

func TestRemoveTracksFromPlaylist_NotFound(t *testing.T) {
	lib := makeLibrary()
	eng := NewEngine(nil, library.NewRekordboxLibrary(lib))

	found, removed := eng.RemoveTracksFromPlaylist("Nonexistent", []string{"1"})
	if found {
		t.Error("expected found=false for missing playlist")
	}
	if removed != 0 {
		t.Errorf("expected removed=0, got %d", removed)
	}
}

func TestRenameNode(t *testing.T) {
	lib := makeLibrary()
	eng := NewEngine(nil, library.NewRekordboxLibrary(lib))

	eng.InjectPlaylist("Summer", []string{"1"})

	ok := eng.RenameNode("Summer", "Summer 2025", 1)
	if !ok {
		t.Error("expected RenameNode to return true")
	}

	folder := findFolder(lib, PlexSyncFolder)
	if folder.Node[0].Name != "Summer 2025" {
		t.Errorf("name after rename: got %q, want %q", folder.Node[0].Name, "Summer 2025")
	}
}

func TestRenameNode_NotFound(t *testing.T) {
	lib := makeLibrary()
	eng := NewEngine(nil, library.NewRekordboxLibrary(lib))

	ok := eng.RenameNode("Nonexistent", "New Name", 1)
	if ok {
		t.Error("expected RenameNode to return false for missing node")
	}
}

func TestMoveNode(t *testing.T) {
	lib := makeLibrary()
	eng := NewEngine(nil, library.NewRekordboxLibrary(lib))

	eng.InjectPlaylist("Summer", []string{"1"})

	ok := eng.MoveNode("Summer", 1, "Archive")
	if !ok {
		t.Error("expected MoveNode to return true")
	}

	// Should no longer be in PlexSyncFolder
	plexFolder := findFolder(lib, PlexSyncFolder)
	for _, n := range plexFolder.Node {
		if n.Name == "Summer" {
			t.Error("playlist still present in original folder after move")
		}
	}

	// Should be in Archive
	archiveFolder := findFolder(lib, "Archive")
	if archiveFolder == nil {
		t.Fatal("Archive folder not created")
	}
	found := false
	for _, n := range archiveFolder.Node {
		if n.Name == "Summer" {
			found = true
		}
	}
	if !found {
		t.Error("playlist not found in Archive after move")
	}
}

func TestRemoveNode_Playlist(t *testing.T) {
	lib := makeLibrary()
	eng := NewEngine(nil, library.NewRekordboxLibrary(lib))

	eng.InjectPlaylist("Summer", []string{"1"})
	eng.InjectPlaylist("Winter", []string{"2"})

	ok := eng.RemoveNode("Summer", 1)
	if !ok {
		t.Error("expected RemoveNode to return true")
	}

	folder := findFolder(lib, PlexSyncFolder)
	if len(folder.Node) != 1 || folder.Node[0].Name != "Winter" {
		t.Errorf("unexpected folder state after remove: %+v", folder.Node)
	}
}

func TestRemoveNode_NotFound(t *testing.T) {
	lib := makeLibrary()
	eng := NewEngine(nil, library.NewRekordboxLibrary(lib))

	ok := eng.RemoveNode("Nonexistent", 1)
	if ok {
		t.Error("expected RemoveNode to return false for missing node")
	}
}

func TestWrapperBackcompat(t *testing.T) {
	lib := makeLibrary()
	eng := NewEngine(nil, library.NewRekordboxLibrary(lib))

	r := eng.InjectPlaylist("Test", []string{"1"})
	if r.Updated || r.TracksInjected != 1 {
		t.Errorf("InjectPlaylist wrapper: unexpected result %+v", r)
	}

	removed := eng.RemovePlaylist("Test")
	if !removed {
		t.Error("RemovePlaylist wrapper: expected true")
	}
}
