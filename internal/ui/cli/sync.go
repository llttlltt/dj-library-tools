package cli

import (
	"fmt"
	"strings"

	"github.com/llttlltt/dj-library-tools/internal/core/models"
	"github.com/llttlltt/dj-library-tools/internal/providers"
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

			src, _, err := ResolveSelection(args[0], queryOverride)
			if err != nil {
				return HandleError(err)
			}

			for _, targetStr := range syncTo {
				tgt, prov, err := ResolveSelection(targetStr, "")
				if err != nil {
					return HandleError(err)
				}
				resolvedTargetID := tgt.Location.Query
				if len(tgt.Groups) > 0 {
					resolvedTargetID = tgt.Groups[0].ID
				}

				// Perform Diffing for CLI Feedback
				if verbose || !apply {
					targetName := tgt.Location.Query
					if targetName == "" {
						targetName = tgt.Location.Resource
					}
					// If the target is a group (playlist/folder), we can get its current members
					if len(tgt.Groups) > 0 {
						group := tgt.Groups[0]
						targetName = group.Name
						
						// Get current IDs by checking all tracks for membership in this SPECIFIC group
						var currentIDs []string
						allTracks, _ := prov.Tracks().List(getExecContext(), "")
						for _, t := range allTracks {
							for _, m := range t.Playlists {
								// Match by name AND folder to ensure we have the right one
								if m.Name == group.Name && m.Folder == group.ParentFolder {
									currentIDs = append(currentIDs, t.ID)
									break
								}
							}
						}

						// Calculate Source IDs
						var sourceIDs []string
						for _, t := range src.Tracks {
							sourceIDs = append(sourceIDs, t.ID)
						}

						added, removed := calculateDiff(currentIDs, sourceIDs)
						
						if verbose {
							trackLookup := make(map[string]models.Track)
							allTracks, _ := prov.Tracks().List(getExecContext(), "")
							for _, t := range allTracks {
								trackLookup[t.ID] = t
							}
							if len(added) > 0 {
								fmt.Printf("\nTracks to ADD:\n")
								printTrackTable(added, trackLookup, getExecContext().Feedback)
							}
							if len(removed) > 0 && !syncAppend {
								fmt.Printf("\nTracks to REMOVE:\n")
								printTrackTable(removed, trackLookup, getExecContext().Feedback)
							}
							fmt.Println()
						}

						fmt.Printf("%s:\n", targetName)
						if syncAppend {
							fmt.Printf("- Total tracks to add: %d\n", len(added))
						} else {
							fmt.Printf("- Current tracks:      %d\n", len(currentIDs))
							fmt.Printf("- Tracks to add:       %d\n", len(added))
							fmt.Printf("- Tracks to remove:    %d\n", len(removed))
							fmt.Printf("- Final count:         %d\n", len(sourceIDs))
						}
						fmt.Println()

						if apply {
							fmt.Printf("Successfully synchronized %q.\n", targetName)
						} else {
							fmt.Printf("Run with --apply to persist changes.\n")
						}
					}
				}

				// Perform Membership and Metadata Sync
				err = prov.System().Sync(getExecContext(), src.Tracks, resolvedTargetID, provider.SyncOptions{
					ExportDest:     exportDest,
					ExportFormat:   exportFormat,
					AppendOnly:     syncAppend,
					MetadataFields: metadataFields,
					MatchFields:    matchFields,
				})
				if err != nil {
					return HandleError(err)
				}

				if apply {
					if err := prov.System().Save(getExecContext(), toFilePath); err != nil {
						return HandleError(err)
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
	cmd.Flags().StringVar(&toFilePath, "to-file", "", "Path to the destination library file for sync/move operations")
	return cmd
}

func calculateDiff(current, target []string) (added, removed []string) {
	currentMap := make(map[string]bool)
	for _, id := range current {
		currentMap[id] = true
	}
	targetMap := make(map[string]bool)
	for _, id := range target {
		targetMap[id] = true
	}

	for _, id := range target {
		if !currentMap[id] {
			added = append(added, id)
		}
	}
	for _, id := range current {
		if !targetMap[id] {
			removed = append(removed, id)
		}
	}
	return
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
