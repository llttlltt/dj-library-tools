package cli

import (
	"fmt"
	"github.com/llttlltt/dj-library-tools/internal/provider"
	"github.com/spf13/cobra"
)

func newFixCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "fix [selection] [query]",
		Short: "Repair formatting or health issues in a library",
		Long: `Enforces formatting standards or repairs broken metadata for a selection.
The exact action depends on the provider (e.g. enriching M3U files with tags).`,
		Args: cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			queryOverride := ""
			if len(args) > 1 {
				queryOverride = args[1]
			}
			sel, err := ResolveSelection(args[0], queryOverride)
			if err != nil {
				return err
			}

			wp, ok := sel.Provider.(provider.WritableProvider)
			if !ok {
				return fmt.Errorf("provider %q is read-only", sel.Location.Provider)
			}

			if dryRun {
				fmt.Printf("[Dry Run] Would perform repair on %s/%s\n", sel.Location.Provider, sel.Location.Resource)
				return nil
			}

			if err := wp.Fix(getExecContext(), sel.Location.Resource, sel.Location.Query); err != nil {
				return err
			}

			fmt.Println("Repair completed successfully.")
			return nil
		},
	}
	return cmd
}
