package cli

import (
	"fmt"
	"strings"

	"github.com/llttlltt/dj-library-tools/internal/engine"
	syncpkg "github.com/llttlltt/dj-library-tools/internal/sync"
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

		if sel.Location.Resource == "tracks" {
			return fmt.Errorf("deleting tracks from collection is not yet supported; use remove to unlink from playlists")
		}

		var nodeType int
		if sel.Location.Resource == "playlists" {
			nodeType = 1
		} else if sel.Location.Resource == "folders" {
			nodeType = 0
		} else {
			return fmt.Errorf("delete only supports rb/playlists and rb/folders")
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
			if !syncEng.RemoveNode(t.Name, int32(nodeType)) {
				fmt.Printf("Warning: failed to delete %q\n", t.Name)
				continue
			}
			fmt.Printf("Deleted %s %q\n", sel.Location.Resource, t.Name)
		}

		return syncEng.SaveXML(path)
	},
}

func init() {
	RootCmd.AddCommand(deleteCmd)
}
