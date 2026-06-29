package cli

import (
	"fmt"
	"strings"

	"github.com/llttlltt/dj-library-tools/internal/provider"
	"github.com/llttlltt/dj-library-tools/internal/resolver"
	"github.com/spf13/cobra"
)

func newFixCmd() *cobra.Command {
	var (
		duplicates []string
		metadata   []string
		paths      []string
		orphans    []string
	)

	cmd := &cobra.Command{
		Use:   "fix [selection] [query]",
		Short: "Perform health, formatting, and structural repairs on the library",
		Long: `A multi-purpose repair command for library maintenance.

Examples:
  # Remove duplicate tracks from specific playlists
  djlt fix rb/playlists "Inbox,Recently Added" --duplicates members

  # Normalize metadata for matching tracks
  djlt fix rb/tracks "genre:Techno" --metadata artist,album

  # Repair broken file paths for missing tracks
  djlt fix rb/tracks --missing --paths relocate`,
		Args: cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			queryOverride := ""
			if len(args) > 1 {
				queryOverride = strings.Join(args[1:], " ")
			}

			opts := resolver.ResolveOptions{
				FilePath: filePath,
				Apply:    apply,
				Verbose:  verbose,
			}

			sel, err := resolver.ResolveSelection(args[0], queryOverride, opts)
			if err != nil {
				return err
			}

			fixOpts := provider.FixOptions{
				Actions: make(map[provider.FixType][]string),
			}

			if len(duplicates) > 0 {
				fixOpts.Actions[provider.FixDuplicates] = duplicates
			}
			if len(metadata) > 0 {
				fixOpts.Actions[provider.FixMetadata] = metadata
			}
			if len(paths) > 0 {
				fixOpts.Actions[provider.FixPaths] = paths
			}
			if len(orphans) > 0 {
				fixOpts.Actions[provider.FixOrphans] = orphans
			}

			if len(fixOpts.Actions) == 0 {
				return fmt.Errorf("no repair actions specified; use --duplicates, --metadata, --paths, or --orphans")
			}

			prov := sel.Provider
			ctx := getExecContext()

			if _, err := prov.System().Fix(ctx, *sel, fixOpts); err != nil {
				return err
			}

			if ctx.Apply {
				fmt.Printf("Successfully performed repairs.\n")
				return prov.System().Save(ctx, "")
			} else {
				fmt.Printf("Run with --apply to persist changes.\n")
			}

			return nil
		},
	}

	cmd.Flags().StringSliceVar(&duplicates, "duplicates", []string{}, "Remove duplicates (targets: members, tracks)")
	cmd.Flags().StringSliceVar(&metadata, "metadata", []string{}, "Fix/normalize metadata (targets: artist, album, etc.)")
	cmd.Flags().StringSliceVar(&paths, "paths", []string{}, "Repair file paths (targets: relocate, normalize)")
	cmd.Flags().StringSliceVar(&orphans, "orphans", []string{}, "Remove orphaned resources (targets: all)")

	return cmd
}
