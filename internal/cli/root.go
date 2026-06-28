package cli

import (
	"github.com/spf13/cobra"
)

var (
	filePath      string
	toFilePath    string
	dryRun        bool
	verbose       bool
	jsonOutput    bool
	filterMissing bool
	filterExists  bool
)

// NewRootCmd builds and returns a fully wired root command.
func NewRootCmd() *cobra.Command {
	root := &cobra.Command{
		Use:   "djlt",
		Short: "DJ Library Tools CLI",
		Long:  `A comprehensive CLI tool for managing DJ libraries across multiple providers.`,
	}
	root.PersistentFlags().StringVarP(&filePath, "file", "f", "", "Path to the primary library file (Rekordbox XML, M3U, etc.)")
	root.PersistentFlags().StringVar(&toFilePath, "to-file", "", "Path to the destination library file for sync/move operations")
	root.PersistentFlags().BoolVar(&dryRun, "dry-run", false, "Preview changes without writing")
	root.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Enable verbose logging")
	root.PersistentFlags().BoolVar(&jsonOutput, "json", false, "Output results in JSON format")
	root.PersistentFlags().BoolVar(&filterMissing, "missing", false, "Filter for tracks where the physical file is missing")
	root.PersistentFlags().BoolVar(&filterExists, "exists", false, "Filter for tracks where the physical file exists")

	root.AddCommand(
		newListCmd(),
		newSyncCmd(),
		newMakeCmd(),
		newMoveCmd(),
		newDeleteCmd(),
		newConfigCmd(),
		newEditCmd(),
	)
	return root
}

// RootCmd is the singleton used by the production binary.
var RootCmd = NewRootCmd()

func Execute() error {
	return RootCmd.Execute()
}
