package cli

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/llttlltt/dj-library-tools/internal/engine"
	syncpkg "github.com/llttlltt/dj-library-tools/internal/sync"
	"github.com/llttlltt/dj-library-tools/internal/utils"
	"github.com/spf13/cobra"
	"github.com/vbauerster/mpb/v8"
	"github.com/vbauerster/mpb/v8/decor"
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
	syncEng := syncpkg.NewEngine(nil, rbXML)

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

		if dryRun {
			for _, origin := range origins {
				fmt.Printf("[Dry Run] Would remove %d tracks from playlist %q\n", len(trackIDs), origin.Node.Name)
			}
			continue
		}

		p := mpb.New(mpb.WithWidth(64))
		for _, origin := range origins {
			bar := p.AddBar(int64(len(trackIDs)),
				mpb.PrependDecorators(
					decor.Name(fmt.Sprintf("Removing from %q", origin.Node.Name), decor.WCSyncSpaceR),
					decor.CountersNoUnit("%d / %d", decor.WCSyncSpace),
				),
				mpb.AppendDecorators(decor.Percentage(decor.WCSyncSpace)),
			)

			if verbose {
				fmt.Printf("Removing %d tracks from playlist %q...\n", len(trackIDs), origin.Node.Name)
			}

			// We process in chunks to show progress
			chunkSize := 10
			if len(trackIDs) < chunkSize {
				chunkSize = len(trackIDs)
			}

			totalRemoved := 0
			for i := 0; i < len(trackIDs); i += chunkSize {
				end := i + chunkSize
				if end > len(trackIDs) {
					end = len(trackIDs)
				}
				chunk := trackIDs[i:end]
				found, removed := syncEng.RemoveTracksFromPlaylist(origin.Node.Name, chunk)
				if !found {
					fmt.Printf("Warning: playlist %q not found during remove\n", origin.Node.Name)
					bar.Abort(false)
					break
				}
				totalRemoved += removed
				bar.IncrBy(len(chunk))
			}
			p.Wait()
			fmt.Printf("Removed %d tracks from %q\n", totalRemoved, origin.Node.Name)
		}
	}

	if dryRun {
		return nil
	}

	return syncEng.SaveXML(path)
}

func init() {
	removeCmd.Flags().StringSliceVar(&removeOrigins, "from", []string{}, "Origin resource(s) to remove from (repeatable)")
	RootCmd.AddCommand(removeCmd)
}
