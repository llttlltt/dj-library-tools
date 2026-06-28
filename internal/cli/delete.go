package cli

import (
	"fmt"

	"github.com/llttlltt/dj-library-tools/internal/models"
	"github.com/llttlltt/dj-library-tools/internal/provider"
	"github.com/spf13/cobra"
)

func newDeleteCmd() *cobra.Command {
	var deleteFrom []string

	cmd := &cobra.Command{
		Use:   "rm [resource] [query]",
		Short: "Permanently delete resources or remove membership",
		Long: `Permanently delete resources or remove track membership from playlists.

Use --from to specify which playlist to remove tracks from.
Without --from, the command deletes the resource itself.

Example:
  djlt rm rb/tracks "artist:Four" --from "rb/playlists name:Inbox"
  djlt rm rb/playlists name:Inbox`,
		Args: cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			sel, err := ResolveSelection(args[0], "")
			if err != nil {
				return err
			}

			wp, ok := sel.Provider.(provider.WritableProvider)
			if !ok {
				return fmt.Errorf("provider %q does not support deleting resources", sel.Location.Provider)
			}

			ctx := getExecContext()

			if len(deleteFrom) == 0 {
				return runDeleteResources(wp, ctx, sel)
			}

			return runRemoveMembership(wp, ctx, sel, deleteFrom)
		},
	}
	cmd.Flags().StringSliceVar(&deleteFrom, "from", []string{}, "Origin resource(s) to remove from")
	return cmd
}

func runDeleteResources(wp provider.WritableProvider, ctx provider.ExecutionContext, sel *Selection) error {
	if len(sel.Items) == 0 {
		fmt.Println("No resources matched the query.")
		return nil
	}

	for _, item := range sel.Items {
		if node, ok := item.(models.ResourceGroup); ok {
			if dryRun {
				fmt.Printf("[Dry Run] Would delete %s %q\n", node.GetKind(), node.Name)
				continue
			}
			if err := wp.DeleteGroup(ctx, node); err != nil {
				return err
			}
			fmt.Printf("Deleted %s %q\n", node.GetKind(), node.Name)
		}
	}

	return wp.Save(ctx, "")
}

func runRemoveMembership(wp provider.WritableProvider, ctx provider.ExecutionContext, sel *Selection, from []string) error {
	if len(sel.Tracks) == 0 {
		fmt.Println("No tracks matched the query.")
		return nil
	}

	for _, fromStr := range from {
		org, err := ResolveSelection(fromStr, "")
		if err != nil {
			return err
		}

		for _, target := range org.Groups {
			if dryRun {
				fmt.Printf("[Dry Run] Would remove %d tracks from playlist %q\n", len(sel.Tracks), target.Name)
				continue
			}
			removed, _ := wp.RemoveTracks(ctx, target, sel.Tracks)
			fmt.Printf("Removed %d tracks from %q\n", removed, target.Name)
		}
	}

	return wp.Save(ctx, "")
}
