package cli

import (
	"bytes"
	"os"
	"strings"
	"testing"

	"github.com/llttlltt/dj-library-tools/pkg/rekordbox"
	"github.com/spf13/cobra"
	pflag "github.com/spf13/pflag"
)

// resetTestState resets all mutable package-level flag vars and clears cobra's
// Changed tracking so that flag state from one test cannot bleed into the next
// when the same RootCmd instance is reused across subtests.
func resetTestState(root *cobra.Command) {
	removeOrigins = nil
	syncTo = nil
	syncAppend = false
	moveTo = ""
	moveFrom = ""
	moveName = ""
	listSort = ""
	listStats = false
	dryRun = false
	verbose = false
	jsonOutput = false

	var resetChanged func(cmd *cobra.Command)
	resetChanged = func(cmd *cobra.Command) {
		cmd.Flags().VisitAll(func(f *pflag.Flag) { f.Changed = false })
		cmd.PersistentFlags().VisitAll(func(f *pflag.Flag) { f.Changed = false })
		for _, child := range cmd.Commands() {
			resetChanged(child)
		}
	}
	resetChanged(root)
}

func executeCommand(root *cobra.Command, args ...string) (string, error) {
	resetTestState(root)
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
			wantIn: "Would append to playlist",
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
			out, err := executeCommand(RootCmd, tt.args...)
			if (err != nil) != tt.wantErr {
				t.Errorf("executeCommand(%v) error = %v, wantErr %v", tt.args, err, tt.wantErr)
				return
			}
			if tt.wantIn != "" && !strings.Contains(out, tt.wantIn) {
				t.Errorf("executeCommand(%v) out = %q, want to contain %q", tt.args, out, tt.wantIn)
			}
		})
	}
}
