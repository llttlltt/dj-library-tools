package cli

import (
	"fmt"
	"strings"

	"github.com/llttlltt/dj-library-tools/internal/engine"
	syncpkg "github.com/llttlltt/dj-library-tools/internal/sync"
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

		syncEng := syncpkg.NewEngine(nil, engine.NewRekordboxLibrary(rbXML))

		queryOverride := ""
		if len(args) > 1 {
			queryOverride = strings.Join(args[1:], " ")
		}
		sel, err := ResolveSelection(args[0], queryOverride)
		if err != nil {
			return err
		}

		if sel.Location.Resource != "playlists" && sel.Location.Resource != "folders" {
			return fmt.Errorf("rename only supports rb/playlists and rb/folders")
		}

		if len(sel.Nodes) == 0 {
			return fmt.Errorf("no resources found matching query %q", sel.Location.Query)
		}
		if len(sel.Nodes) > 1 {
			return fmt.Errorf("rename matched %d resources; refine your query to match exactly one", len(sel.Nodes))
		}

		target := sel.Nodes[0]
		// Since we don't have nodeType in NodeResult yet, we'll assume based on Resource
		nodeType := int32(1)
		if sel.Location.Resource == "folders" {
			nodeType = 0
		}

		if verbose {
			fmt.Printf("Renaming %s %q -> %q...\n", sel.Location.Resource, target.Name, renameTo)
		}

		if dryRun {
			fmt.Printf("[Dry Run] Would rename %q to %q\n", target.Name, renameTo)
			return nil
		}

		if !syncEng.RenameNode(target.Name, renameTo, nodeType) {
			return fmt.Errorf("failed to rename %q", target.Name)
		}

		fmt.Printf("Renamed %q -> %q\n", target.Name, renameTo)
		return syncEng.SaveXML(path)
	},
}

func init() {
	renameCmd.Flags().StringVar(&renameTo, "to", "", "The new name for the resource")
	RootCmd.AddCommand(renameCmd)
}
