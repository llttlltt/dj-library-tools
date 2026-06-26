package cli

import (
	"fmt"
	"os"

	"github.com/llttlltt/dj-library-tools/internal/config"
	"github.com/llttlltt/dj-library-tools/internal/utils"
	"github.com/llttlltt/dj-library-tools/pkg/rekordbox"
	"github.com/spf13/cobra"
)

var (
	xmlPath string
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

var RootCmd = &cobra.Command{
	Use:   "djlt",
	Short: "DJ Library Tools CLI",
	Long:  `A comprehensive CLI tool for managing DJ libraries and Rekordbox XMLs.`,
}

func init() {
	RootCmd.PersistentFlags().StringVarP(&xmlPath, "xml", "x", "", "Path to the Rekordbox XML library")
}

func Execute() error {
	return RootCmd.Execute()
}

func main() {
	if err := Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
