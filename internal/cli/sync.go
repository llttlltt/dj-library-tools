package cli

import (
	"fmt"
	"strings"

	"github.com/llttlltt/dj-library-tools/internal/provider"
	"github.com/llttlltt/dj-library-tools/internal/sync"
	"github.com/spf13/cobra"
)

func newSyncCmd() *cobra.Command {
	var syncTo []string
	var exportDest, exportFormat string
	var syncAppend bool
	var metadataFields, matchFields []string

	cmd := &cobra.Command{
		Use:   "sync [source-resource] [source-query] --to [target-resource] [target-query]",
		Short: "Keep a playlist or metadata in sync with a track query",
		Long: `Synchronizes a target (like a Rekordbox playlist or M3U file) with a source query.

The sync command is "surgical"—it only adds or removes tracks necessary to make the target
match the source. By default, it removes tracks from the target that no longer match the query.

### Metadata Reconciliation
If --metadata is specified, djlt will match tracks between the source and target using the --match keys
and synchronize specific metadata fields (e.g. beatgrids, rating).

### Examples

**Keep an "Inbox" playlist matched to specific criteria:**
  djlt sync "rb/tracks added:>today" --to "rb/playlists name:Inbox"

**Sync beatgrids from a backup Rekordbox XML to your primary library:**
  djlt sync "rb/tracks" --file backup.xml --to "rb/tracks" --metadata beatgrids

**Sync ratings from Plex to Rekordbox matching by filename:**
  djlt sync "plex/tracks" --to "rb/tracks" --metadata rating --match filename
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
				return HandleError(err)
			}

			for _, targetStr := range syncTo {
				tgt, err := ResolveSelection(targetStr, "")
				if err != nil {
					return HandleError(err)
				}

				prov := tgt.Provider

				if dryRun {
					action := "sync"
					if syncAppend {
						action = "append to"
					}
					fmt.Printf("[Dry Run] Would %s playlist %q with %d tracks\n", action, tgt.Location.Query, len(src.Tracks))
					if len(metadataFields) > 0 {
						fmt.Printf("[Dry Run] Would sync metadata fields %v for matched tracks\n", metadataFields)
					}
					continue
				}

				// 1. Membership Sync
				err = prov.System().Sync(getExecContext(), src.Tracks, src.Location.Query, tgt.Location.Query, provider.SyncOptions{
					ExportDest:   exportDest,
					ExportFormat: exportFormat,
					AppendOnly:   syncAppend,
				})
				if err != nil {
					return HandleError(err)
				}

				// 2. Metadata Sync (if requested)
				if len(metadataFields) > 0 {
					// We need to resolve the library of the target provider to perform the join
					// This is a slight leak - let's see how we can make Join more agnostic.
					// For now, let's allow it if we can get target tracks.
					targetTracks, err := prov.Tracks().List(getExecContext(), "")
					if err == nil {
						// Match datasets
						matcher := sync.NewMatcher(targetTracks).WithKeys(matchFields)
						matches := sync.NewOrchestrator(nil, dryRun, verbose).WithMatcher(matcher).Join(src.Tracks, matchFields)
						
						if err := prov.Tracks().UpdateBatch(getExecContext(), matches, metadataFields); err != nil {
							return HandleError(err)
						}
					}
				}
			}
			return nil
		},
	}
	cmd.Flags().StringSliceVar(&syncTo, "to", []string{}, "Target resource(s) to sync to (repeatable)")
	cmd.Flags().StringVar(&exportDest, "dest", "", "Destination directory for exported files")
	cmd.Flags().StringVar(&exportFormat, "format", "mp3", "Target format for exported files")
	cmd.Flags().BoolVar(&syncAppend, "append", false, "Append new tracks without removing existing ones")
	cmd.Flags().StringSliceVar(&metadataFields, "metadata", []string{}, "Metadata fields to synchronize (e.g. beatgrids, rating)")
	cmd.Flags().StringSliceVar(&matchFields, "match", []string{"artist", "title"}, "Fields to use for matching tracks")
	return cmd
}
