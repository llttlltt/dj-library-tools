package cli

import (
	"fmt"
	"strings"

	"github.com/llttlltt/dj-library-tools/internal/models"
	"github.com/spf13/cobra"
)

func newMakeCmd() *cobra.Command {
	var createIn, createFrom string
	var createAt int
	var parents bool

	cmd := &cobra.Command{
		Use:     "mk [resource] [name]",
		Short:   "Create a new playlist or folder",
		Long: `Create a new Rekordbox playlist or folder.
You can optionally populate it immediately using items from a source.

Example:
  djlt mk rb/playlists "New Arrivals" --from "rb/tracks added:>2024-01-01"
  djlt mk rb/playlists "2024/Jan/Inbox" --parents`,
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runCreateCmd(args, createIn, createFrom, createAt, parents)
		},
	}
	cmd.Flags().StringVar(&createIn, "in", "", "Parent folder for the new resource")
	cmd.Flags().IntVar(&createAt, "at", -1, "Insert at this 0-indexed position (-1 for end)")
	cmd.Flags().StringVar(&createFrom, "from", "", "Initial items to populate the resource with")
	cmd.Flags().BoolVarP(&parents, "parents", "p", false, "Create parent folders if they don't exist")
	return cmd
}

func runCreateCmd(args []string, createIn, createFrom string, createAt int, parents bool) error {
	sel, err := ResolveSelection(args[0], "")
	if err != nil {
		return HandleError(err)
	}
	name := args[1]

	prov := sel.Provider

	// Handle --parents by ensuring the folder path exists
	if parents && createIn != "" {
		parts := strings.Split(createIn, "/")
		currentParent := ""
		for _, part := range parts {
			if part == "" { continue }
			// Check if part exists as folder
			query := fmt.Sprintf("name:%q", part)
			if currentParent != "" {
				query += fmt.Sprintf(" && parent:%q", currentParent)
			}
			
			res, _ := prov.Groups().List(getExecContext(), query)
			if len(res) == 0 {
				if dryRun {
					fmt.Printf("[Dry Run] Would create folder %q in %q\n", part, currentParent)
				} else {
					_, err := prov.Groups().Create(getExecContext(), models.ResourceGroup{Name: currentParent}, part, models.GroupKindFolder, -1)
					if err != nil { return HandleError(err) }
				}
			}
			currentParent = part
		}
	}

	var tracks []models.Track
	if createFrom != "" {
		src, err := ResolveSelection(createFrom, "")
		if err != nil {
			return HandleError(err)
		}
		tracks = src.Tracks
	}

	groupType := models.GroupKindPlaylist
	if sel.Location.Resource == "folders" {
		groupType = models.GroupKindFolder
	}

	// Handle structural validation by checking provider capabilities or specific constraints
	// (Note: In the new architecture, validation logic can be embedded in the Create call 
	// or specific policy checks on System().Containment())

	ctx := getExecContext()

	if dryRun {
		fmt.Printf("[Dry Run] Would create %s %q in folder %q at position %d with %d tracks\n", sel.Location.Resource, name, createIn, createAt, len(tracks))
		return nil
	}

	newNode, err := prov.Groups().Create(ctx, models.ResourceGroup{Name: createIn}, name, groupType, createAt)
	if err != nil {
		return HandleError(err)
	}

	if len(tracks) > 0 {
		added, _ := prov.Tracks().Groups().Add(ctx, tracks, newNode)
		fmt.Printf("Created %s %q with %d tracks\n", sel.Location.Resource, name, added)
	} else {
		fmt.Printf("Created %s %q\n", sel.Location.Resource, name)
	}

	return prov.System().Save(ctx, "")
}
