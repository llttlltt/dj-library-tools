package cli

import (
	"fmt"
	"strconv"

	"github.com/llttlltt/dj-library-tools/internal/core/models"
	"github.com/spf13/cobra"
)

func newMakeCmd() *cobra.Command {
	var createIn, createPopulate string
	var createAt string
	var parents bool

	cmd := &cobra.Command{
		Use:   "mk [resource] [name]",
		Short: "Create a new playlist or folder",
		Long: `Create a new Rekordbox playlist or folder.
You can optionally populate it immediately using items from a source.

The --at flag controls insertion position using 1-based indexing or named sentinels:
  --at start   Insert at the first position
  --at end     Append to the end (default when flag is omitted)
  --at 2       Insert at the second position (1-based)
Omitting --at is equivalent to --at end.

Examples:
  djlt mk rb/playlists "New Arrivals" --populate "rb/tracks added:>2024-01-01"
  djlt mk rb/playlists "2024/Jan/Inbox" --parents
  djlt mk rb/playlists "Inbox" --in "Sorting" --at start
  djlt mk rb/playlists "Archive" --in "Sorting" --at end
  djlt mk rb/playlists "Featured" --in "Sorting" --at 2`,

		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			orch := getOrchestrator()
			runOpts := getRunOptions()

			groupKind := models.GroupKindPlaylist
			if args[0] == "rb/folders" { // Simplification for now
				groupKind = models.GroupKindFolder
			}

			var position int
			switch createAt {
			case "", "end":
				position = -1
			case "start":
				position = 0
			default:
				n, err := strconv.Atoi(createAt)
				if err != nil {
					return fmt.Errorf("--at: unrecognised value %q; use start, end, or a positive integer", createAt)
				}
				if n == 0 {
					return fmt.Errorf("--at: 0 is not valid; use --at start, --at end, or a positive integer")
				}
				if n < 0 {
					return fmt.Errorf("--at: negative values are not valid; use --at start, --at end, or a positive integer")
				}
				position = n - 1
			}

			_, err := orch.Make(cmd.Context(), createIn, args[1], runOpts, groupKind, position, createPopulate)
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
	cmd.Flags().StringVar(&createAt, "at", "", `Insert position: a positive integer (1-based), "start", or "end" (default: end)`)
	cmd.Flags().StringVar(&createPopulate, "populate", "", "Source selection to populate the new resource with")
	cmd.Flags().BoolVarP(&parents, "parents", "p", false, "Create parent folders if they don't exist")
	return cmd
}
