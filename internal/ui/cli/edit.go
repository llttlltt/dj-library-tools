package cli

import (
	"fmt"
	"strings"

	"github.com/llttlltt/dj-library-tools/internal/services/resolver"
	"github.com/spf13/cobra"
)

func newEditCmd() *cobra.Command {
	var setFields []string
	var filterMissing bool
	var filterExists bool

	cmd := &cobra.Command{
		Use:   "edit [selection] [query]",
		Short: "Update metadata for resources",
		Long: `Modify metadata fields for tracks or other resources.
For library maintenance (deduplication, path repair), use 'djlt fix'.

Examples:
  # Set a comment for tracks
  djlt edit rb/tracks playlists:Inbox --set comment:Great`,
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

			sel, prov, err := resolver.ResolveSelection(args[0], queryOverride, opts)
			if err != nil {
				return err
			}

			ctx := getExecContext()

			// Handle Metadata Updates
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
					return prov.System().Save(ctx, "")
				} else {
					fmt.Println("Run with --apply to persist changes.")
					return nil
				}
			}

			return cmd.Help()
		},
	}

	cmd.Flags().StringSliceVar(&setFields, "set", []string{}, "Metadata fields to update (key:value)")
	cmd.Flags().BoolVar(&filterMissing, "missing", false, "Filter for tracks where the physical file is missing")
	cmd.Flags().BoolVar(&filterExists, "exists", false, "Filter for tracks where the physical file exists")

	return cmd
}
