package cli

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/llttlltt/dj-library-tools/internal/engine"
	syncpkg "github.com/llttlltt/dj-library-tools/internal/sync"
	"github.com/llttlltt/dj-library-tools/internal/utils"
	"github.com/spf13/cobra"
)

var (
	moveTo     string
	moveFrom   string
	moveDryRun bool
)

var moveCmd = &cobra.Command{
	Use:   "move [resource] [query] --to [destination] [--from origin]",
	Short: "Move items between locations",
	Long: `Move items between Rekordbox locations.
For tracks, both --from and --to are required.
For playlists and folders, only --to (the parent folder) is required.

Example:
  djlt move rb/tracks "bpm:>130" --from "Inbox" --to "High Energy"
  djlt move rb/playlists "Deep House" --to "Genres"`,
	Args: cobra.MinimumNArgs(1),
	RunE: runMoveCmd,
}

func runMoveCmd(cmd *cobra.Command, args []string) error {
	if moveTo == "" {
		return fmt.Errorf("--to destination is required")
	}

	rbXML, path, err := loadXML()
	if err != nil {
		return err
	}

	eng := engine.NewEngine(rbXML)
	syncEng := syncpkg.NewEngine(nil, rbXML)

	query := ""
	if len(args) > 1 {
		query = strings.Join(args[1:], " ")
	}
	loc := utils.ParseLocation(args[0], query)

	if loc.Resource == "tracks" {
		if moveFrom == "" {
			return fmt.Errorf("--from origin is required when moving tracks")
		}
		return runMoveTracks(eng, syncEng, loc.Query, path)
	}

	return runMoveNodes(eng, syncEng, loc, path)
}

func runMoveTracks(eng *engine.Engine, syncEng *syncpkg.Engine, sourceQuery, path string) error {
	// 1. Resolve source tracks
	tracks, err := eng.Ls(sourceQuery)
	if err != nil {
		return err
	}
	if len(tracks) == 0 {
		fmt.Println("No tracks matched the source query.")
		return nil
	}
	var trackIDs []string
	for _, t := range tracks {
		trackIDs = append(trackIDs, strconv.Itoa(t.TrackID))
	}

	// 2. Resolve origin playlists
	org := utils.ParseLocation(moveFrom, "")
	origins, err := eng.LsPlaylists(org.Query)
	if err != nil || len(origins) == 0 {
		return fmt.Errorf("could not find origin playlist(s) matching %q", moveFrom)
	}

	// 3. Resolve target playlists
	tgt := utils.ParseLocation(moveTo, "")
	targets, err := eng.LsPlaylists(tgt.Query)
	if err != nil || len(targets) == 0 {
		return fmt.Errorf("could not find target playlist(s) matching %q", moveTo)
	}

	if moveDryRun {
		fmt.Printf("[Dry Run] Would move %d tracks from %d origins to %d targets\n", len(trackIDs), len(origins), len(targets))
		return nil
	}

	// 4. Perform Move
	for _, origin := range origins {
		syncEng.RemoveTracksFromPlaylist(origin.Node.Name, trackIDs)
	}
	for _, target := range targets {
		syncEng.AddTracksToPlaylist(target.Node.Name, trackIDs)
	}

	return syncEng.SaveXML(path)
}

func runMoveNodes(eng *engine.Engine, syncEng *syncpkg.Engine, loc utils.Location, path string) error {
	var targets []engine.NodeResult
	var nodeType int
	if loc.Resource == "playlists" {
		targets, _ = eng.LsPlaylists(loc.Query)
		nodeType = 1
	} else if loc.Resource == "folders" {
		targets, _ = eng.LsFolders(loc.Query)
		nodeType = 0
	} else {
		return fmt.Errorf("move only supports rb/tracks, rb/playlists, and rb/folders")
	}

	if len(targets) == 0 {
		fmt.Println("No resources found matching query.")
		return nil
	}

	if moveDryRun {
		for _, t := range targets {
			fmt.Printf("[Dry Run] Would move %s %q to folder %q\n", loc.Resource, t.Node.Name, moveTo)
		}
		return nil
	}

	for _, t := range targets {
		if !syncEng.MoveNode(t.Node.Name, int32(nodeType), moveTo) {
			fmt.Printf("Warning: failed to move %q\n", t.Node.Name)
			continue
		}
		fmt.Printf("Moved %s %q -> %q\n", loc.Resource, t.Node.Name, moveTo)
	}

	return syncEng.SaveXML(path)
}

func init() {
	moveCmd.Flags().StringVar(&moveTo, "to", "", "Destination playlist or folder")
	moveCmd.Flags().StringVar(&moveFrom, "from", "", "Origin playlist (required for tracks)")
	moveCmd.Flags().BoolVar(&moveDryRun, "dry-run", false, "Preview changes without writing")
	RootCmd.AddCommand(moveCmd)
}
