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
	xmlPath    string
	dryRun     bool
	verbose    bool
	jsonOutput bool
	// loadXMLFunc allows overriding the XML loading logic for testing.
	loadXMLFunc = loadXML
)

// loadXML resolves and loads the Rekordbox XML library, preferring --xml flag over config.
func loadXML() (*rekordbox.RekordboxLibraryXML, string, error) {
	cfg, _ := config.LoadAppConfig()
	path := utils.ExpandPath(xmlPath)
	if path == "" {
		path = utils.ExpandPath(cfg.RekordboxXMLPath)
	}
	if path == "" {
		return nil, "", fmt.Errorf("rekordbox XML path required; use --xml or run 'djlt config rekordbox --xml PATH'")
	}
	rbXML, err := rekordbox.ReadRekordboxLibrary(path)
	if err != nil {
		return nil, "", fmt.Errorf("failed to read rekordbox library: %w", err)
	}
	return rbXML, path, nil
}

// NewRootCmd builds and returns a fully wired root command. Each call
// produces an independent command tree — verb-specific flag vars live in
// closures so there is no shared mutable state between instances. Tests
// call NewRootCmd() directly to get isolation without a reset helper.
func NewRootCmd() *cobra.Command {
	root := &cobra.Command{
		Use:   "djlt",
		Short: "DJ Library Tools CLI",
		Long:  `A comprehensive CLI tool for managing DJ libraries and Rekordbox XMLs.`,
	}
	root.PersistentFlags().StringVarP(&xmlPath, "xml", "x", "", "Path to the Rekordbox XML library")
	root.PersistentFlags().BoolVar(&dryRun, "dry-run", false, "Preview changes without writing")
	root.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Enable verbose logging")
	root.PersistentFlags().BoolVar(&jsonOutput, "json", false, "Output results in JSON format")

	root.AddCommand(
		newListCmd(),
		newSyncCmd(),
		newMakeCmd(),
		newMoveCmd(),
		newDeleteCmd(),
		newAuthCmd(),
		newConfigCmd(),
		newFixCmd(),
		newUpdateCmd(),
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
