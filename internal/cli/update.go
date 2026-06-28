package cli

import (
	"fmt"
	"github.com/llttlltt/dj-library-tools/internal/provider"
	"github.com/llttlltt/dj-library-tools/internal/sync"
	"github.com/spf13/cobra"
)

func newUpdateCmd() *cobra.Command {
	var updateFrom string
	var matchFields []string

	cmd := &cobra.Command{
		Use:   "update [target-selection] --from [source-selection]",
		Short: "Reconcile metadata between two libraries",
		Long: `Matches tracks between a source and target selection and synchronizes metadata.

Example:
  djlt update rb/tracks --from plex/tracks --metadata beatgrids`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if updateFrom == "" {
				return fmt.Errorf("--from source-selection is required")
			}

			targetSel, err := ResolveSelection(args[0], "")
			if err != nil {
				return err
			}

			sourceSel, err := ResolveSelection(updateFrom, "")
			if err != nil {
				return err
			}

			wp, ok := targetSel.Provider.(provider.WritableProvider)
			if !ok {
				return fmt.Errorf("target provider is not writable")
			}

			// Agnostic Metadata Sync orchestration
			orch := sync.NewOrchestrator(nil, dryRun, verbose) // Placeholder
			matches := orch.Join(sourceSel.Tracks, matchFields)
			
			// We should get fields from a flag, defaulting to beatgrids
			fields := []string{"beatgrids"} 

			return wp.UpdateMetadata(getExecContext(), matches, fields)
		},
	}
	cmd.Flags().StringVarP(&updateFrom, "from", "F", "", "Source selection to read metadata from")
	cmd.Flags().StringSliceVar(&matchFields, "match", []string{"artist", "title"}, "Fields to use for matching tracks")
	return cmd
}
