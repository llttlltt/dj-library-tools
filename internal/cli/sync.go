package cli

import (
	"fmt"
	"strings"

	"github.com/llttlltt/dj-library-tools/internal/config"
	"github.com/llttlltt/dj-library-tools/internal/engine"
	"github.com/llttlltt/dj-library-tools/internal/provider"
	"github.com/llttlltt/dj-library-tools/internal/sync"
	"github.com/llttlltt/dj-library-tools/internal/utils"
	"github.com/spf13/cobra"
)

var (
	exportDest   string
	exportFormat string
	syncTo       []string
	syncAppend   bool
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

			if _, ok := tgt.Provider.(provider.WritableProvider); ok {
				if err := syncToRekordbox(src, tgt); err != nil {
					return err
				}
			} else if tgt.Location.Provider == "m3u8" {
				if err := syncPlexToM3U8(src.Location, tgt.Location); err != nil {
					return err
				}
			} else {
				return fmt.Errorf("unsupported sync target: %s → %s", src.Location.Provider, tgt.Location.Provider)
			}
		}
		return nil
	},
}

func syncToRekordbox(src, tgt *Selection) error {
	cfg, _ := config.LoadAppConfig()
	rbXML, path, err := loadXMLFunc()
	if err != nil {
		return err
	}

	orch := sync.NewOrchestrator(nil, engine.NewRekordboxLibrary(rbXML), dryRun, verbose)

	tracks, err := src.Provider.GetTracks(src.Location.Query)
	if err != nil {
		return err
	}

	err = orch.SyncToLibrary(tracks, src.Location.Query, tgt.Location.Query, sync.SyncOptions{
		ExportDest:   exportDest,
		ExportFormat: exportFormat,
		PathMaps:     cfg.PathMaps,
	}, syncAppend)
	if err != nil {
		return err
	}

	if !dryRun {
		return engine.NewRekordboxLibrary(rbXML).Save(path)
	}
	return nil
}

func syncPlexToM3U8(src, tgt utils.Location) error {
	fmt.Printf("M3U8 sync not yet refactored to Orchestrator\n")
	return nil
}

func init() {
	syncCmd.Flags().StringSliceVar(&syncTo, "to", []string{}, "Target resource(s) to sync to (repeatable)")
	syncCmd.Flags().StringVar(&exportDest, "dest", "", "Destination directory for exported files")
	syncCmd.Flags().StringVar(&exportFormat, "format", "mp3", "Target format for exported files")
	syncCmd.Flags().BoolVar(&syncAppend, "append", false, "Append new tracks without removing existing ones")
	RootCmd.AddCommand(syncCmd)
}
