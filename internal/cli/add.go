package cli

import (
	"fmt"
	"strings"

	"github.com/llttlltt/dj-library-tools/internal/engine"
	syncpkg "github.com/llttlltt/dj-library-tools/internal/sync"
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

	syncEng := syncpkg.NewEngine(nil, engine.NewRekordboxLibrary(rbXML))

	// 1. Resolve source
	queryOverride := ""
	if len(args) > 1 {
		queryOverride = strings.Join(args[1:], " ")
	}
	src, err := ResolveSelection(args[0], queryOverride)
	if err != nil {
		return err
	}

	if src.Location.Provider != "rb" || src.Location.Resource != "tracks" {
		return fmt.Errorf("currently only rb/tracks is supported as a source for add")
	}

	var trackIDs []string
	for _, t := range src.Tracks {
		trackIDs = append(trackIDs, t.ID)
	}

	// 2. Resolve targets and apply
	var targetNames []string
	for _, targetStr := range addTargets {
		tgt, err := ResolveSelection(targetStr, "")
		if err != nil {
			return err
		}
		if tgt.Location.Provider != "rb" || tgt.Location.Resource != "playlists" {
			return fmt.Errorf("currently only rb/playlists is supported as a target for add, got %q", targetStr)
		}
		for _, n := range tgt.Nodes {
			targetNames = append(targetNames, n.Name)
		}
	}

	RunBulkOperation("add", targetNames, trackIDs, func(targetName string, items []string) (bool, int) {
		return syncEng.AddTracksToPlaylist(targetName, items)
	})

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
