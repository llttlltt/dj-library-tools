package cli

import (
	"fmt"
	"strings"

	"github.com/llttlltt/dj-library-tools/internal/models"
	"github.com/llttlltt/dj-library-tools/internal/provider"
	"github.com/spf13/cobra"
)

func newDeleteCmd() *cobra.Command {
	var removeOrigins []string

	cmd := &cobra.Command{
		Use:     "rm [resource] [query]",
		Short:   "Remove a resource or membership from the library",
		Long: `Permanently delete resources or remove track membership from playlists.

Use --from to specify which playlist to remove tracks from.
Without --from, the command deletes the resource itself.

Example:
  djlt rm rb/tracks artist:Four --from "rb/playlists name:Inbox"
  djlt rm rb/playlists name:Inbox`,
		Args: cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			queryOverride := ""
			if len(args) > 1 {
				queryOverride = strings.Join(args[1:], " ")
			}
			sel, err := ResolveSelection(args[0], queryOverride)
			if err != nil {
				return err
			}

			wp, ok := sel.Provider.(provider.WritableProvider)
			if !ok {
				return fmt.Errorf("provider %q does not support removal", sel.Location.Provider)
			}

			// Membership removal case
			if cmd.Flags().Changed("from") {
				if sel.Location.Resource != "tracks" {
					return fmt.Errorf("can only remove tracks from playlists")
				}
				return runRemoveMembership(wp, sel, removeOrigins)
			}

			// Resource deletion case
			if sel.Location.Resource == "tracks" {
				return fmt.Errorf("deleting tracks from collection is not yet supported; use --from to unlink from playlists")
			}

			if len(sel.Nodes) == 0 {
				fmt.Println("No resources found matching query.")
				return nil
			}

			if dryRun {
				kind := strings.TrimSuffix(sel.Location.Resource, "s")
				for _, t := range sel.Nodes {
					fmt.Printf("[Dry Run] Would delete %s %q\n", kind, t.Name)
				}
				return nil
			}

			for _, t := range sel.Nodes {
				if verbose {
					fmt.Printf("Deleting %s %q...\n", sel.Location.Resource, t.Name)
				}
				if err := wp.DeleteNode(t); err != nil {
					fmt.Printf("Warning: failed to delete %q: %v\n", t.Name, err)
					continue
				}
				fmt.Printf("Deleted %s %q\n", sel.Location.Resource, t.Name)
			}

			_, path, _ := loadXMLFunc()
			return wp.Save(path)
		},
	}
	cmd.Flags().StringSliceVar(&removeOrigins, "from", []string{}, "Origin resource(s) to remove from (repeatable)")
	return cmd
}

func runRemoveMembership(wp provider.WritableProvider, src *Selection, removeOrigins []string) error {
	var targetNodes []models.Node

	for _, originStr := range removeOrigins {
		org, err := ResolveSelection(originStr, "")
		if err != nil {
			return err
		}
		if org.Location.Resource != "playlists" {
			return fmt.Errorf("can only remove from playlists, got %q", org.Location.Resource)
		}
		targetNodes = append(targetNodes, org.Nodes...)
	}

	if dryRun {
		for _, n := range targetNodes {
			fmt.Printf("[Dry Run] Would remove %d tracks from playlist %q\n", len(src.Tracks), n.Name)
		}
		return nil
	}

	for _, n := range targetNodes {
		removed, err := wp.RemoveTracks(n, src.Tracks)
		if err != nil {
			return err
		}
		fmt.Printf("Removed %d tracks from %q\n", removed, n.Name)
	}

	_, path, _ := loadXMLFunc()
	return wp.Save(path)
}
