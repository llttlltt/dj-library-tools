package cli

import (
	"fmt"
	"strings"

	"github.com/llttlltt/dj-library-tools/internal/models"
	"github.com/llttlltt/dj-library-tools/internal/provider"
	"github.com/llttlltt/dj-library-tools/internal/resolver"
	"github.com/spf13/cobra"
)

func newEditCmd() *cobra.Command {
	var setFields []string
	var relocateDir string
	var matchFields []string
	var repair bool
	var filterMissing bool
	var filterExists bool

	cmd := &cobra.Command{
		Use:   "edit [selection] [query]",
		Short: "Update metadata, repair paths, or fix library issues",
		Long: `A unified command for modifying resource state.

Examples:
  # Set a comment for tracks
  djlt edit rb/tracks playlists:Inbox --set comment:Great

  # Relocate missing files
  djlt edit rb/tracks --missing --relocate "/Volumes/Media/Music"`,
		Args: cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			queryOverride := ""
			if len(args) > 1 {
				queryOverride = strings.Join(args[1:], " ")
			}

			opts := resolver.ResolveOptions{
				FilePath:      filePath,
				FilterMissing: filterMissing,
				FilterExists:  filterExists,
				Apply:         apply,
				Verbose:       verbose,
			}

			sel, err := resolver.ResolveSelection(args[0], queryOverride, opts)
			if err != nil {
				return err
			}

			prov := sel.Provider
			ctx := getExecContext()

			// 1. Handle Repairs
			if repair {
				fixOpts := provider.FixOptions{
					Actions: map[provider.FixType][]string{
						provider.FixPaths: {"relocate"},
					},
				}
				if _, err := prov.System().Fix(ctx, *sel, fixOpts); err != nil {
					return err
				}
				if ctx.Apply {
					fmt.Println("Repair completed successfully.")
					return prov.System().Save(ctx, "")
				}
				return nil
			}

			// 2. Handle Relocation
			if relocateDir != "" {
				// Redirect to 'fix --paths relocate' logic if preferred, 
				// but for now keeping compatible with existing 'edit' workflow.
				relocated := prov.(interface {
					Relocate(tracks []models.Track, dir string, match []string) map[string]string
				}).Relocate(sel.Tracks, relocateDir, matchFields)
				
				if len(relocated) == 0 {
					fmt.Println("No tracks were relocated.")
					return nil
				}

				for id, newPath := range relocated {
					changes := map[string]string{"location": newPath}
					prov.Tracks().Update(ctx, "id:"+id, changes)
				}
				
				return prov.System().Save(ctx, "")
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

				count, err := prov.Tracks().Update(ctx, sel.Location.Query, changes)
				if err != nil {
					return err
				}

				if ctx.Apply {
					fmt.Printf("Successfully modified %d tracks.\n", count)
				}
				return prov.System().Save(ctx, "")
			}

			return cmd.Help()
		},
	}

	cmd.Flags().StringSliceVar(&setFields, "set", []string{}, "Metadata fields to update (key:value)")
	cmd.Flags().StringVar(&relocateDir, "relocate", "", "Search this directory to repair missing file paths")
	cmd.Flags().StringSliceVar(&matchFields, "match", []string{"filename"}, "Criteria to use for relocation matching")
	cmd.Flags().BoolVar(&repair, "repair", false, "Perform provider-specific health/formatting repairs")
	cmd.Flags().BoolVar(&filterMissing, "missing", false, "Filter for tracks where the physical file is missing")
	cmd.Flags().BoolVar(&filterExists, "exists", false, "Filter for tracks where the physical file exists")

	cmd.Flags().MarkHidden("repair")
	cmd.Flags().MarkHidden("relocate")

	return cmd
}
