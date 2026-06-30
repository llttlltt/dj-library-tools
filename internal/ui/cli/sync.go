package cli

import (
	"fmt"
	"strings"

	"github.com/llttlltt/dj-library-tools/internal/core/models"
	"github.com/llttlltt/dj-library-tools/internal/providers"
	"github.com/llttlltt/dj-library-tools/internal/services/orchestrator"
	"github.com/spf13/cobra"
)

func newSyncCmd() *cobra.Command {
	var syncTo []string
	var exportDest, exportFormat string
	var syncAppend bool
	var metadataFields, matchFields []string
	var toFilePath string

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

			orch := getOrchestrator()
			runOpts := getRunOptions()

			for _, targetStr := range syncTo {
				// Perform Diffing for CLI Feedback
				if verbose || !apply {
					diff, err := orch.GetSyncDiff(cmd.Context(), args[0], targetStr, queryOverride, runOpts, syncAppend)
					if err != nil {
						return HandleError(err)
					}

					if verbose {
						if len(diff.AddedIDs) > 0 {
							fmt.Printf("\nTracks to ADD:\n")
							printTrackTable(diff.AddedIDs, diff.TrackLookup, &TerminalFeedback{})
						}
						if len(diff.RemovedIDs) > 0 && !syncAppend {
							fmt.Printf("\nTracks to REMOVE:\n")
							printTrackTable(diff.RemovedIDs, diff.TrackLookup, &TerminalFeedback{})
						}
						fmt.Println()
					}

					fmt.Printf("%s:\n", diff.TargetName)
					if syncAppend {
						fmt.Printf("- Total tracks to add: %d\n", len(diff.AddedIDs))
					} else {
						fmt.Printf("- Current tracks:      %d\n", len(diff.CurrentIDs))
						fmt.Printf("- Tracks to add:       %d\n", len(diff.AddedIDs))
						fmt.Printf("- Tracks to remove:    %d\n", len(diff.RemovedIDs))
						fmt.Printf("- Final count:         %d\n", len(diff.SourceIDs))
					}
					fmt.Println()

					if apply {
						fmt.Printf("Successfully synchronized %q.\n", diff.TargetName)
					} else {
						fmt.Printf("Run with --apply to persist changes.\n")
					}
				}

				// Perform Membership and Metadata Sync
				err := orch.Sync(cmd.Context(), args[0], targetStr, queryOverride, runOpts, orchestrator.SyncOptions{
					ExportDest:     exportDest,
					ExportFormat:   exportFormat,
					AppendOnly:     syncAppend,
					MetadataFields: metadataFields,
					MatchFields:    matchFields,
				})
				if err != nil {
					return HandleError(err)
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
	cmd.Flags().StringVar(&toFilePath, "to-file", "", "Path to the destination library file for sync/move operations")
	return cmd
}

func printTrackTable(ids []string, lookup map[string]models.Track, feedback provider.Feedback) {
	headers := []string{"id", "title", "artist"}
	var rows [][]string
	for _, id := range ids {
		if t, ok := lookup[id]; ok {
			rows = append(rows, []string{id, t.Title, t.Artist})
		} else {
			rows = append(rows, []string{id, "[Unknown]", "[Unknown]"})
		}
	}
	feedback.OnTable(headers, rows)
}
