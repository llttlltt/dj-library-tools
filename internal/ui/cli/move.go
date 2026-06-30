package cli

import (
	"fmt"
	"strings"

	"github.com/llttlltt/dj-library-tools/internal/providers"
	"github.com/llttlltt/dj-library-tools/internal/services/resolver"
	"github.com/spf13/cobra"
)

func newMoveCmd() *cobra.Command {
	var moveTo, moveFrom, moveName string

	cmd := &cobra.Command{
		Use:     "mv [resource] [query] --to [destination] [--from origin]",
		Short:   "Move items between locations",
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
			src, prov, err := ResolveSelection(args[0], queryOverride)
			if err != nil {
				return HandleError(err)
			}

			if moveName != "" {
				return runRenameGroups(prov, src, moveName)
			}

			if src.Location.Resource == "tracks" {
				if moveFrom == "" {
					return fmt.Errorf("--from origin is required when moving tracks")
				}
				return runMoveTracks(prov, src, moveFrom, moveTo)
			}

			return runMoveGroups(prov, src, moveTo)
		},
	}
	cmd.Flags().StringVar(&moveTo, "to", "", "Destination playlist or folder")
	cmd.Flags().StringVar(&moveFrom, "from", "", "Origin playlist (required for tracks)")
	cmd.Flags().StringVar(&moveName, "name", "", "New name for the resource (renames)")
	return cmd
}

func runMoveTracks(prov provider.Provider, src *resolver.Selection, moveFrom, moveTo string) error {
	if len(src.Tracks) == 0 {
		fmt.Println("No tracks matched the source query.")
		return nil
	}

	org, _, err := ResolveSelection(moveFrom, "")
	if err != nil || len(org.Groups) == 0 {
		return fmt.Errorf("could not find origin playlist(s) matching %q", moveFrom)
	}

	tgt, _, err := ResolveSelection(moveTo, "")
	if err != nil || len(tgt.Groups) == 0 {
		return fmt.Errorf("could not find target playlist(s) matching %q", moveTo)
	}

	ctx := getExecContext()
	totalMoved := 0
	for _, origin := range org.Groups {
		for _, target := range tgt.Groups {
			moved, err := prov.Tracks().Groups().Move(ctx, src.Tracks, origin, target)
			if err != nil {
				return HandleError(err)
			}
			totalMoved += moved
		}
	}

	if ctx.Apply {
		fmt.Printf("Successfully moved %d tracks.\n", totalMoved)
		return prov.System().Save(ctx, "")
	} else {
		fmt.Println("Run with --apply to persist changes.")
		return nil
	}
}

func runMoveGroups(prov provider.Provider, src *resolver.Selection, moveTo string) error {
	if len(src.Groups) == 0 {
		fmt.Println("No resources found matching query.")
		return nil
	}

	tgt, _, err := ResolveSelection(moveTo, "")
	if err != nil || len(tgt.Groups) == 0 {
		return fmt.Errorf("could not find target folder matching %q", moveTo)
	}
	targetParent := tgt.Groups[0]

	ctx := getExecContext()

	for _, t := range src.Groups {
		if verbose && ctx.Apply {
			fmt.Printf("Moving %s %q into folder %q...\n", src.Location.Resource, t.Name, targetParent.Name)
		}
		if err := prov.Groups().Update(ctx, t, "", &targetParent); err != nil {
			fmt.Printf("Warning: failed to move %q: %v\n", t.Name, err)
			continue
		}
		if ctx.Apply {
			fmt.Printf("Moved %s %q -> %q\n", src.Location.Resource, t.Name, targetParent.Name)
		}
	}

	if ctx.Apply {
		return prov.System().Save(ctx, "")
	} else {
		fmt.Println("Run with --apply to persist changes.")
		return nil
	}
}

func runRenameGroups(prov provider.Provider, src *resolver.Selection, newName string) error {
	if len(src.Groups) == 0 {
		return fmt.Errorf("no resources found matching query %q", src.Location.Query)
	}
	if len(src.Groups) > 1 {
		return fmt.Errorf("rename matched %d resources; refine your query to match exactly one", len(src.Groups))
	}

	target := src.Groups[0]
	ctx := getExecContext()

	if verbose && ctx.Apply {
		fmt.Printf("Renaming %s %q -> %q...\n", src.Location.Resource, target.Name, newName)
	}

	if err := prov.Groups().Update(ctx, target, newName, nil); err != nil {
		return fmt.Errorf("failed to rename %q: %v", target.Name, err)
	}

	if ctx.Apply {
		fmt.Printf("Renamed %q -> %q\n", target.Name, newName)
		return prov.System().Save(ctx, "")
	} else {
		fmt.Println("Run with --apply to persist changes.")
		return nil
	}
}
