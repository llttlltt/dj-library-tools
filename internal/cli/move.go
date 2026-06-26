package cli

import (
	"fmt"
	"strings"

	"github.com/llttlltt/dj-library-tools/internal/engine"
	"github.com/llttlltt/dj-library-tools/internal/provider"
	"github.com/spf13/cobra"
)

var (
	moveTo   string
	moveFrom string
	moveName string
)

var moveCmd = &cobra.Command{
	Use:     "move [resource] [query] --to [destination] [--from origin]",
	Aliases: []string{"mv"},
	Short:   "Move items between locations",
	Long: `Move items between locations.
For tracks, both --from and --to are required.
For playlists and folders, only --to (the parent folder) is required.

Use the --name flag to rename a resource.

Example:
  djlt mv rb/tracks "bpm:>130" --from "name:Inbox" --to "name:'High Energy'"
  djlt mv rb/playlists name:Inbox --name "Processed"`,
	Args: cobra.MinimumNArgs(1),
	RunE: runMoveCmd,
}

func runMoveCmd(cmd *cobra.Command, args []string) error {
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
		return runMoveTracks(wp, src)
	}

	return runMoveNodes(wp, src)
}

func runMoveTracks(wp provider.WritableProvider, src *Selection) error {
	if len(src.Tracks) == 0 {
		fmt.Println("No tracks matched the source query.")
		return nil
	}

	// 2. Resolve origin playlists
	org, err := ResolveSelection(moveFrom, "")
	if err != nil || len(org.Nodes) == 0 {
		return fmt.Errorf("could not find origin playlist(s) matching %q", moveFrom)
	}

	// 3. Resolve target playlists
	tgt, err := ResolveSelection(moveTo, "")
	if err != nil || len(tgt.Nodes) == 0 {
		return fmt.Errorf("could not find target playlist(s) matching %q", moveTo)
	}

	if dryRun {
		fmt.Printf("[Dry Run] Would move %d tracks from %d origins to %d targets\n", len(src.Tracks), len(org.Nodes), len(tgt.Nodes))
		return nil
	}

	// 4. Perform Move
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

	// Save Rekordbox
	if rb, ok := wp.(*provider.RekordboxProvider); ok {
		_, path, _ := loadXMLFunc()
		return rb.Engine.Library.(engine.WritableLibrary).Save(path)
	}

	return nil
}

func runMoveNodes(wp provider.WritableProvider, src *Selection) error {
	if len(src.Nodes) == 0 {
		fmt.Println("No resources found matching query.")
		return nil
	}

	// Resolve target parent
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

	// Save Rekordbox
	if rb, ok := wp.(*provider.RekordboxProvider); ok {
		_, path, _ := loadXMLFunc()
		return rb.Engine.Library.(engine.WritableLibrary).Save(path)
	}

	return nil
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

	// Save Rekordbox
	if rb, ok := wp.(*provider.RekordboxProvider); ok {
		_, path, _ := loadXMLFunc()
		return rb.Engine.Library.(engine.WritableLibrary).Save(path)
	}

	return nil
}

func init() {
	moveCmd.Flags().StringVar(&moveTo, "to", "", "Destination playlist or folder")
	moveCmd.Flags().StringVar(&moveFrom, "from", "", "Origin playlist (required for tracks)")
	moveCmd.Flags().StringVar(&moveName, "name", "", "New name for the resource (renames)")
	RootCmd.AddCommand(moveCmd)
}
