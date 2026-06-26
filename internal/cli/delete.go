package cli

import (
	"fmt"
	"strings"

	"github.com/llttlltt/dj-library-tools/internal/engine"
	syncpkg "github.com/llttlltt/dj-library-tools/internal/sync"
	"github.com/llttlltt/dj-library-tools/internal/utils"
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

		eng := engine.NewEngine(engine.NewRekordboxLibrary(rbXML))
		syncEng := syncpkg.NewEngine(nil, engine.NewRekordboxLibrary(rbXML))

		query := ""
		if len(args) > 1 {
			query = strings.Join(args[1:], " ")
		}
		loc := utils.ParseLocation(args[0], query)

		if loc.Resource == "tracks" {
			return fmt.Errorf("deleting tracks from collection is not yet supported; use remove to unlink from playlists")
		}

		var targets []engine.NodeResult
		var nodeType int
		if loc.Resource == "playlists" {
			targets, _ = eng.LsPlaylists(loc.Query)
			nodeType = 1
		} else if loc.Resource == "folders" {
			targets, _ = eng.LsFolders(loc.Query)
			nodeType = 0
		} else {
			return fmt.Errorf("delete only supports rb/playlists and rb/folders")
		}

		if len(targets) == 0 {
			fmt.Println("No resources found matching query.")
			return nil
		}

		if dryRun {
			for _, t := range targets {
				fmt.Printf("[Dry Run] Would delete %s %q\n", loc.Resource, t.Node.Name)
			}
			return nil
		}

		for _, t := range targets {
			if verbose {
				fmt.Printf("Deleting %s %q...\n", loc.Resource, t.Node.Name)
			}
			if !syncEng.RemoveNode(t.Node.Name, int32(nodeType)) {
				fmt.Printf("Warning: failed to delete %q\n", t.Node.Name)
				continue
			}
			fmt.Printf("Deleted %s %q\n", loc.Resource, t.Node.Name)
		}

		return syncEng.SaveXML(path)
	},
}

func init() {
	RootCmd.AddCommand(deleteCmd)
}
