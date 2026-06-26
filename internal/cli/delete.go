package cli

import (
	"fmt"
	"strings"

	"github.com/llttlltt/dj-library-tools/internal/engine"
	"github.com/llttlltt/dj-library-tools/internal/provider"
	"github.com/spf13/cobra"
)

var deleteCmd = &cobra.Command{
	Use:     "delete [resource] [query]",
	Aliases: []string{"del", "rm"},
	Short:   "Delete a resource from the library (destructive)",
	Long: `Permanently delete playlists or folders from the Rekordbox XML.
Warning: This is destructive to the resource, but does not delete tracks from the collection.

Example:
  djlt delete rb/playlists "name:'Old Mixes'"`,
	Args: cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
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
			return fmt.Errorf("provider %q does not support deleting resources", sel.Location.Provider)
		}

		if sel.Location.Resource == "tracks" {
			return fmt.Errorf("deleting tracks from collection is not yet supported; use remove to unlink from playlists")
		}

		if len(sel.Nodes) == 0 {
			fmt.Println("No resources found matching query.")
			return nil
		}

		if dryRun {
			for _, t := range sel.Nodes {
				fmt.Printf("[Dry Run] Would delete %s %q\n", sel.Location.Resource, t.Name)
			}
			return nil
		}

		for _, t := range sel.Nodes {
			if verbose {
				fmt.Printf("Deleting %s %q...\n", sel.Location.Resource, t.Name)
			}
			if err := wp.DeleteNode(t); err != nil {
				fmt.Printf("Warning: failed to delete %q: %v\n", t.Name, err)
				continue
			}
			fmt.Printf("Deleted %s %q\n", sel.Location.Resource, t.Name)
		}

		// For Rekordbox we still need to save.
		if rb, ok := wp.(*provider.RekordboxProvider); ok {
			_, path, _ := loadXMLFunc()
			return rb.Engine.Library.(engine.WritableLibrary).Save(path)
		}

		return nil
	},
}

func init() {
	RootCmd.AddCommand(deleteCmd)
}
