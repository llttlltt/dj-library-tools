package main

import (
	"fmt"
	"strings"

	"github.com/llttlltt/dj-library-tools/internal/config"
	"github.com/spf13/cobra"
)

var (
	cfgHost      string
	cfgPort      int
	cfgToken     string
	cfgMap       string
	cfgRemoveMap string
	cfgXMLPath   string
)

var configCmd = &cobra.Command{
	Use:   "config [flags]",
	Short: "View or update application configuration",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, _ := config.LoadAppConfig()
		dirty := false

		if cfgHost != "" {
			cfg.PlexHost = cfgHost
			fmt.Printf("Plex host set to: %s\n", cfgHost)
			dirty = true
		}

		if cmd.Flags().Changed("port") {
			cfg.PlexPort = cfgPort
			fmt.Printf("Plex port set to: %d\n", cfgPort)
			dirty = true
		}

		if cfgToken != "" {
			cfg.PlexToken = cfgToken
			fmt.Printf("Plex token updated.\n")
			dirty = true
		}

		if cfgMap != "" {
			parts := strings.SplitN(cfgMap, ":", 2)
			if len(parts) != 2 {
				return fmt.Errorf("invalid --map format; use remote:local")
			}
			if cfg.PathMaps == nil {
				cfg.PathMaps = make(map[string]string)
			}
			cfg.PathMaps[parts[0]] = parts[1]
			fmt.Printf("Added path map: %s -> %s\n", parts[0], parts[1])
			dirty = true
		}

		if cfgRemoveMap != "" {
			if cfg.PathMaps != nil {
				if _, exists := cfg.PathMaps[cfgRemoveMap]; exists {
					delete(cfg.PathMaps, cfgRemoveMap)
					fmt.Printf("Removed path map for: %s\n", cfgRemoveMap)
				} else {
					fmt.Printf("No path map found for: %s\n", cfgRemoveMap)
				}
			}
			dirty = true
		}

		if cfgXMLPath != "" {
			cfg.RekordboxXMLPath = cfgXMLPath
			fmt.Printf("Rekordbox XML path set to: %s\n", cfgXMLPath)
			dirty = true
		}

		if dirty {
			if err := config.SaveAppConfig(cfg); err != nil {
				return fmt.Errorf("failed to save config: %w", err)
			}
		}

		fmt.Printf("\nDJ Library Tools Config:\n")
		fmt.Printf("  Plex Host:   %s\n", cfg.PlexHost)
		fmt.Printf("  Plex Port:   %d\n", cfg.PlexPort)
		fmt.Printf("  Plex Token:  %s\n", maskToken(cfg.PlexToken))
		fmt.Printf("  RB XML Path: %s\n", cfg.RekordboxXMLPath)
		if len(cfg.PathMaps) > 0 {
			fmt.Printf("  Path Maps:\n")
			for k, v := range cfg.PathMaps {
				fmt.Printf("    %s -> %s\n", k, v)
			}
		}
		return nil
	},
}

func maskToken(t string) string {
	if len(t) < 8 {
		return "****"
	}
	return t[:4] + "...." + t[len(t)-4:]
}

func init() {
	configCmd.Flags().StringVar(&cfgHost, "host", "", "Plex server host (e.g. 10.0.0.5)")
	configCmd.Flags().IntVar(&cfgPort, "port", 32400, "Plex server port (default: 32400)")
	configCmd.Flags().StringVar(&cfgToken, "token", "", "Plex authentication token")
	configCmd.Flags().StringVar(&cfgMap, "map", "", "Add a path map (remote:local)")
	configCmd.Flags().StringVar(&cfgRemoveMap, "remove-map", "", "Remove a path map by its remote key")
	configCmd.Flags().StringVar(&cfgXMLPath, "xml-path", "", "Path to the Rekordbox XML library")
	rootCmd.AddCommand(configCmd)
}
