package cli

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/llttlltt/dj-library-tools/internal/engine"
	syncpkg "github.com/llttlltt/dj-library-tools/internal/sync"
	"github.com/llttlltt/dj-library-tools/internal/utils"
	"github.com/spf13/cobra"
)

var (
	createIn     string
	createAt     int
	createFrom   string
	createDryRun bool
)

var createCmd = &cobra.Command{
	Use:   "create [resource] [name]",
	Short: "Create a new playlist or folder",
	Long: `Create a new Rekordbox playlist or folder.
You can optionally populate it immediately using items from a source.

Example:
  djlt create rb/playlists "New Arrivals" --from "rb/tracks added:>2024-01-01"`,
	Args: cobra.ExactArgs(2),
	RunE: runCreateCmd,
}

func runCreateCmd(cmd *cobra.Command, args []string) error {
	rbXML, path, err := loadXML()
	if err != nil {
		return err
	}

	eng := engine.NewEngine(rbXML)
	syncEng := syncpkg.NewEngine(nil, rbXML)

	loc := utils.ParseLocation(args[0], "")
	name := args[1]

	var trackIDs []string
	if createFrom != "" {
		src := utils.ParseLocation(createFrom, "")
		if src.Provider != "rb" || src.Resource != "tracks" {
			return fmt.Errorf("--from currently only supports rb/tracks for initial population")
		}

		tracks, err := eng.Ls(src.Query)
		if err != nil {
			return fmt.Errorf("failed to resolve source tracks: %w", err)
		}
		for _, t := range tracks {
			trackIDs = append(trackIDs, strconv.Itoa(t.TrackID))
		}
	}

	if createDryRun {
		fmt.Printf("[Dry Run] Would create %s %q in folder %q with %d tracks\n", loc.Resource, name, createIn, len(trackIDs))
		return nil
	}

	if loc.Resource == "playlists" {
		result := syncEng.UpsertPlaylist(createIn, name, trackIDs, createAt)
		if result.Updated {
			fmt.Printf("Updated existing playlist %q (%d tracks)\n", result.PlaylistName, result.TracksInjected)
		} else {
			fmt.Printf("Created playlist %q (%d tracks)\n", result.PlaylistName, result.TracksInjected)
		}
	} else if loc.Resource == "folders" {
		if !syncEng.CreateFolder(createIn, name, createAt) {
			return fmt.Errorf("failed to create folder %q", name)
		}
		fmt.Printf("Created folder %q\n", name)
	} else {
		return fmt.Errorf("create only supports rb/playlists and rb/folders")
	}

	return syncEng.SaveXML(path)
}

func init() {
	createCmd.Flags().StringVar(&createIn, "in", "", "Parent folder for the new resource")
	createCmd.Flags().IntVar(&createAt, "at", -1, "Insert at this 0-indexed position (-1 for end)")
	createCmd.Flags().StringVar(&createFrom, "from", "", "Initial items to populate the resource with")
	createCmd.Flags().BoolVar(&createDryRun, "dry-run", false, "Preview changes without writing")
	RootCmd.AddCommand(createCmd)
}
