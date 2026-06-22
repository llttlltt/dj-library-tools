package main

import (
	"fmt"
	"strings"

	"github.com/llttlltt/dj-library-tools/internal/config"
	"github.com/spf13/cobra"
)

var (
	cfgPlexHost      string
	cfgPlexPort      int
	cfgPlexToken     string
	cfgPlexMap       string
	cfgPlexRemoveMap string
	cfgRekordboxXML  string
)

var configCmd = &cobra.Command{
	Use:   "config [flags]",
	Short: "View or update application configuration",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, _ := config.LoadAppConfig()
		dirty := false

		if cfgPlexHost != "" {
			cfg.PlexHost = cfgPlexHost
			fmt.Printf("Plex host set to: %s\n", cfgPlexHost)
			dirty = true
		}

		if cmd.Flags().Changed("plex-port") {
			cfg.PlexPort = cfgPlexPort
			fmt.Printf("Plex port set to: %d\n", cfgPlexPort)
			dirty = true
		}

		if cfgPlexToken != "" {
			cfg.PlexToken = cfgPlexToken
			fmt.Printf("Plex token updated.\n")
			dirty = true
		}

		if cfgPlexMap != "" {
			parts := strings.SplitN(cfgPlexMap, ":", 2)
			if len(parts) != 2 {
				return fmt.Errorf("invalid --plex-map format; use remote:local")
			}
			if cfg.PathMaps == nil {
				cfg.PathMaps = make(map[string]string)
			}
			cfg.PathMaps[parts[0]] = parts[1]
			fmt.Printf("Added path map: %s -> %s\n", parts[0], parts[1])
			dirty = true
		}

		if cfgPlexRemoveMap != "" {
			if cfg.PathMaps != nil {
				if _, exists := cfg.PathMaps[cfgPlexRemoveMap]; exists {
					delete(cfg.PathMaps, cfgPlexRemoveMap)
					fmt.Printf("Removed path map for: %s\n", cfgPlexRemoveMap)
				} else {
					fmt.Printf("No path map found for: %s\n", cfgPlexRemoveMap)
				}
			}
			dirty = true
		}

		if cfgRekordboxXML != "" {
			cfg.RekordboxXMLPath = cfgRekordboxXML
			fmt.Printf("Rekordbox XML path set to: %s\n", cfgRekordboxXML)
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
	configCmd.Flags().StringVar(&cfgPlexHost, "plex-host", "", "Plex server host (e.g. 10.0.0.5)")
	configCmd.Flags().IntVar(&cfgPlexPort, "plex-port", 32400, "Plex server port (default: 32400)")
	configCmd.Flags().StringVar(&cfgPlexToken, "plex-token", "", "Plex authentication token")
	configCmd.Flags().StringVar(&cfgPlexMap, "plex-map", "", "Add a Plex path map (remote:local)")
	configCmd.Flags().StringVar(&cfgPlexRemoveMap, "plex-remove-map", "", "Remove a Plex path map by its remote key")
	configCmd.Flags().StringVar(&cfgRekordboxXML, "rekordbox-xml-path", "", "Path to the Rekordbox XML library")
	rootCmd.AddCommand(configCmd)
}
