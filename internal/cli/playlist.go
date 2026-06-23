package cli

import (
	"fmt"
	"os"
	"strconv"

	"github.com/llttlltt/dj-library-tools/internal/config"
	"github.com/llttlltt/dj-library-tools/internal/engine"
	"github.com/llttlltt/dj-library-tools/internal/playlist"
	syncpkg "github.com/llttlltt/dj-library-tools/internal/sync"
	"github.com/llttlltt/dj-library-tools/internal/utils"
	"github.com/llttlltt/dj-library-tools/pkg/rekordbox"
	"github.com/spf13/cobra"
)

var (
	// playlist fix flags
	extsFlag       []string
	m3u8Flag       bool
	removeOriginal bool
	forceOverwrite bool
	outputFileFlag string
	verboseFlag    bool
	dryRunFlag     bool

	// playlist god command flags
	plAddFlag    string
	plRemoveFlag string
	plSyncFlag   string
	plNewFlag    string
	plRenameFlag string
	plMoveFlag   string
	plDeleteFlag bool
	plFolderFlag string
	plAtFlag     int
	plDryRun     bool
)

var playlistCmd = &cobra.Command{
	Use:   "playlist [rb/playlists query] [flags]",
	Short: "Manage rekordbox playlists",
	RunE:  runPlaylistCmd,
}

func runPlaylistCmd(cmd *cobra.Command, args []string) error {
	// Tally mutually exclusive ops. --add without --new counts as one.
	exclusiveOps := 0
	if plAddFlag != "" && plNewFlag == "" {
		exclusiveOps++
	}
	if plRemoveFlag != "" {
		exclusiveOps++
	}
	if plSyncFlag != "" {
		exclusiveOps++
	}
	if plRenameFlag != "" {
		exclusiveOps++
	}
	if plMoveFlag != "" {
		exclusiveOps++
	}
	if plDeleteFlag {
		exclusiveOps++
	}

	if plNewFlag == "" && exclusiveOps == 0 {
		return cmd.Help()
	}
	if plNewFlag != "" && exclusiveOps > 0 {
		return fmt.Errorf("--new cannot be combined with --remove, --sync, --rename, --move, or --delete")
	}
	if exclusiveOps > 1 {
		return fmt.Errorf("only one of --add, --remove, --sync, --rename, --move, --delete may be specified at a time")
	}
	if plNewFlag == "" && len(args) == 0 {
		return fmt.Errorf("a playlist query (e.g. rb/playlists:name:MyPlaylist) is required")
	}

	rbXML, path, err := loadXML()
	if err != nil {
		return err
	}

	eng := engine.NewEngine(rbXML)
	syncEng := syncpkg.NewEngine(nil, rbXML)

	// Resolve playlist query when provided.
	var targets []engine.NodeResult
	if len(args) > 0 {
		loc := utils.ParseLocation(args[0])
		if loc.Provider != "rb" || loc.Resource != "playlists" {
			return fmt.Errorf("playlist query must use rb/playlists: syntax, got %q", args[0])
		}
		targets, err = eng.LsPlaylists(loc.Query)
		if err != nil {
			return fmt.Errorf("failed to resolve playlist query: %w", err)
		}
		if len(targets) == 0 {
			return fmt.Errorf("no playlists matched query %q", args[0])
		}
	}

	switch {
	case plNewFlag != "":
		return runPlaylistNew(syncEng, eng, path)
	case plAddFlag != "":
		return runPlaylistAdd(syncEng, eng, targets, path)
	case plRemoveFlag != "":
		return runPlaylistRemoveTracks(syncEng, eng, targets, path)
	case plSyncFlag != "":
		return runPlaylistSync(syncEng, eng, targets, path)
	case plRenameFlag != "":
		return runPlaylistRename(syncEng, targets, path)
	case plMoveFlag != "":
		return runPlaylistMove(syncEng, targets, path)
	case plDeleteFlag:
		return runPlaylistDelete(syncEng, targets, path)
	}
	return nil
}

