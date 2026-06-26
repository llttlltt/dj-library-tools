package cli

import (
	"fmt"
	"strings"

	"github.com/llttlltt/dj-library-tools/internal/provider"
	"github.com/spf13/cobra"
)

func newMoveCmd() *cobra.Command {
	var moveTo, moveFrom, moveName string

	cmd := &cobra.Command{
		Use:     "mv [resource] [query] --to [destination] [--from origin]",
		Short:   "Move items between locations",
		Long: `Move items between locations.
For tracks, both --from and --to are required.
For playlists and folders, only --to (the parent folder) is required.

Use the --name flag to rename a resource.

Example:
  djlt mv rb/tracks "bpm:>130" --from "name:Inbox" --to "name:'High Energy'"
  djlt mv rb/playlists name:Inbox --name "Processed"`,
		Args: cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if moveTo == "" && moveName == "" {
				return fmt.Errorf("either --to destination or --name (for rename) is required")
			}

			queryOverride := ""
			if len(args) > 1 {
				queryOverride = strings.Join(args[1:], " ")
			}
			src, err := ResolveSelection(args[0], queryOverride)
			if err != nil {
				return err
			}

			wp, ok := src.Provider.(provider.WritableProvider)
			if !ok {
				return fmt.Errorf("provider %q does not support moving resources", src.Location.Provider)
			}

			if moveName != "" {
				return runRenameNodes(wp, src, moveName)
			}

			if src.Location.Resource == "tracks" {
				if moveFrom == "" {
					return fmt.Errorf("--from origin is required when moving tracks")
				}
				return runMoveTracks(wp, src, moveFrom, moveTo)
			}

			return runMoveNodes(wp, src, moveTo)
		},
	}
	cmd.Flags().StringVar(&moveTo, "to", "", "Destination playlist or folder")
	cmd.Flags().StringVar(&moveFrom, "from", "", "Origin playlist (required for tracks)")
	cmd.Flags().StringVar(&moveName, "name", "", "New name for the resource (renames)")
	return cmd
}

func runMoveTracks(wp provider.WritableProvider, src *Selection, moveFrom, moveTo string) error {
	if len(src.Tracks) == 0 {
		fmt.Println("No tracks matched the source query.")
		return nil
	}

	org, err := ResolveSelection(moveFrom, "")
	if err != nil || len(org.Nodes) == 0 {
		return fmt.Errorf("could not find origin playlist(s) matching %q", moveFrom)
	}

	tgt, err := ResolveSelection(moveTo, "")
	if err != nil || len(tgt.Nodes) == 0 {
		return fmt.Errorf("could not find target playlist(s) matching %q", moveTo)
	}

	if dryRun {
		fmt.Printf("[Dry Run] Would move %d tracks from %d origins to %d targets\n", len(src.Tracks), len(org.Nodes), len(tgt.Nodes))
		return nil
	}

	for _, origin := range org.Nodes {
		if verbose {
			fmt.Printf("Removing tracks from origin playlist %q...\n", origin.Name)
		}
		wp.RemoveTracks(origin, src.Tracks)
	}
	for _, target := range tgt.Nodes {
		if verbose {
			fmt.Printf("Adding tracks to target playlist %q...\n", target.Name)
		}
		wp.AddTracks(target, src.Tracks)
	}

	savePath := ""
	if src.Location.Provider == "rb" || src.Location.Provider == "rekordbox" {
		_, p, _ := loadXMLFunc()
		savePath = p
	}
	return wp.Save(savePath)
}

func runMoveNodes(wp provider.WritableProvider, src *Selection, moveTo string) error {
	if len(src.Nodes) == 0 {
		fmt.Println("No resources found matching query.")
		return nil
	}

	tgt, err := ResolveSelection(moveTo, "")
	if err != nil || len(tgt.Nodes) == 0 {
		return fmt.Errorf("could not find target folder matching %q", moveTo)
	}
	targetParent := tgt.Nodes[0]

	if dryRun {
		for _, t := range src.Nodes {
			fmt.Printf("[Dry Run] Would move %s %q to folder %q\n", src.Location.Resource, t.Name, targetParent.Name)
		}
		return nil
	}

	for _, t := range src.Nodes {
		if verbose {
			fmt.Printf("Moving %s %q into folder %q...\n", src.Location.Resource, t.Name, targetParent.Name)
		}
		if err := wp.MoveNode(t, targetParent); err != nil {
			fmt.Printf("Warning: failed to move %q: %v\n", t.Name, err)
			continue
		}
		fmt.Printf("Moved %s %q -> %q\n", src.Location.Resource, t.Name, targetParent.Name)
	}

	savePath := ""
	if src.Location.Provider == "rb" || src.Location.Provider == "rekordbox" {
		_, p, _ := loadXMLFunc()
		savePath = p
	}
	return wp.Save(savePath)
}

func runRenameNodes(wp provider.WritableProvider, src *Selection, newName string) error {
	if len(src.Nodes) == 0 {
		return fmt.Errorf("no resources found matching query %q", src.Location.Query)
	}
	if len(src.Nodes) > 1 {
		return fmt.Errorf("rename matched %d resources; refine your query to match exactly one", len(src.Nodes))
	}

	target := src.Nodes[0]
	if verbose {
		fmt.Printf("Renaming %s %q -> %q...\n", src.Location.Resource, target.Name, newName)
	}

	if dryRun {
		fmt.Printf("[Dry Run] Would rename %q to %q\n", target.Name, newName)
		return nil
	}

	if err := wp.RenameNode(target, newName); err != nil {
		return fmt.Errorf("failed to rename %q: %v", target.Name, err)
	}

	fmt.Printf("Renamed %q -> %q\n", target.Name, newName)

	savePath := ""
	if src.Location.Provider == "rb" || src.Location.Provider == "rekordbox" {
		_, p, _ := loadXMLFunc()
		savePath = p
	}
	return wp.Save(savePath)
}
