package cli

import (
	"fmt"
	"strings"

	"github.com/llttlltt/dj-library-tools/internal/models"
	"github.com/llttlltt/dj-library-tools/internal/provider"
	"github.com/spf13/cobra"
)

func newEditCmd() *cobra.Command {
	var setFields []string
	var relocateDir string
	var matchFields []string
	var repair bool

	cmd := &cobra.Command{
		Use:   "edit [selection] [query]",
		Short: "Update metadata, repair paths, or fix library issues",
		Long: `A unified command for modifying resource state.

Examples:
  # Set a comment for tracks
  djlt edit rb/tracks playlists:Inbox --set comment:Great

  # Relocate missing files
  djlt edit rb/tracks --missing --relocate "/Volumes/Media/Music"

  # Run provider-specific repairs (formerly 'fix')
  djlt edit rb/tracks --repair`,
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

			// 1. Handle Repairs
			if repair {
				if dryRun {
					fmt.Printf("[Dry Run] Would perform repair on %s/%s\n", sel.Location.Provider, sel.Location.Resource)
					return nil
				}
				if err := wp.Fix(ctx, sel.Location.Resource, sel.Location.Query); err != nil {
					return err
				}
				fmt.Println("Repair completed successfully.")
				return wp.Save(ctx, "")
			}

			// 2. Handle Relocation
			if relocateDir != "" {
				// We keep the relocation logic here as it's a cross-provider 'Search & Patch' orchestration
				// but it calls ModifyTracks on the provider for the actual write.
				relocated := wp.(interface {
					Relocate(tracks []models.Track, dir string, match []string) map[string]string
				}).Relocate(sel.Tracks, relocateDir, matchFields)
				
				if len(relocated) == 0 {
					fmt.Println("No tracks were relocated.")
					return nil
				}

				if dryRun {
					fmt.Printf("[Dry Run] Would update paths for %d tracks\n", len(relocated))
					return nil
				}

				for id, newPath := range relocated {
					changes := map[string]string{"location": newPath}
					wp.UpdateTracks(ctx, "id:"+id, changes)
				}
				
				return wp.Save(ctx, "")
			}

			// 3. Handle Metadata Updates
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

				count, err := wp.UpdateTracks(ctx, sel.Location.Query, changes)
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
	cmd.Flags().BoolVar(&repair, "repair", false, "Perform provider-specific health/formatting repairs")

	return cmd
}