func runPlaylistNew(syncEng *syncpkg.Engine, eng *engine.Engine, path string) error {
	var trackIDs []string
	if plAddFlag != "" {
		loc := utils.ParseLocation(plAddFlag)
		if loc.Provider != "rb" || loc.Resource != "tracks" {
			return fmt.Errorf("--add must use rb/tracks: syntax, got %q", plAddFlag)
		}
		tracks, err := eng.Ls(loc.Query)
		if err != nil {
			return fmt.Errorf("failed to resolve track query: %w", err)
		}
		for _, t := range tracks {
			trackIDs = append(trackIDs, strconv.Itoa(t.TrackID))
		}
	}

	if plDryRun {
		fmt.Printf("[Dry Run] Would create playlist %q in folder %q at position %d with %d tracks\n", plNewFlag, plFolderFlag, plAtFlag, len(trackIDs))
		return nil
	}

	result := syncEng.UpsertPlaylist(plFolderFlag, plNewFlag, trackIDs, plAtFlag)
	if result.Updated {
		fmt.Printf("Updated existing playlist %q (%d tracks)\n", result.PlaylistName, result.TracksInjected)
	} else {
		fmt.Printf("Created playlist %q (%d tracks)\n", result.PlaylistName, result.TracksInjected)
	}
	return syncEng.SaveXML(path)
}

func runPlaylistAdd(syncEng *syncpkg.Engine, eng *engine.Engine, targets []engine.NodeResult, path string) error {
	loc := utils.ParseLocation(plAddFlag)
	if loc.Provider != "rb" || loc.Resource != "tracks" {
		return fmt.Errorf("--add must use rb/tracks: syntax, got %q", plAddFlag)
	}
	tracks, err := eng.Ls(loc.Query)
	if err != nil {
		return fmt.Errorf("failed to resolve track query: %w", err)
	}
	var trackIDs []string
	for _, t := range tracks {
		trackIDs = append(trackIDs, strconv.Itoa(t.TrackID))
	}

	if plDryRun {
		for _, target := range targets {
			fmt.Printf("[Dry Run] Would add %d tracks to playlist %q\n", len(trackIDs), target.Node.Name)
		}
		return nil
	}

	for _, target := range targets {
		found, added := syncEng.AddTracksToPlaylist(target.Node.Name, trackIDs)
		if !found {
			fmt.Printf("Warning: playlist %q not found during add\n", target.Node.Name)
			continue
		}
		fmt.Printf("Added %d tracks to %q\n", added, target.Node.Name)
	}
	return syncEng.SaveXML(path)
}

func runPlaylistRename(syncEng *syncpkg.Engine, targets []engine.NodeResult, path string) error {
	if len(targets) > 1 {
		return fmt.Errorf("--rename matched %d playlists; refine your query to match exactly one", len(targets))
	}
	oldName := targets[0].Node.Name

	if plDryRun {
		fmt.Printf("[Dry Run] Would rename playlist %q -> %q\n", oldName, plRenameFlag)
		return nil
	}

	if !syncEng.RenameNode(oldName, plRenameFlag, 1) {
		return fmt.Errorf("failed to rename playlist %q", oldName)
	}
	fmt.Printf("Renamed playlist %q -> %q\n", oldName, plRenameFlag)
	return syncEng.SaveXML(path)
}

func runPlaylistMove(syncEng *syncpkg.Engine, targets []engine.NodeResult, path string) error {
	if plDryRun {
		for _, target := range targets {
			fmt.Printf("[Dry Run] Would move playlist %q -> folder %q\n", target.Node.Name, plMoveFlag)
		}
		return nil
	}

	for _, target := range targets {
		if !syncEng.MoveNode(target.Node.Name, 1, plMoveFlag) {
			fmt.Printf("Warning: could not move playlist %q\n", target.Node.Name)
			continue
		}
		fmt.Printf("Moved playlist %q -> folder %q\n", target.Node.Name, plMoveFlag)
	}
	return syncEng.SaveXML(path)
}

func runPlaylistRemoveTracks(syncEng *syncpkg.Engine, eng *engine.Engine, targets []engine.NodeResult, path string) error {
	loc := utils.ParseLocation(plRemoveFlag)
	if loc.Provider != "rb" || loc.Resource != "tracks" {
		return fmt.Errorf("--remove must use rb/tracks: syntax, got %q", plRemoveFlag)
	}
	tracks, err := eng.Ls(loc.Query)
	if err != nil {
		return fmt.Errorf("failed to resolve track query: %w", err)
	}
	var trackIDs []string
	for _, t := range tracks {
		trackIDs = append(trackIDs, strconv.Itoa(t.TrackID))
	}

	if plDryRun {
		for _, target := range targets {
			fmt.Printf("[Dry Run] Would remove %d tracks from playlist %q\n", len(trackIDs), target.Node.Name)
		}
		return nil
	}

	for _, target := range targets {
		found, removed := syncEng.RemoveTracksFromPlaylist(target.Node.Name, trackIDs)
		if !found {
			fmt.Printf("Warning: playlist %q not found\n", target.Node.Name)
			continue
		}
		fmt.Printf("Removed %d tracks from %q\n", removed, target.Node.Name)
	}
	return syncEng.SaveXML(path)
}

