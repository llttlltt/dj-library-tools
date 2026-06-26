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
	removeOrigins []string
)

var removeCmd = &cobra.Command{
	Use:   "remove [source-resource] [source-query] --from [origin-resource] [origin-query]",
	Short: "Remove items from one or more origins",
	Long: `Remove items matching a source selection from one or more origin resources.
Currently supports removing tracks (rb/tracks) from playlists (rb/playlists).

Example:
  djlt remove rb/tracks artist:Four --from "rb/playlists name:Inbox"`,
	Args: cobra.MinimumNArgs(1),
	RunE: runRemoveCmd,
}

func runRemoveCmd(cmd *cobra.Command, args []string) error {
	if len(removeOrigins) == 0 {
		return fmt.Errorf("at least one --from origin is required")
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

	// 2. Resolve origins and apply
	var targetNodes []models.Node
	var targetProv provider.WritableProvider

	for _, originStr := range removeOrigins {
		org, err := ResolveSelection(originStr, "")
		if err != nil {
			return err
		}

		wp, ok := org.Provider.(provider.WritableProvider)
		if !ok {
			return fmt.Errorf("provider %q does not support removing tracks", org.Location.Provider)
		}
		targetProv = wp

		if org.Location.Resource != "playlists" {
			return fmt.Errorf("can only remove from playlists, got %q", org.Location.Resource)
		}
		targetNodes = append(targetNodes, org.Nodes...)
	}

	if dryRun {
		for _, n := range targetNodes {
			fmt.Printf("[Dry Run] Would remove %d tracks from playlist %q\n", len(src.Tracks), n.Name)
		}
		return nil
	}

	for _, n := range targetNodes {
		removed, err := targetProv.RemoveTracks(n, src.Tracks)
		if err != nil {
			return err
		}
		fmt.Printf("Removed %d tracks from %q\n", removed, n.Name)
	}

	// For Rekordbox we still need to save.
	if rb, ok := targetProv.(*provider.RekordboxProvider); ok {
		_, path, _ := loadXMLFunc()
		return rb.Engine.Library.(engine.WritableLibrary).Save(path)
	}

	return nil
}

func init() {
	removeCmd.Flags().StringSliceVar(&removeOrigins, "from", []string{}, "Origin resource(s) to remove from (repeatable)")
	RootCmd.AddCommand(removeCmd)
}
