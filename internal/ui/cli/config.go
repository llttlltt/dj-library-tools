package cli

import (
	"fmt"

	"github.com/llttlltt/dj-library-tools/internal/config"
	"github.com/spf13/cobra"
)

func newConfigCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "config",
		Short: "Manage application configuration",
	}
	cmd.AddCommand(newSourceListCmd())
	return cmd
}

// newSourceListCmd prints all configured Sources — useful for debugging and
// confirming that ~/.config/djlt/sources/ is populated correctly.
func newSourceListCmd() *cobra.Command {
	sourceCmd := &cobra.Command{
		Use:   "source",
		Short: "Manage Sources",
	}
	sourceCmd.AddCommand(&cobra.Command{
		Use:   "list",
		Short: "List all configured Sources",
		RunE: func(cmd *cobra.Command, args []string) error {
			sources, err := config.LoadSources()
			if err != nil {
				return err
			}
			if len(sources) == 0 {
				fmt.Println("No Sources configured. Add one via the GUI.")
				return nil
			}
			for _, s := range sources {
				fmt.Printf("%-38s  %-8s  %s\n", s.ID, s.Provider, s.Name)
			}
			return nil
		},
	})
	return sourceCmd
}