func runPlaylistSync(syncEng *syncpkg.Engine, eng *engine.Engine, targets []engine.NodeResult, path string) error {
	loc := utils.ParseLocation(plSyncFlag)
	if loc.Provider != "rb" || loc.Resource != "tracks" {
		return fmt.Errorf("--sync must use rb/tracks: syntax, got %q", plSyncFlag)
	}
	tracks, err := eng.Ls(loc.Query)
	if err != nil {
		return fmt.Errorf("failed to resolve track query: %w", err)
	}
	newIDs := make(map[string]struct{}, len(tracks))
	var newIDSlice []string
	for _, t := range tracks {
		id := strconv.Itoa(t.TrackID)
		newIDs[id] = struct{}{}
		newIDSlice = append(newIDSlice, id)
	}

	for _, target := range targets {
		// Determine current IDs in the playlist.
		currentIDs := make(map[string]struct{}, len(target.Node.TRACK))
		for _, tr := range target.Node.TRACK {
			currentIDs[tr.Key] = struct{}{}
		}

		// Tracks to add: in new set but not currently in playlist.
		var toAdd []string
		for _, id := range newIDSlice {
			if _, exists := currentIDs[id]; !exists {
				toAdd = append(toAdd, id)
			}
		}

		// Tracks to remove: currently in playlist but not in new set.
		var toRemove []string
		for id := range currentIDs {
			if _, keep := newIDs[id]; !keep {
				toRemove = append(toRemove, id)
			}
		}

		if plDryRun {
			fmt.Printf("[Dry Run] %q: would add %d, remove %d tracks\n", target.Node.Name, len(toAdd), len(toRemove))
			continue
		}

		if len(toAdd) > 0 {
			syncEng.AddTracksToPlaylist(target.Node.Name, toAdd)
		}
		if len(toRemove) > 0 {
			syncEng.RemoveTracksFromPlaylist(target.Node.Name, toRemove)
		}
		fmt.Printf("Synced %q: +%d / -%d tracks\n", target.Node.Name, len(toAdd), len(toRemove))
	}

	if plDryRun {
		return nil
	}
	return syncEng.SaveXML(path)
}

func runPlaylistDelete(syncEng *syncpkg.Engine, targets []engine.NodeResult, path string) error {
	if plDryRun {
		for _, target := range targets {
			fmt.Printf("[Dry Run] Would delete playlist %q\n", target.Node.Name)
		}
		return nil
	}

	for _, target := range targets {
		if !syncEng.RemoveNode(target.Node.Name, 1) {
			fmt.Printf("Warning: could not delete playlist %q\n", target.Node.Name)
			continue
		}
		fmt.Printf("Deleted playlist %q\n", target.Node.Name)
	}
	return syncEng.SaveXML(path)
}

// loadXML resolves and loads the Rekordbox XML library, preferring --xml flag over config.
func loadXML() (*rekordbox.RekordboxLibraryXML, string, error) {
	cfg, _ := config.LoadAppConfig()
	path := utils.ExpandPath(xmlPath)
	if path == "" {
		path = utils.ExpandPath(cfg.RekordboxXMLPath)
	}
	if path == "" {
		return nil, "", fmt.Errorf("rekordbox XML path required; use --xml or run 'djlt config rekordbox --xml PATH'")
	}
	rbXML, err := rekordbox.ReadRekordboxLibrary(path)
	if err != nil {
		return nil, "", fmt.Errorf("failed to read rekordbox library: %w", err)
	}
	return rbXML, path, nil
}

// ── playlist fix subcommand ────────────────────────────────────────────────

