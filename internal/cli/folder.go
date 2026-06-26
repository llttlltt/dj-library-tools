package cli

import (
	"fmt"

	"github.com/llttlltt/dj-library-tools/internal/engine"
	syncpkg "github.com/llttlltt/dj-library-tools/internal/sync"
	"github.com/llttlltt/dj-library-tools/internal/utils"
	"github.com/spf13/cobra"
)

var (
	folderNewFlag    string
	folderRenameFlag string
	folderMoveFlag   string
	folderRemoveFlag bool
	folderParentFlag string
	folderDryRun     bool
)

var folderCmd = &cobra.Command{
	Use:   "folder [rb/folders query] [flags]",
	Short: "Manage rekordbox playlist folders",
	RunE:  runFolderCmd,
}

func runFolderCmd(cmd *cobra.Command, args []string) error {
	exclusiveOps := 0
	if folderRenameFlag != "" {
		exclusiveOps++
	}
	if folderMoveFlag != "" {
		exclusiveOps++
	}
	if folderRemoveFlag {
		exclusiveOps++
	}

	if folderNewFlag == "" && exclusiveOps == 0 {
		return cmd.Help()
	}
	if folderNewFlag != "" && exclusiveOps > 0 {
		return fmt.Errorf("--new cannot be combined with --rename, --move, or --remove")
	}
	if exclusiveOps > 1 {
		return fmt.Errorf("only one of --rename, --move, --remove may be specified at a time")
	}
	if folderNewFlag == "" && len(args) == 0 {
		return fmt.Errorf("a folder query (e.g. rb/folders:name:My Sets) is required")
	}

	rbXML, path, err := loadXML()
	if err != nil {
		return err
	}

	eng := engine.NewEngine(rbXML)
	syncEng := syncpkg.NewEngine(nil, rbXML)

	var targets []engine.NodeResult
	if len(args) > 0 {
		loc := utils.ParseLocation(args[0], "")
		if loc.Provider != "rb" || loc.Resource != "folders" {
			return fmt.Errorf("folder query must use rb/folders: syntax, got %q", args[0])
		}
		targets, err = eng.LsFolders(loc.Query)
		if err != nil {
			return fmt.Errorf("failed to resolve folder query: %w", err)
		}
		if len(targets) == 0 {
			return fmt.Errorf("no folders matched query %q", args[0])
		}
	}

	switch {
	case folderNewFlag != "":
		return runFolderNew(syncEng, path)
	case folderRenameFlag != "":
		return runFolderRename(syncEng, targets, path)
	case folderMoveFlag != "":
		return runFolderMove(syncEng, targets, path)
	case folderRemoveFlag:
		return runFolderRemove(syncEng, targets, path)
	}
	return nil
}

func runFolderNew(syncEng *syncpkg.Engine, path string) error {
	if folderDryRun {
		fmt.Printf("[Dry Run] Would create folder %q under parent %q\n", folderNewFlag, folderParentFlag)
		return nil
	}

	// Create the folder by upserting a temporary placeholder playlist inside it,
	// then removing that placeholder. This triggers findOrCreateFolder internally.
	const tmp = "__djlt_tmp__"
	syncEng.UpsertPlaylist(folderNewFlag, tmp, nil, -1)
	syncEng.RemoveNode(tmp, 1)

	if folderParentFlag != "" {
		syncEng.MoveNode(folderNewFlag, 0, folderParentFlag)
	}

	fmt.Printf("Created folder %q\n", folderNewFlag)
	return syncEng.SaveXML(path)
}

func runFolderRename(syncEng *syncpkg.Engine, targets []engine.NodeResult, path string) error {
	if len(targets) > 1 {
		return fmt.Errorf("--rename matched %d folders; refine your query to match exactly one", len(targets))
	}
	oldName := targets[0].Node.Name

	if folderDryRun {
		fmt.Printf("[Dry Run] Would rename folder %q -> %q\n", oldName, folderRenameFlag)
		return nil
	}

	if !syncEng.RenameNode(oldName, folderRenameFlag, 0) {
		return fmt.Errorf("failed to rename folder %q", oldName)
	}
	fmt.Printf("Renamed folder %q -> %q\n", oldName, folderRenameFlag)
	return syncEng.SaveXML(path)
}

func runFolderMove(syncEng *syncpkg.Engine, targets []engine.NodeResult, path string) error {
	if folderDryRun {
		for _, target := range targets {
			fmt.Printf("[Dry Run] Would move folder %q -> %q\n", target.Node.Name, folderMoveFlag)
		}
		return nil
	}

	for _, target := range targets {
		if !syncEng.MoveNode(target.Node.Name, 0, folderMoveFlag) {
			fmt.Printf("Warning: could not move folder %q\n", target.Node.Name)
			continue
		}
		fmt.Printf("Moved folder %q -> %q\n", target.Node.Name, folderMoveFlag)
	}
	return syncEng.SaveXML(path)
}

func runFolderRemove(syncEng *syncpkg.Engine, targets []engine.NodeResult, path string) error {
	if folderDryRun {
		for _, target := range targets {
			fmt.Printf("[Dry Run] Would remove folder %q\n", target.Node.Name)
		}
		return nil
	}

	for _, target := range targets {
		if !syncEng.RemoveNode(target.Node.Name, 0) {
			fmt.Printf("Warning: could not remove folder %q\n", target.Node.Name)
			continue
		}
		fmt.Printf("Removed folder %q\n", target.Node.Name)
	}
	return syncEng.SaveXML(path)
}

func init() {
	folderCmd.Flags().StringVar(&folderNewFlag, "new", "", "Create a new folder with this name")
	folderCmd.Flags().StringVar(&folderRenameFlag, "rename", "", "Rename matched folder to this name")
	folderCmd.Flags().StringVar(&folderMoveFlag, "move", "", "Move matched folder into this parent folder")
	folderCmd.Flags().BoolVar(&folderRemoveFlag, "remove", false, "Remove matched folder")
	folderCmd.Flags().StringVar(&folderParentFlag, "parent", "", "Parent folder for --new (default: root level)")
	folderCmd.Flags().BoolVar(&folderDryRun, "dry-run", false, "Preview changes without writing")

	RootCmd.AddCommand(folderCmd)
}
