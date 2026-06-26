package cli

import (
	"fmt"
	"strings"

	"github.com/llttlltt/dj-library-tools/internal/engine"
	syncpkg "github.com/llttlltt/dj-library-tools/internal/sync"
	"github.com/spf13/cobra"
)

var (
	moveTo   string
	moveFrom string
)

var moveCmd = &cobra.Command{
	Use:     "move [resource] [query] --to [destination] [--from origin]",
	Aliases: []string{"mv"},
	Short:   "Move items between locations",
	Long: `Move items between Rekordbox locations.
For tracks, both --from and --to are required.
For playlists and folders, only --to (the parent folder) is required.

Example:
  djlt move rb/tracks "bpm:>130" --from "name:Inbox" --to "name:'High Energy'"
  djlt move rb/playlists "name:'Deep House'" --to "name:Genres"`,
	Args: cobra.MinimumNArgs(1),
	RunE: runMoveCmd,
}

func runMoveCmd(cmd *cobra.Command, args []string) error {
	if moveTo == "" {
		return fmt.Errorf("--to destination is required")
	}

	rbXML, path, err := loadXMLFunc()
	if err != nil {
		return err
	}

	syncEng := syncpkg.NewEngine(nil, engine.NewRekordboxLibrary(rbXML))

	queryOverride := ""
	if len(args) > 1 {
		queryOverride = strings.Join(args[1:], " ")
	}
	sel, err := ResolveSelection(args[0], queryOverride)
	if err != nil {
		return err
	}

	if sel.Location.Resource == "tracks" {
		if moveFrom == "" {
			return fmt.Errorf("--from origin is required when moving tracks")
		}
		return runMoveTracks(syncEng, sel, path)
	}

	return runMoveNodes(syncEng, sel, path)
}

func runMoveTracks(syncEng *syncpkg.Engine, sel *Selection, path string) error {
	if len(sel.Tracks) == 0 {
		fmt.Println("No tracks matched the source query.")
		return nil
	}
	var trackIDs []string
	for _, t := range sel.Tracks {
		trackIDs = append(trackIDs, t.ID)
	}

	// 2. Resolve origin playlists
	org, err := ResolveSelection(moveFrom, "")
	if err != nil || len(org.Nodes) == 0 {
		return fmt.Errorf("could not find origin playlist(s) matching %q", moveFrom)
	}

	// 3. Resolve target playlists
	tgt, err := ResolveSelection(moveTo, "")
	if err != nil || len(tgt.Nodes) == 0 {
		return fmt.Errorf("could not find target playlist(s) matching %q", moveTo)
	}

	if dryRun {
		fmt.Printf("[Dry Run] Would move %d tracks from %d origins to %d targets\n", len(trackIDs), len(org.Nodes), len(tgt.Nodes))
		return nil
	}

	// 4. Perform Move
	for _, origin := range org.Nodes {
		if verbose {
			fmt.Printf("Removing tracks from origin playlist %q...\n", origin.Name)
		}
		syncEng.RemoveTracksFromPlaylist(origin.Name, trackIDs)
	}
	for _, target := range tgt.Nodes {
		if verbose {
			fmt.Printf("Adding tracks to target playlist %q...\n", target.Name)
		}
		syncEng.AddTracksToPlaylist(target.Name, trackIDs)
	}

	return syncEng.SaveXML(path)
}

func runMoveNodes(syncEng *syncpkg.Engine, sel *Selection, path string) error {
	var nodeType int
	if sel.Location.Resource == "playlists" {
		nodeType = 1
	} else if sel.Location.Resource == "folders" {
		nodeType = 0
	} else {
		return fmt.Errorf("move only supports rb/tracks, rb/playlists, and rb/folders")
	}

	if len(sel.Nodes) == 0 {
		fmt.Println("No resources found matching query.")
		return nil
	}

	if dryRun {
		for _, t := range sel.Nodes {
			fmt.Printf("[Dry Run] Would move %s %q to folder %q\n", sel.Location.Resource, t.Name, moveTo)
		}
		return nil
	}

	for _, t := range sel.Nodes {
		if verbose {
			fmt.Printf("Moving %s %q into folder %q...\n", sel.Location.Resource, t.Name, moveTo)
		}
		if !syncEng.MoveNode(t.Name, int32(nodeType), moveTo) {
			fmt.Printf("Warning: failed to move %q\n", t.Name)
			continue
		}
		fmt.Printf("Moved %s %q -> %q\n", sel.Location.Resource, t.Name, moveTo)
	}

	return syncEng.SaveXML(path)
}

func init() {
	moveCmd.Flags().StringVar(&moveTo, "to", "", "Destination playlist or folder")
	moveCmd.Flags().StringVar(&moveFrom, "from", "", "Origin playlist (required for tracks)")
	RootCmd.AddCommand(moveCmd)
}
