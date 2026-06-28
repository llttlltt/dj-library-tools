package cli

import (
	"fmt"

	"github.com/llttlltt/dj-library-tools/internal/models"
	"github.com/llttlltt/dj-library-tools/internal/provider"
	"github.com/spf13/cobra"
)

func newMakeCmd() *cobra.Command {
	var createIn, createFrom string
	var createAt int

	cmd := &cobra.Command{
		Use:     "mk [resource] [name]",
		Short:   "Create a new playlist or folder",
		Long: `Create a new Rekordbox playlist or folder.
You can optionally populate it immediately using items from a source.

Example:
  djlt mk rb/playlists "New Arrivals" --from "rb/tracks added:>2024-01-01"`,
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runCreateCmd(args, createIn, createFrom, createAt)
		},
	}
	cmd.Flags().StringVar(&createIn, "in", "", "Parent folder for the new resource")
	cmd.Flags().IntVar(&createAt, "at", -1, "Insert at this 0-indexed position (-1 for end)")
	cmd.Flags().StringVar(&createFrom, "from", "", "Initial items to populate the resource with")
	return cmd
}

func runCreateCmd(args []string, createIn, createFrom string, createAt int) error {
	sel, err := ResolveSelection(args[0], "")
	if err != nil {
		return err
	}
	name := args[1]

	wp, ok := sel.Provider.(provider.WritableProvider)
	if !ok {
		return fmt.Errorf("provider %q does not support creating resources", sel.Location.Provider)
	}

	var tracks []models.Track
	if createFrom != "" {
		src, err := ResolveSelection(createFrom, "")
		if err != nil {
			return err
		}
		tracks = src.Tracks
	}

	groupType := models.GroupTypePlaylist
	if sel.Location.Resource == "folders" {
		groupType = models.GroupTypeFolder
	}

	// Agnostic Pre-flight Validation
	if err := wp.ValidateCreateGroup(models.ResourceGroup{Name: createIn}, groupType); err != nil {
		return err
	}
	if len(tracks) > 0 {
		if err := wp.ValidateAddTracks(models.ResourceGroup{Name: name, Type: groupType}); err != nil {
			return err
		}
	}

	ctx := getExecContext()

	if dryRun {
		fmt.Printf("[Dry Run] Would create %s %q in folder %q at position %d with %d tracks\n", sel.Location.Resource, name, createIn, createAt, len(tracks))
		return nil
	}

	newNode, err := wp.CreateGroup(ctx, models.ResourceGroup{Name: createIn}, name, int(groupType), createAt)
	if err != nil {
		return err
	}

	if len(tracks) > 0 {
		added, _ := wp.AddTracks(ctx, newNode, tracks)
		fmt.Printf("Created %s %q with %d tracks\n", sel.Location.Resource, name, added)
	} else {
		fmt.Printf("Created %s %q\n", sel.Location.Resource, name)
	}

	return wp.Save(ctx, "")
}
