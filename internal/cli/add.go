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
	addTargets []string
	addForce   bool
)

var addCmd = &cobra.Command{
	Use:   "add [source-resource] [source-query] --to [target-resource] [target-query]",
	Short: "Add items from a source to one or more targets",
	Long: `Add items from a source selection to one or more target resources.
Currently supports adding tracks (rb/tracks) to playlists (rb/playlists).

Example:
  djlt add rb/tracks artist:Four --to "rb/playlists name:Inbox"`,
	Args: cobra.MinimumNArgs(1),
	RunE: runAddCmd,
}

func runAddCmd(cmd *cobra.Command, args []string) error {
	if len(addTargets) == 0 {
		return fmt.Errorf("at least one --to target is required")
	}

	rbXML, path, err := loadXMLFunc()
	if err != nil {
		return err
	}

	eng := engine.NewEngine(rbXML)
	syncEng := syncpkg.NewEngine(nil, rbXML)

	// 1. Resolve source
	sourceQuery := ""
	if len(args) > 1 {
		sourceQuery = strings.Join(args[1:], " ")
	}
	src := utils.ParseLocation(args[0], sourceQuery)

	if src.Provider != "rb" || src.Resource != "tracks" {
		return fmt.Errorf("currently only rb/tracks is supported as a source for add")
	}

	tracks, err := eng.Ls(src.Query)
	if err != nil {
		return fmt.Errorf("failed to resolve source tracks: %w", err)
	}
	if len(tracks) == 0 {
		fmt.Println("No source tracks found matching query.")
		return nil
	}

	var trackIDs []string
	for _, t := range tracks {
		trackIDs = append(trackIDs, strconv.Itoa(t.TrackID))
	}

	// 2. Resolve targets and apply
	for _, targetStr := range addTargets {
		tgt := utils.ParseLocation(targetStr, "")
		if tgt.Provider != "rb" || tgt.Resource != "playlists" {
			return fmt.Errorf("currently only rb/playlists is supported as a target for add, got %q", targetStr)
		}

		targets, err := eng.LsPlaylists(tgt.Query)
		if err != nil {
			return fmt.Errorf("failed to resolve target playlists: %w", err)
		}
		if len(targets) == 0 {
			return fmt.Errorf("no target playlists matched query %q", tgt.Query)
		}

		if dryRun {
			for _, target := range targets {
				fmt.Printf("[Dry Run] Would add %d tracks to playlist %q\n", len(trackIDs), target.Node.Name)
			}
			continue
		}

		for _, target := range targets {
			if verbose {
				fmt.Printf("Adding %d tracks to playlist %q...\n", len(trackIDs), target.Node.Name)
				for _, id := range trackIDs {
					fmt.Printf("  + Track ID: %s\n", id)
				}
			}
			found, added := syncEng.AddTracksToPlaylist(target.Node.Name, trackIDs)
			if !found {
				fmt.Printf("Warning: playlist %q not found during add\n", target.Node.Name)
				continue
			}
			fmt.Printf("Added %d tracks to %q\n", added, target.Node.Name)
		}
	}

	if dryRun {
		return nil
	}

	return syncEng.SaveXML(path)
}

func init() {
	addCmd.Flags().StringSliceVar(&addTargets, "to", []string{}, "Target resource(s) to add to (repeatable)")
	addCmd.Flags().BoolVar(&addForce, "force", false, "Allow adding duplicates (if supported by target)")
	
	RootCmd.AddCommand(addCmd)
}
