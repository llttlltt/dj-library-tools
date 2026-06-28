package cli

import (
	"fmt"
	"strings"

	"github.com/llttlltt/dj-library-tools/internal/provider"
	"github.com/llttlltt/dj-library-tools/internal/sync"
	"github.com/spf13/cobra"
)

func newModifyCmd() *cobra.Command {
	var setFields []string
	var relocateDir string
	var matchFields []string

	cmd := &cobra.Command{
		Use:   "modify [selection] [query]",
		Short: "Bulk update track metadata or relocate missing files",
		Long: `Modify track metadata or repair broken file paths.

Examples:
  # Set a comment for all tracks in a playlist
  djlt modify rb/tracks playlists:Inbox --set comment:Great

  # Relocate missing files by searching a directory
  djlt modify rb/tracks --missing --relocate "/Volumes/Media/Music" --match filename`,
		Args: cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			queryOverride := ""
			if len(args) > 1 {
				queryOverride = strings.Join(args[1:], " ")
			}

			sel, err := ResolveSelection(args[0], queryOverride)
			if err != nil {
				return err
			}

			wp, ok := sel.Provider.(provider.WritableProvider)
			if !ok {
				return fmt.Errorf("provider %q is read-only", sel.Location.Provider)
			}

			ctx := getExecContext()

			// 1. Handle Relocation
			if relocateDir != "" {
				orch := sync.NewOrchestrator(nil, dryRun, verbose)
				relocated := orch.Relocate(sel.Tracks, relocateDir, matchFields)
				
				if len(relocated) == 0 {
					fmt.Println("No tracks were relocated.")
					return nil
				}

				if dryRun {
					fmt.Printf("[Dry Run] Would update paths for %d tracks\n", len(relocated))
					return nil
				}

				changes := make(map[string]string)
				for id, newPath := range relocated {
					changes["location"] = newPath
					_, err := wp.ModifyTracks(ctx, "id:"+id, changes)
					if err != nil {
						fmt.Printf("Warning: failed to update path for track %s: %v\n", id, err)
					}
				}
				
				return wp.Save(ctx, "")
			}

			// 2. Handle Metadata Updates
			if len(setFields) > 0 {
				changes := make(map[string]string)
				for _, f := range setFields {
					parts := strings.SplitN(f, ":", 2)
					if len(parts) != 2 {
						return fmt.Errorf("invalid set format %q; use key:value", f)
					}
					changes[parts[0]] = parts[1]
				}

				if dryRun {
					fmt.Printf("[Dry Run] Would apply changes %v to %d tracks\n", changes, len(sel.Tracks))
					return nil
				}

				count, err := wp.ModifyTracks(ctx, sel.Location.Query, changes)
				if err != nil {
					return err
				}

				fmt.Printf("Successfully modified %d tracks.\n", count)
				return wp.Save(ctx, "")
			}

			return cmd.Help()
		},
	}

	cmd.Flags().StringSliceVar(&setFields, "set", []string{}, "Metadata fields to update (key:value)")
	cmd.Flags().StringVar(&relocateDir, "relocate", "", "Search this directory to repair missing file paths")
	cmd.Flags().StringSliceVar(&matchFields, "match", []string{"filename"}, "Criteria to use for relocation matching")

	return cmd
}
