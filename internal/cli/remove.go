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

	eng := engine.NewEngine(engine.NewRekordboxLibrary(rbXML))
	syncEng := syncpkg.NewEngine(nil, engine.NewRekordboxLibrary(rbXML))

	// 1. Resolve source
	sourceQuery := ""
	if len(args) > 1 {
		sourceQuery = strings.Join(args[1:], " ")
	}
	src := utils.ParseLocation(args[0], sourceQuery)

	if src.Provider != "rb" || src.Resource != "tracks" {
		return fmt.Errorf("currently only rb/tracks is supported as a source for remove")
	}

	tracks, err := eng.Ls(src.Query)
	if err != nil {
		return fmt.Errorf("failed to resolve source tracks: %w", err)
	}
	if len(tracks) == 0 {
		fmt.Println("No tracks found matching query.")
		return nil
	}

	var trackIDs []string
	for _, t := range tracks {
		trackIDs = append(trackIDs, strconv.Itoa(t.TrackID))
	}

	// 2. Resolve origins and apply
	var originNames []string
	for _, originStr := range removeOrigins {
		org := utils.ParseLocation(originStr, "")
		if org.Provider != "rb" || org.Resource != "playlists" {
			return fmt.Errorf("currently only rb/playlists is supported as an origin for remove, got %q", originStr)
		}

		origins, err := eng.LsPlaylists(org.Query)
		if err != nil {
			return fmt.Errorf("failed to resolve origin playlists: %w", err)
		}
		if len(origins) == 0 {
			return fmt.Errorf("no origin playlists matched query %q", org.Query)
		}
		for _, o := range origins {
			originNames = append(originNames, o.Node.Name)
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
