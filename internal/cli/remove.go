package cli

import (
	"fmt"
	"strings"

	"github.com/llttlltt/dj-library-tools/internal/engine"
	syncpkg "github.com/llttlltt/dj-library-tools/internal/sync"
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
		return fmt.Errorf("currently only rb/tracks is supported as a source for remove")
	}

	var trackIDs []string
	for _, t := range src.Tracks {
		trackIDs = append(trackIDs, t.ID)
	}

	// 2. Resolve origins and apply
	var originNames []string
	for _, originStr := range removeOrigins {
		org, err := ResolveSelection(originStr, "")
		if err != nil {
			return err
		}
		if org.Location.Provider != "rb" || org.Location.Resource != "playlists" {
			return fmt.Errorf("currently only rb/playlists is supported as an origin for remove, got %q", originStr)
		}
		for _, o := range org.Nodes {
			originNames = append(originNames, o.Name)
		}
	}

	RunBulkOperation("remove", originNames, trackIDs, func(targetName string, items []string) (bool, int) {
		return syncEng.RemoveTracksFromPlaylist(targetName, items)
	})

	if dryRun {
		return nil
	}

	return syncEng.SaveXML(path)
}

func init() {
	removeCmd.Flags().StringSliceVar(&removeOrigins, "from", []string{}, "Origin resource(s) to remove from (repeatable)")
	RootCmd.AddCommand(removeCmd)
}
