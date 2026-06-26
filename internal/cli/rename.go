package cli

import (
	"fmt"
	"strings"

	"github.com/llttlltt/dj-library-tools/internal/engine"
	syncpkg "github.com/llttlltt/dj-library-tools/internal/sync"
	"github.com/llttlltt/dj-library-tools/internal/utils"
	"github.com/spf13/cobra"
)

var (
	renameTo string
)

var renameCmd = &cobra.Command{
	Use:   "rename [resource] [query] --to [new-name]",
	Short: "Rename a playlist or folder",
	Long: `Rename a Rekordbox playlist or folder.
The target must resolve to a single resource.

Example:
  djlt rename rb/playlists name:Inbox --to "Inbox (Processed)"`,
	Args: cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if renameTo == "" {
			return fmt.Errorf("--to new-name is required")
		}

		rbXML, path, err := loadXMLFunc()
		if err != nil {
			return err
		}

		eng := engine.NewEngine(engine.NewRekordboxLibrary(rbXML))
		syncEng := syncpkg.NewEngine(nil, engine.NewRekordboxLibrary(rbXML))

		query := ""
		if len(args) > 1 {
			query = strings.Join(args[1:], " ")
		}
		loc := utils.ParseLocation(args[0], query)

		var targets []engine.NodeResult
		if loc.Resource == "playlists" {
			targets, _ = eng.LsPlaylists(loc.Query)
		} else if loc.Resource == "folders" {
			targets, _ = eng.LsFolders(loc.Query)
		} else {
			return fmt.Errorf("rename only supports rb/playlists and rb/folders")
		}

		if len(targets) == 0 {
			return fmt.Errorf("no resources found matching query %q", loc.Query)
		}
		if len(targets) > 1 {
			return fmt.Errorf("rename matched %d resources; refine your query to match exactly one", len(targets))
		}

		target := targets[0]
		nodeType := target.Node.Type

		if verbose {
			fmt.Printf("Renaming %s %q -> %q...\n", loc.Resource, target.Node.Name, renameTo)
		}

		if dryRun {
			fmt.Printf("[Dry Run] Would rename %q to %q\n", target.Node.Name, renameTo)
			return nil
		}

		if !syncEng.RenameNode(target.Node.Name, renameTo, nodeType) {
			return fmt.Errorf("failed to rename %q", target.Node.Name)
		}

		fmt.Printf("Renamed %q -> %q\n", target.Node.Name, renameTo)
		return syncEng.SaveXML(path)
	},
}

func init() {
	renameCmd.Flags().StringVar(&renameTo, "to", "", "The new name for the resource")
	RootCmd.AddCommand(renameCmd)
}
