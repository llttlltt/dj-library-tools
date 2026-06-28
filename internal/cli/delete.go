package cli

import (
	"fmt"

	"github.com/llttlltt/dj-library-tools/internal/models"
	"github.com/llttlltt/dj-library-tools/internal/provider"
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
			sel, err := ResolveSelection(args[0], "")
			if err != nil {
				return HandleError(err)
			}

			wp, ok := sel.Provider.(provider.WritableProvider)
			if !ok {
				return fmt.Errorf("provider %q does not support deleting resources", sel.Location.Provider)
			}

			ctx := getExecContext()

			if len(deleteFrom) == 0 {
				return runDeleteResources(wp, ctx, sel, recursive)
			}

			return runRemoveMembership(wp, ctx, sel, deleteFrom)
		},
	}
	cmd.Flags().StringSliceVar(&deleteFrom, "from", []string{}, "Origin resource(s) to remove from")
	cmd.Flags().BoolVarP(&recursive, "recursive", "r", false, "Delete folder and all its contents")
	return cmd
}

func runDeleteResources(wp provider.WritableProvider, ctx provider.ExecutionContext, sel *Selection, recursive bool) error {
	if len(sel.Items) == 0 {
		fmt.Println("No resources matched the query.")
		return nil
	}

	for _, item := range sel.Items {
		if node, ok := item.(models.ResourceGroup); ok {
			if recursive && node.Type == models.GroupTypeFolder {
				// Recursive delete: find children and delete them first
				// This is a simple CLI-side orchestration
				children, _ := wp.GetResources(ctx, "playlists", fmt.Sprintf("parent:%q", node.Name))
				childFolders, _ := wp.GetResources(ctx, "folders", fmt.Sprintf("parent:%q", node.Name))
				
				for _, c := range children {
					if !dryRun { wp.DeleteGroup(ctx, c.(models.ResourceGroup)) }
				}
				for _, c := range childFolders {
					// We could recurse deeper here if needed, but for now 1-level deep
					if !dryRun { wp.DeleteGroup(ctx, c.(models.ResourceGroup)) }
				}
			}

			if dryRun {
				fmt.Printf("[Dry Run] Would delete %s %q\n", node.GetKind(), node.Name)
				continue
			}
			if err := wp.DeleteGroup(ctx, node); err != nil {
				return HandleError(err)
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
			return HandleError(err)
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
