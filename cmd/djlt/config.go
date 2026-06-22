package main

import (
	"fmt"
	"strings"

	"github.com/llttlltt/dj-library-tools/internal/config"
	"github.com/spf13/cobra"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage application configuration",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, _ := config.LoadAppConfig()
		fmt.Printf("DJ Library Tools Config:\n")
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

var plexConfigCmd = &cobra.Command{
	Use:   "plex",
	Short: "Configure Plex settings",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, _ := config.LoadAppConfig()

		if plexHost != "" {
			cfg.PlexHost = plexHost
			fmt.Printf("Plex host set to: %s\n", plexHost)
		}

		if cmd.Flags().Changed("port") {
			cfg.PlexPort = plexPort
			fmt.Printf("Plex port set to: %d\n", plexPort)
		}

		if plexToken != "" {
			cfg.PlexToken = plexToken
			fmt.Printf("Plex token updated.\n")
		}

		if pathMap != "" {
			parts := strings.SplitN(pathMap, ":", 2)
			if len(parts) == 2 {
				if cfg.PathMaps == nil {
					cfg.PathMaps = make(map[string]string)
				}
				cfg.PathMaps[parts[0]] = parts[1]
				fmt.Printf("Added path map: %s -> %s\n", parts[0], parts[1])
			} else {
				return fmt.Errorf("invalid map format. Use remote:local")
			}
		}

		if removeMap != "" {
			if cfg.PathMaps != nil {
				if _, exists := cfg.PathMaps[removeMap]; exists {
					delete(cfg.PathMaps, removeMap)
					fmt.Printf("Removed path map for: %s\n", removeMap)
				} else {
					fmt.Printf("No path map found for: %s\n", removeMap)
				}
			}
		}

		if err := config.SaveAppConfig(cfg); err != nil {
			return fmt.Errorf("failed to save config: %w", err)
		}

		fmt.Printf("Current Plex Config:\n")
		fmt.Printf("  Host:   %s\n", cfg.PlexHost)
		fmt.Printf("  Port:   %d\n", cfg.PlexPort)
		fmt.Printf("  Token:  %s\n", maskToken(cfg.PlexToken))
		if len(cfg.PathMaps) > 0 {
			fmt.Printf("  Maps:\n")
			for k, v := range cfg.PathMaps {
				fmt.Printf("    %s -> %s\n", k, v)
			}
		}

		return nil
	},
}

var rbConfigCmd = &cobra.Command{
	Use:   "rekordbox",
	Short: "Configure Rekordbox settings",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, _ := config.LoadAppConfig()

		if rbXMLPath != "" {
			cfg.RekordboxXMLPath = rbXMLPath
			fmt.Printf("Rekordbox XML path set to: %s\n", rbXMLPath)
		}

		if err := config.SaveAppConfig(cfg); err != nil {
			return fmt.Errorf("failed to save config: %w", err)
		}

		fmt.Printf("Current Rekordbox Config:\n")
		fmt.Printf("  XML: %s\n", cfg.RekordboxXMLPath)

		return nil
	},
}

var (
	plexHost  string
	plexPort  int
	plexToken string
	pathMap   string
	removeMap string
	rbXMLPath string
)

func maskToken(t string) string {
	if len(t) < 8 {
		return "****"
	}
	return t[:4] + "...." + t[len(t)-4:]
}

func init() {
	plexConfigCmd.Flags().StringVar(&plexHost, "host", "", "Host of the Plex server (e.g. 10.0.10.151)")
	plexConfigCmd.Flags().IntVar(&plexPort, "port", 32400, "Port of the Plex server (default: 32400)")
	plexConfigCmd.Flags().StringVar(&plexToken, "token", "", "Plex authentication token")
	plexConfigCmd.Flags().StringVar(&pathMap, "map", "", "Map a remote path to a local path (remote:local)")
	plexConfigCmd.Flags().StringVar(&removeMap, "remove-map", "", "Remove a path map by its remote path key")
	configCmd.AddCommand(plexConfigCmd)

	rbConfigCmd.Flags().StringVar(&rbXMLPath, "xml", "", "Path to the Rekordbox XML library")
	configCmd.AddCommand(rbConfigCmd)

	rootCmd.AddCommand(configCmd)
}
