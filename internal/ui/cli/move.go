package cli

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

func newMoveCmd() *cobra.Command {
	var moveTo, moveFrom, moveName string

	cmd := &cobra.Command{
		Use:   "mv [resource] [query] --to [destination] [--from origin]",
		Short: "Move items between locations",
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

			orch := getOrchestrator()
			runOpts := getRunOptions()

			count, err := orch.Move(cmd.Context(), args[0], queryOverride, runOpts, moveTo, moveFrom, moveName)
			if err != nil {
				return HandleError(err)
			}

			if apply {
				fmt.Printf("✔ Moved/renamed %d item(s).\n", count)
			} else {
				fmt.Printf("Scope: %d item(s) would be moved/renamed.\n", count)
				fmt.Println("Run with --apply to commit.")
			}

			return nil
		},
	}
	cmd.Flags().StringVar(&moveTo, "to", "", "Destination playlist or folder")
	cmd.Flags().StringVar(&moveFrom, "from", "", "Origin playlist (required for tracks)")
	cmd.Flags().StringVar(&moveName, "name", "", "New name for the resource (renames)")
	return cmd
}
