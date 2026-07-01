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
	cmd.AddCommand(newConnectionListCmd())
	return cmd
}

// newConnectionListCmd prints all configured Connections — useful for debugging and
// confirming that ~/.config/djlt/connections/ is populated correctly.
func newConnectionListCmd() *cobra.Command {
	connectionCmd := &cobra.Command{
		Use:   "connection",
		Short: "Manage Connection",
	}
	connectionCmd.AddCommand(&cobra.Command{
		Use:   "list",
		Short: "List all configured Connections",
		RunE: func(cmd *cobra.Command, args []string) error {
			connections, err := config.LoadConnections()
			if err != nil {
				return err
			}
			if len(connections) == 0 {
				fmt.Println("No Connections configured. Add one via the GUI.")
				return nil
			}
			for _, c := range connections {
				fmt.Printf("%-38s  %-8s  %s\n", c.ID, c.Provider, c.Name)
			}
			return nil
		},
	})
	return connectionCmd
}
