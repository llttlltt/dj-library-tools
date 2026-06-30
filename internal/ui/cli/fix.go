package cli

import (
	"fmt"
	"strings"

	"github.com/llttlltt/dj-library-tools/internal/services/orchestrator"
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
  djlt fix rb/tracks "missing:true" --paths normalize`,
		Args: cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			queryOverride := ""
			if len(args) > 1 {
				queryOverride = strings.Join(args[1:], " ")
			}

			orch := getOrchestrator()
			runOpts := getRunOptions()

			fixOpts := orchestrator.FixOptions{
				Actions: make(map[orchestrator.FixKind][]string),
			}

			if len(duplicates) > 0 {
				fixOpts.Actions[orchestrator.FixDuplicates] = duplicates
			}
			if len(metadata) > 0 {
				fixOpts.Actions[orchestrator.FixMetadata] = metadata
			}
			if len(paths) > 0 {
				fixOpts.Actions[orchestrator.FixPaths] = paths
			}
			if len(orphans) > 0 {
				fixOpts.Actions[orchestrator.FixOrphans] = orphans
			}

			if len(fixOpts.Actions) == 0 {
				return fmt.Errorf("no repair actions specified; use --duplicates, --metadata, --paths, or --orphans")
			}

			count, err := orch.Fix(cmd.Context(), args[0], queryOverride, runOpts, fixOpts)
			if err != nil {
				return HandleError(err)
			}

			if apply {
				fmt.Printf("✔ Repaired %d item(s).\n", count)
			} else {
				fmt.Printf("Scope: %d item(s) would be repaired.\n", count)
				fmt.Println("Run with --apply to commit.")
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
