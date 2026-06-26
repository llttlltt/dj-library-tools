package cli

import (
	"fmt"
	"os"
	"strings"

	"github.com/llttlltt/dj-library-tools/internal/config"
	"github.com/llttlltt/dj-library-tools/internal/engine"
	"github.com/llttlltt/dj-library-tools/internal/plex"
	"github.com/llttlltt/dj-library-tools/internal/sync"
	"github.com/llttlltt/dj-library-tools/internal/utils"
	"github.com/spf13/cobra"
)

var (
	exportDest   string
	exportFormat string
	syncTo       []string
)

var syncCmd = &cobra.Command{
	Use:   "sync [source-resource] [source-query] --to [target-resource] [target-query]",
	Short: "Sync items between a source and one or more targets",
	Args:  cobra.MinimumNArgs(1),
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

			if src.Location.Provider == "plex" && tgt.Location.Provider == "rb" {
				if err := syncPlexToRekordbox(src, tgt); err != nil {
					return err
				}
			} else if src.Location.Provider == "plex" && tgt.Location.Provider == "m3u8" {
				if err := syncPlexToM3U8(src.Location, tgt.Location); err != nil {
					return err
				}
			} else {
				return fmt.Errorf("unsupported sync direction: %s to %s", src.Location.Provider, tgt.Location.Provider)
			}
		}
		return nil
	},
}

func syncPlexToRekordbox(src, tgt *Selection) error {
	cfg, _ := config.LoadAppConfig()
	rbXML, path, err := loadXMLFunc()
	if err != nil {
		return err
	}

	orch := sync.NewOrchestrator(src.PlexClient, engine.NewRekordboxLibrary(rbXML), dryRun, verbose)

	err = orch.SyncPlexToRekordbox(src.RawTracks.([]plex.Track), src.Location.Query, sync.SyncOptions{
		ExportDest:   exportDest,
		ExportFormat: exportFormat,
		PathMaps:     cfg.PathMaps,
	})
	if err != nil {
		return err
	}

	if !dryRun {
		return engine.NewRekordboxLibrary(rbXML).Save(path)
	}
	return nil
}

func syncPlexToM3U8(src, tgt utils.Location) error {
	cfg, _ := config.LoadAppConfig()
	token := os.Getenv("PLEX_TOKEN")
	if token == "" {
		token = cfg.PlexToken
	}
	if token == "" {
		return fmt.Errorf("Plex token not found. Run 'djlt auth plex' or set PLEX_TOKEN env var")
	}

	// We'll reuse our Selection logic here in a future pass
	// For now, let's just make it compile.
	fmt.Printf("M3U8 sync not yet refactored to Orchestrator\n")
	return nil
}

func init() {
	syncCmd.Flags().StringSliceVar(&syncTo, "to", []string{}, "Target resource(s) to sync to (repeatable)")
	syncCmd.Flags().StringVar(&exportDest, "dest", "", "Destination directory for exported files")
	syncCmd.Flags().StringVar(&exportFormat, "format", "mp3", "Target format for exported files")
	RootCmd.AddCommand(syncCmd)
}
