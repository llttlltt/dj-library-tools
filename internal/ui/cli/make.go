package cli

import (
	"fmt"

	"github.com/llttlltt/dj-library-tools/internal/core/models"
	"github.com/spf13/cobra"
)

func newMakeCmd() *cobra.Command {
	var createIn, createPopulate string
	var createAt int
	var parents bool

	cmd := &cobra.Command{
		Use:   "mk [resource] [name]",
		Short: "Create a new playlist or folder",
		Long: `Create a new Rekordbox playlist or folder.
You can optionally populate it immediately using items from a source.

Example:
  djlt mk rb/playlists "New Arrivals" --populate "rb/tracks added:>2024-01-01"
  djlt mk rb/playlists "2024/Jan/Inbox" --parents`,
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			orch := getOrchestrator()
			runOpts := getRunOptions()

			groupKind := models.GroupKindPlaylist
			if args[0] == "rb/folders" { // Simplification for now
				groupKind = models.GroupKindFolder
			}

			_, err := orch.Make(cmd.Context(), createIn, args[1], runOpts, groupKind, createAt, createPopulate)
			if err != nil {
				return HandleError(err)
			}

			if apply {
				fmt.Printf("Created %s %q\n", args[0], args[1])
			} else {
				fmt.Println("Run with --apply to persist changes.")
			}

			return nil
		},
	}
	cmd.Flags().StringVar(&createIn, "in", "", "Parent folder for the new resource")
	cmd.Flags().IntVar(&createAt, "at", -1, "Insert at this 0-indexed position (-1 for end)")
	cmd.Flags().StringVar(&createPopulate, "populate", "", "Source selection to populate the new resource with")
	cmd.Flags().BoolVarP(&parents, "parents", "p", false, "Create parent folders if they don't exist")
	return cmd
}