var fixCmd = &cobra.Command{
	Use:   "fix [file...]",
	Short: "Fix playlist extensions and/or enrich with M3U8 metadata",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		for _, inputPath := range args {
			opts := playlist.FixOptions{
				Exts:           extsFlag,
				M3U8:           m3u8Flag,
				RemoveOriginal: removeOriginal,
				Force:          forceOverwrite,
				OutputPath:     outputFileFlag,
				Verbose:        verboseFlag,
				DryRun:         dryRunFlag,
			}

			if len(args) > 1 && outputFileFlag != "" {
				fmt.Printf("Warning: --output ignored when processing multiple files. Using default names for %s\n", inputPath)
				opts.OutputPath = ""
			}

			result, err := playlist.FixPlaylist(inputPath, opts)
			if err != nil {
				fmt.Printf("Error processing %s: %v\n", inputPath, err)
				continue
			}

			if opts.DryRun {
				fmt.Printf("DRY RUN: Would process '%s' -> '%s'\n", inputPath, result.OutputPath)
			} else {
				fmt.Printf("Successfully processed '%s' -> '%s'\n", inputPath, result.OutputPath)
			}
			fmt.Printf("Total tracks found: %d\n", result.TotalTracks-len(result.SkippedTracks))
			if len(result.SkippedTracks) > 0 {
				fmt.Printf("Skipped tracks (not found): %d\n", len(result.SkippedTracks))
				if verboseFlag {
					for _, p := range result.SkippedTracks {
						fmt.Printf("  - %s\n", p)
					}
				} else {
					fmt.Println("  (Use -v to see full list of skipped tracks)")
				}
			}

			if !opts.DryRun && removeOriginal && inputPath != result.OutputPath {
				fmt.Printf("\nRemove original file '%s'? (y/N): ", inputPath)
				var response string
				fmt.Scanln(&response)
				if response == "y" || response == "Y" {
					if err := os.Remove(inputPath); err != nil {
						return fmt.Errorf("failed to remove original file: %w", err)
					}
					fmt.Println("Original file removed.")
				} else {
					fmt.Println("Original file retained.")
				}
			}
			fmt.Println("---")
		}
		return nil
	},
}

func init() {
	// playlist fix flags
	fixCmd.Flags().StringSliceVarP(&extsFlag, "ext", "e", []string{}, "Priority list of file extensions (comma-separated)")
	fixCmd.Flags().BoolVar(&m3u8Flag, "m3u8", false, "Enrich playlist with M3U8 #EXTINF tags")
	fixCmd.Flags().BoolVarP(&removeOriginal, "remove-original", "r", false, "Remove the original playlist file after processing")
	fixCmd.Flags().BoolVarP(&forceOverwrite, "force", "f", false, "Force overwrite if output file exists")
	fixCmd.Flags().StringVarP(&outputFileFlag, "output", "o", "", "Specific output path (optional)")
	fixCmd.Flags().BoolVarP(&verboseFlag, "verbose", "v", false, "Enable verbose logging")
	fixCmd.Flags().BoolVar(&dryRunFlag, "dry-run", false, "Show what would be done without modifying files")

	// playlist god command flags
	playlistCmd.Flags().StringVar(&plAddFlag, "add", "", "Add tracks matching this rb/tracks: query")
	playlistCmd.Flags().StringVar(&plRemoveFlag, "remove", "", "Remove tracks matching this rb/tracks: query from matched playlists")
	playlistCmd.Flags().StringVar(&plSyncFlag, "sync", "", "Sync matched playlists to exactly match this rb/tracks: query")
	playlistCmd.Flags().StringVar(&plNewFlag, "new", "", "Create a new playlist with this name")
	playlistCmd.Flags().StringVar(&plRenameFlag, "rename", "", "Rename matched playlists to this name")
	playlistCmd.Flags().StringVar(&plMoveFlag, "move", "", "Move matched playlists into this folder")
	playlistCmd.Flags().BoolVar(&plDeleteFlag, "delete", false, "Delete matched playlists")
	playlistCmd.Flags().StringVar(&plFolderFlag, "folder", "", "Parent folder for --new (default: root level)")
	playlistCmd.Flags().IntVar(&plAtFlag, "at", -1, "Insert new playlist at this 0-indexed position (-1 for end)")
	playlistCmd.Flags().BoolVar(&plDryRun, "dry-run", false, "Preview changes without writing")

	playlistCmd.AddCommand(fixCmd)
	RootCmd.AddCommand(playlistCmd)
}
