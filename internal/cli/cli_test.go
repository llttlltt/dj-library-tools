package cli

import (
	"bytes"
	"os"
	"strings"
	"testing"

	"github.com/llttlltt/dj-library-tools/internal/rekordbox"
)

func resetTestState() {
	dryRun = false
	verbose = false
	jsonOutput = false
	filePath = "mock.xml"
}

func executeCommand(args ...string) (string, error) {
	resetTestState()
	root := NewRootCmd()
	buf := new(bytes.Buffer)
	root.SetOut(buf)
	root.SetErr(buf)
	root.SetArgs(args[:])

	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	err := root.Execute()

	w.Close()
	os.Stdout = old
	var outBuf bytes.Buffer
	outBuf.ReadFrom(r)

	return buf.String() + outBuf.String(), err
}

func mockLoadXML() (*rekordbox.RekordboxLibraryXML, string, error) {
	return &rekordbox.RekordboxLibraryXML{
		Collection: rekordbox.Collection{
			TRACK: []rekordbox.Track{
				{TrackID: 1, Name: "Test Track", Artist: "Test Artist", Location: "file://localhost/test.mp3"},
			},
		},
		Playlists: rekordbox.Playlists{
			Node: rekordbox.RootNode{
				Name: "ROOT",
				Type: 0,
				Node: []rekordbox.Node{
					{Name: "Inbox", Type: 1, Entries: rekordbox.PtrInt32(0)},
				},
			},
		},
	}, "mock.xml", nil
}

func TestCommandConsistency(t *testing.T) {
	loadXMLFunc = mockLoadXML

	tests := []struct {
		name     string
		args     []string
		wantIn   string
		wantErr  bool
	}{
		{
			name: "list rb/tracks positional query",
			args: []string{"ls", "rb/tracks", "title:'Test Track'"},
			wantIn: "Test Track",
		},
		{
			name: "list rb/playlists positional query",
			args: []string{"ls", "rb/playlists", "name:Inbox"},
			wantIn: "Inbox",
		},
		{
			name: "stat merged into list --stats",
			args: []string{"ls", "rb/tracks", "title:'Test Track'", "--stats"},
			wantIn: "Selection Summary",
		},
		{
			name: "add tracks merged into sync --append",
			args: []string{"sync", "rb/tracks", "title:Test", "--to", "rb/playlists name:Inbox", "--append", "--dry-run"},
			wantIn: "Would append to",
		},
		{
			name: "remove tracks from playlist",
			args: []string{"rm", "rb/tracks", "title:'Test Track'", "--from", "rb/playlists name:Inbox", "--dry-run"},
			wantIn: "Would remove 1 tracks from playlist \"Inbox\"",
		},
		{
			name: "move tracks requires --from and --to",
			args: []string{"mv", "rb/tracks", "title:Test", "--to", "Target"},
			wantErr: true,
		},
		{
			name: "rename merged into move --name",
			args: []string{"mv", "rb/playlists", "name:Inbox", "--name", "NewInbox", "--dry-run"},
			wantIn: "Would rename \"Inbox\" to \"NewInbox\"",
		},
		{
			name: "remove playlist resource",
			args: []string{"rm", "rb/playlists", "name:Inbox", "--dry-run"},
			wantIn: "Would delete playlist \"Inbox\"",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			out, err := executeCommand(tt.args...)
			if (err != nil) != tt.wantErr {
				t.Errorf("executeCommand(%v) error = %v, wantErr %v", tt.args, err, tt.wantErr)
				return
			}
			if tt.wantIn != "" && !strings.Contains(out, tt.wantIn) && !tt.wantErr {
				t.Errorf("executeCommand(%v) out = %q, want to contain %q", tt.args, out, tt.wantIn)
			}
		})
	}
}
