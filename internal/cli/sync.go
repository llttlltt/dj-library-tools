package cli

import (
	"fmt"
	"strings"

	"github.com/llttlltt/dj-library-tools/internal/config"
	"github.com/llttlltt/dj-library-tools/internal/library"
	"github.com/llttlltt/dj-library-tools/internal/models"
	"github.com/llttlltt/dj-library-tools/internal/provider"
	"github.com/llttlltt/dj-library-tools/internal/sync"
	"github.com/llttlltt/dj-library-tools/internal/utils"
	"github.com/spf13/cobra"
)

func newSyncCmd() *cobra.Command {
	var syncTo []string
	var exportDest, exportFormat string
	var syncAppend bool

	cmd := &cobra.Command{
		Use:   "sync [source-resource] [source-query] --to [target-resource] [target-query]",
		Short: "Keep a playlist in sync with a track query",
		Long: `Synchronizes a target (like a Rekordbox playlist or M3U file) with a source query.

The sync command is "surgical"—it only adds or removes tracks necessary to make the target
match the source. By default, it removes tracks from the target that no longer match the query.

### Examples

**Keep an "Inbox" playlist matched to specific criteria:**
  djlt sync "rb/tracks added:>today" --to "rb/playlists name:Inbox"

**Add new tracks to a playlist without removing existing ones:**
  djlt sync "rb/tracks rating:5" --to "rb/playlists name:Favorites" --append

**Sync a query to an external M3U playlist file:**
  djlt sync "rb/tracks genre:House" --to "m3u/path/to/playlist.m3u"
`,
		Args: cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(syncTo) == 0 {
				return fmt.Errorf("at least one --to target is required")
			}

			queryOverride := ""
			if len(args) > 1 {
				queryOverride = strings.Join(args[1:], " ")
			}
			src, err := ResolveSelection(args[0], queryOverride)
			if err != nil {
				return err
			}

			for _, targetStr := range syncTo {
				tgt, err := ResolveSelection(targetStr, "")
				if err != nil {
					return err
				}

				if _, ok := tgt.Provider.(provider.WritableProvider); ok {
					if tgt.Location.Provider == "m3u" || tgt.Location.Provider == "m3u8" {
						if err := syncToM3U(src, tgt); err != nil {
							return err
						}
					} else {
						if err := syncToRekordbox(src, tgt, exportDest, exportFormat, syncAppend); err != nil {
							return err
						}
					}
				} else {
					return fmt.Errorf("unsupported sync target: %s → %s", src.Location.Provider, tgt.Location.Provider)
				}
			}
			return nil
		},
	}
	cmd.Flags().StringSliceVar(&syncTo, "to", []string{}, "Target resource(s) to sync to (repeatable)")
	cmd.Flags().StringVar(&exportDest, "dest", "", "Destination directory for exported files")
	cmd.Flags().StringVar(&exportFormat, "format", "mp3", "Target format for exported files")
	cmd.Flags().BoolVar(&syncAppend, "append", false, "Append new tracks without removing existing ones")
	return cmd
}

func syncToRekordbox(src, tgt *Selection, exportDest, exportFormat string, appendOnly bool) error {
	cfg, _ := config.LoadAppConfig()
	rbXML, path, err := loadXMLFunc()
	if err != nil {
		return err
	}

	orch := sync.NewOrchestrator(nil, library.NewRekordboxLibrary(rbXML), dryRun, verbose)

	tracks, err := src.Provider.GetTracks(src.Location.Query)
	if err != nil {
		return err
	}

	err = orch.SyncToLibrary(tracks, src.Location.Query, tgt.Location.Query, sync.SyncOptions{
		ExportDest:   exportDest,
		ExportFormat: exportFormat,
		PathMaps:     cfg.PathMaps,
	}, appendOnly)
	if err != nil {
		return err
	}

	if !dryRun {
		return library.NewRekordboxLibrary(rbXML).Save(path)
	}
	return nil
}

func syncToM3U(src, tgt *Selection) error {
	tracks, err := src.Provider.GetTracks(src.Location.Query)
	if err != nil {
		return err
	}

	if dryRun {
		fmt.Printf("[Dry Run] Would sync %d tracks to %s\n", len(tracks), tgt.Location.Resource)
		return nil
	}

	wp := tgt.Provider.(provider.WritableProvider)
	added, err := wp.AddTracks(models.ResourceGroup{}, tracks)
	if err != nil {
		return err
	}

	fmt.Printf("Synced %d tracks to %s\n", added, tgt.Location.Resource)
	return wp.Save(tgt.Location.Resource)
}

func syncPlexToM3U8(src, tgt utils.Location) error {
	fmt.Printf("M3U8 sync not yet refactored to Orchestrator\n")
	return nil
}
