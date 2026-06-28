package cli

import (
	"fmt"
	"strings"

	"github.com/llttlltt/dj-library-tools/internal/provider"
	"github.com/spf13/cobra"
)

func newSyncCmd() *cobra.Command {
	var syncTo []string
	var exportDest, exportFormat string
	var syncAppend bool

	cmd := &cobra.Command{
		Use:   "sync [source-resource] [source-query] --to [target-resource] [target-query]",
		Short: "Keep a playlist in sync with a track query",
		Long: `Synchronizes a target (like a Rekordbox playlist or M3U file) with a source query.

The sync command is "surgical"—it only adds or removes tracks necessary to make the target
match the source. By default, it removes tracks from the target that no longer match the query.

### Examples

**Keep an "Inbox" playlist matched to specific criteria:**
  djlt sync "rb/tracks added:>today" --to "rb/playlists name:Inbox"

**Add new tracks to a playlist without removing existing ones:**
  djlt sync "rb/tracks rating:5" --to "rb/playlists name:Favorites" --append

**Sync a query to an external M3U playlist file:**
  djlt sync "rb/tracks genre:House" --to "m3u/path/to/playlist.m3u"
`,
		Args: cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(syncTo) == 0 {
				return fmt.Errorf("at least one --to target is required")
			}

			queryOverride := ""
			if len(args) > 1 {
				queryOverride = strings.Join(args[1:], " ")
			}
			src, err := ResolveSelection(args[0], queryOverride)
			if err != nil {
				return err
			}

			for _, targetStr := range syncTo {
				tgt, err := ResolveSelection(targetStr, "")
				if err != nil {
					return err
				}

				wp, ok := tgt.Provider.(provider.WritableProvider)
				if !ok {
					return fmt.Errorf("unsupported sync target: %s (provider is read-only)", tgt.Location.Provider)
				}

				if dryRun {
					action := "sync"
					if syncAppend {
						action = "append to"
					}
					fmt.Printf("[Dry Run] Would %s playlist %q with %d tracks\n", action, tgt.Location.Query, len(src.Tracks))
					continue
				}

				err = wp.Sync(getExecContext(), src.Tracks, src.Location.Query, tgt.Location.Query, provider.SyncOptions{
					ExportDest:   exportDest,
					ExportFormat: exportFormat,
					AppendOnly:   syncAppend,
				})
				if err != nil {
					return err
				}
			}
			return nil
		},
	}
	cmd.Flags().StringSliceVar(&syncTo, "to", []string{}, "Target resource(s) to sync to (repeatable)")
	cmd.Flags().StringVar(&exportDest, "dest", "", "Destination directory for exported files")
	cmd.Flags().StringVar(&exportFormat, "format", "mp3", "Target format for exported files")
	cmd.Flags().BoolVar(&syncAppend, "append", false, "Append new tracks without removing existing ones")
	return cmd
}
