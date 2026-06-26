package cli

import (
	"bytes"
	"os"
	"strings"
	"testing"

	"github.com/llttlltt/dj-library-tools/pkg/rekordbox"
	"github.com/spf13/cobra"
)

func executeCommand(root *cobra.Command, args ...string) (string, error) {
	buf := new(bytes.Buffer)
	root.SetOut(buf)
	root.SetErr(buf)
	root.SetArgs(args)

	// Since we use fmt.Printf in many places, we need to capture stdout
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	err := root.Execute()

	w.Close()
	os.Stdout = old
	var outBuf bytes.Buffer
	outBuf.ReadFrom(r)

	// Combine both buffers
	return buf.String() + outBuf.String(), err
}

func mockLoadXML() (*rekordbox.RekordboxLibraryXML, string, error) {
	return &rekordbox.RekordboxLibraryXML{
		Collection: rekordbox.Collection{
			TRACK: []rekordbox.Track{
				{TrackID: 1, Name: "Test Track", Artist: "Test Artist"},
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
	// Override the XML loader for all tests
	loadXMLFunc = mockLoadXML

	tests := []struct {
		name     string
		args     []string
		wantIn   string
		wantOut  string
		wantErr  bool
	}{
		{
			name: "list rb/tracks positional query",
			args: []string{"list", "rb/tracks", "title:'Test Track'"},
			wantIn: "Test Track",
		},
		{
			name: "list rb/playlists positional query",
			args: []string{"list", "rb/playlists", "name:Inbox"},
			wantIn: "Inbox",
		},
		{
			name: "stat rb/tracks positional query",
			args: []string{"stat", "rb/tracks", "title:'Test Track'"},
			wantIn: "Total Tracks:   1",
		},
		{
			name: "add requires --to",
			args: []string{"add", "rb/tracks", "title:Test"},
			wantErr: true,
		},
		{
			name: "add rb/tracks to playlist",
			args: []string{"add", "rb/tracks", "title:'Test Track'", "--to", "rb/playlists name:Inbox", "--dry-run"},
			wantIn: "Would add 1 tracks to playlist \"Inbox\"",
		},
		{
			name: "remove requires --from",
			args: []string{"remove", "rb/tracks", "title:Test"},
			wantErr: true,
		},
		{
			name: "remove rb/tracks from playlist",
			args: []string{"remove", "rb/tracks", "title:'Test Track'", "--from", "rb/playlists name:Inbox", "--dry-run"},
			wantIn: "Would remove 1 tracks from playlist \"Inbox\"",
		},
		{
			name: "move tracks requires --from and --to",
			args: []string{"move", "rb/tracks", "title:Test", "--to", "Target"},
			wantErr: true,
		},
		{
			name: "rename requires --to",
			args: []string{"rename", "rb/playlists", "name:Inbox"},
			wantErr: true,
		},
		{
			name: "delete requires resource",
			args: []string{"delete"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			out, err := executeCommand(RootCmd, tt.args...)
			if (err != nil) != tt.wantErr {
				t.Errorf("executeCommand() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantIn != "" && !strings.Contains(out, tt.wantIn) {
				t.Errorf("executeCommand() out = %q, want to contain %q", out, tt.wantIn)
			}
		})
	}
}
