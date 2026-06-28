package cli

import (
	"fmt"

	"github.com/llttlltt/dj-library-tools/internal/models"
	"github.com/llttlltt/dj-library-tools/internal/resolver"
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

			prov := sel.Provider

			ctx := getExecContext()

			if len(deleteFrom) == 0 {
				return runDeleteResources(prov, ctx, sel, recursive)
			}

			return runRemoveMembership(prov, ctx, sel, deleteFrom)
		},
	}
	cmd.Flags().StringSliceVar(&deleteFrom, "from", []string{}, "Origin resource(s) to remove from")
	cmd.Flags().BoolVarP(&recursive, "recursive", "r", false, "Delete folder and all its contents")
	return cmd
}

func runDeleteResources(prov provider.Provider, ctx provider.ExecutionContext, sel *resolver.Selection, recursive bool) error {
	if len(sel.Items) == 0 {
		fmt.Println("No resources matched the query.")
		return nil
	}

	for _, item := range sel.Items {
		if node, ok := item.(models.ResourceGroup); ok {
			if recursive && node.Kind == models.GroupKindFolder {
				// Recursive delete: find children and delete them first
				// This is a simple CLI-side orchestration
				children, _ := prov.Groups().List(ctx, fmt.Sprintf("parent:%q", node.Name))
				
				for _, c := range children {
					if !dryRun { prov.Groups().Delete(ctx, c) }
				}
			}

			if dryRun {
				fmt.Printf("[Dry Run] Would delete %s %q\n", node.GetKind(), node.Name)
				continue
			}
			if err := sel.Provider.Groups().Delete(ctx, node); err != nil {
				return HandleError(err)
			}
			fmt.Printf("Deleted %s %q\n", node.GetKind(), node.Name)
		}
	}

	return sel.Provider.System().Save(ctx, "")
}

func runRemoveMembership(prov provider.Provider, ctx provider.ExecutionContext, sel *resolver.Selection, from []string) error {
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
			removed, _ := prov.Tracks().Groups().Remove(ctx, sel.Tracks, target)
			fmt.Printf("Removed %d tracks from %q\n", removed, target.Name)
		}
	}

	return sel.Provider.System().Save(ctx, "")
}
