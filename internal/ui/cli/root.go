package cli

import (
	"context"
	"github.com/spf13/cobra"
)

var (
	filePath   string
	apply      bool
	verbose    bool
)

// NewRootCmd builds and returns a fully wired root command.
func NewRootCmd() *cobra.Command {
	root := &cobra.Command{
		Use:   "djlt",
		Short: "DJ Library Tools CLI",
		Long:  `A comprehensive CLI tool for managing DJ libraries across multiple providers.`,
	}
	root.PersistentFlags().StringVarP(&filePath, "file", "f", "", "Path to the primary library file (Rekordbox XML, M3U, etc.)")
	root.PersistentFlags().BoolVar(&apply, "apply", false, "Actually apply changes to the library (destructive)")
	root.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Enable verbose logging")

	root.AddCommand(
		newListCmd(),
		newSyncCmd(),
		newMakeCmd(),
		newMoveCmd(),
		newDeleteCmd(),
		newConfigCmd(),
		newEditCmd(),
		newFixCmd(),
	)
	return root
}

// RootCmd is the singleton used by the production binary.
var RootCmd = NewRootCmd()

func ExecuteContext(ctx context.Context) error {
	return RootCmd.ExecuteContext(ctx)
}
