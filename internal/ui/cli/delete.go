package cli

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

func newDeleteCmd() *cobra.Command {
	var deleteFrom []string
	var recursive bool

	cmd := &cobra.Command{
		Use:   "rm [resource] [query]",
		Short: "Permanently delete resources or remove membership",
		Long: `Permanently delete resources or remove track membership from playlists.

Use --from to specify which playlist to remove tracks from.
Without --from, the command deletes the resource itself.

Example:
  djlt rm rb/tracks "artist:Four" --from "rb/playlists name:Inbox"
  djlt rm rb/playlists name:Inbox
  djlt rm rb/folders name:OldSets --recursive`,
		Args: cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			queryOverride := ""
			if len(args) > 1 {
				queryOverride = strings.Join(args[1:], " ")
			}

			orch := getOrchestrator()
			runOpts := getRunOptions()

			count, err := orch.Delete(cmd.Context(), args[0], queryOverride, runOpts, deleteFrom, recursive)
			if err != nil {
				return HandleError(err)
			}

			if apply {
				fmt.Printf("✔ Removed %d item(s).\n", count)
			} else {
				fmt.Printf("Scope: %d item(s) would be removed.\n", count)
				fmt.Println("Run with --apply to commit.")
			}

			return nil
		},
	}
	cmd.Flags().StringSliceVar(&deleteFrom, "from", []string{}, "Origin resource(s) to remove from")
	cmd.Flags().BoolVarP(&recursive, "recursive", "r", false, "Delete folder and all its contents")
	return cmd
}
