package cli

import (
	"fmt"
	"strings"

	"github.com/llttlltt/dj-library-tools/internal/engine"
	"github.com/llttlltt/dj-library-tools/internal/models"
	"github.com/llttlltt/dj-library-tools/internal/provider"
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

	// 1. Resolve source
	queryOverride := ""
	if len(args) > 1 {
		queryOverride = strings.Join(args[1:], " ")
	}
	src, err := ResolveSelection(args[0], queryOverride)
	if err != nil {
		return err
	}

	// 2. Resolve targets and apply
	var targetNodes []models.Node
	var targetProv provider.WritableProvider

	for _, targetStr := range addTargets {
		tgt, err := ResolveSelection(targetStr, "")
		if err != nil {
			return err
		}

		wp, ok := tgt.Provider.(provider.WritableProvider)
		if !ok {
			return fmt.Errorf("provider %q does not support adding tracks", tgt.Location.Provider)
		}
		targetProv = wp

		if tgt.Location.Resource != "playlists" {
			return fmt.Errorf("can only add to playlists, got %q", tgt.Location.Resource)
		}
		targetNodes = append(targetNodes, tgt.Nodes...)
	}

	if dryRun {
		for _, n := range targetNodes {
			fmt.Printf("[Dry Run] Would add %d tracks to playlist %q\n", len(src.Tracks), n.Name)
		}
		return nil
	}

	for _, n := range targetNodes {
		added, err := targetProv.AddTracks(n, src.Tracks)
		if err != nil {
			return err
		}
		fmt.Printf("Added %d tracks to %q\n", added, n.Name)
	}

	// For Rekordbox we still need to save.
	if rb, ok := targetProv.(*provider.RekordboxProvider); ok {
		_, path, _ := loadXMLFunc()
		return rb.Engine.Library.(engine.WritableLibrary).Save(path)
	}

	return nil
}

func init() {
	addCmd.Flags().StringSliceVar(&addTargets, "to", []string{}, "Target resource(s) to add to (repeatable)")
	addCmd.Flags().BoolVar(&addForce, "force", false, "Allow adding duplicates (if supported by target)")

	RootCmd.AddCommand(addCmd)
}
