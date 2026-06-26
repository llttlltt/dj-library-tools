package cli

import (
	"fmt"
	"strings"

	"github.com/llttlltt/dj-library-tools/internal/engine"
	"github.com/llttlltt/dj-library-tools/internal/provider"
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

		queryOverride := ""
		if len(args) > 1 {
			queryOverride = strings.Join(args[1:], " ")
		}
		sel, err := ResolveSelection(args[0], queryOverride)
		if err != nil {
			return err
		}

		wp, ok := sel.Provider.(provider.WritableProvider)
		if !ok {
			return fmt.Errorf("provider %q does not support renaming resources", sel.Location.Provider)
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

		if verbose {
			fmt.Printf("Renaming %s %q -> %q...\n", sel.Location.Resource, target.Name, renameTo)
		}

		if dryRun {
			fmt.Printf("[Dry Run] Would rename %q to %q\n", target.Name, renameTo)
			return nil
		}

		if err := wp.RenameNode(target, renameTo); err != nil {
			return fmt.Errorf("failed to rename %q: %v", target.Name, err)
		}

		fmt.Printf("Renamed %q -> %q\n", target.Name, renameTo)

		// Save Rekordbox
		if rb, ok := wp.(*provider.RekordboxProvider); ok {
			_, path, _ := loadXMLFunc()
			return rb.Engine.Library.(engine.WritableLibrary).Save(path)
		}

		return nil
	},
}

func init() {
	renameCmd.Flags().StringVar(&renameTo, "to", "", "The new name for the resource")
	RootCmd.AddCommand(renameCmd)
}
