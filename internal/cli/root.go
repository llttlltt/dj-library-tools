package cli

import (
	"fmt"
	"os"

	"github.com/llttlltt/dj-library-tools/internal/config"
	"github.com/llttlltt/dj-library-tools/internal/utils"
	"github.com/llttlltt/dj-library-tools/internal/rekordbox"
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

// loadXML resolves and loads the Rekordbox XML library, preferring --file flag over config.
func loadXML() (*rekordbox.RekordboxLibraryXML, string, error) {
	cfg, _ := config.LoadAppConfig()
	path := utils.ExpandPath(filePath)
	if path == "" {
		path = utils.ExpandPath(cfg.PrimaryFilePath)
	}
	if path == "" {
		return nil, "", fmt.Errorf("library path required; use --file or run 'djlt config rekordbox --file PATH'")
	}
	rbXML, err := rekordbox.ReadRekordboxLibrary(path)
	if err != nil {
		return nil, "", fmt.Errorf("failed to read rekordbox library: %w", err)
	}
	return rbXML, path, nil
}

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
		newAuthCmd(),
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

func main() {
	if err := Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
