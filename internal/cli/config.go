package cli

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/llttlltt/dj-library-tools/internal/config"
	"github.com/spf13/cobra"
)

func newConfigCmd() *cobra.Command {
	var cfgUnset, cfgList bool

	cmd := &cobra.Command{
	Use:   "config [key] [value]",
	Short: "View or update application configuration",
	Long: `Manage djlt configuration using dot-namespaced keys. Settings are stored in ~/.config/djlt/config.json.

## Keys

- **plex.host**: Plex server hostname or IP.
- **plex.port**: Plex server port (default: 32400).
- **plex.token**: Plex authentication token (usually set via 'djlt auth --plex').
- **plex.map**: Remote-to-local path map entry. Used to bridge Plex remote paths to your local mount points.
- **rekordbox.file-path**: Absolute path to your Rekordbox XML export file.

## Examples

**List all settings**
  djlt config --list

**Configure Rekordbox library**
  djlt config rekordbox.file-path ~/Documents/rekordbox.xml

**Set up Plex connection**
  djlt config plex.host 192.168.1.50
  djlt config plex.port 32400

**Add a Plex path mapping**
  djlt config plex.map /music/remote:/Volumes/Music

**Remove a specific path mapping**
  djlt config --unset plex.map /music/remote

**Unset a scalar value**
  djlt config --unset plex.host`,
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, _ := config.LoadAppConfig()

			if cfgList || len(args) == 0 {
			printConfig(cfg)
			return nil
		}

			key := args[0]

			if cfgUnset {
				return runConfigUnset(cfg, key, args[1:])
			}

			if len(args) == 1 {
				return runConfigGet(cfg, key)
			}

			return runConfigSet(cfg, key, args[1])
		},
	}
	cmd.Flags().BoolVar(&cfgList, "list", false, "Show all configuration values")
	cmd.Flags().BoolVar(&cfgUnset, "unset", false, "Remove a configuration value")
	return cmd
}

func runConfigSet(cfg *config.AppConfig, key, value string) error {
	switch key {
	case "plex.host":
		cfg.PlexHost = value
		fmt.Printf("plex.host = %s\n", value)
	case "plex.port":
		port, err := strconv.Atoi(value)
		if err != nil {
			return fmt.Errorf("plex.port must be an integer, got %q", value)
		}
		cfg.PlexPort = port
		fmt.Printf("plex.port = %d\n", port)
	case "plex.token":
		cfg.PlexToken = value
		fmt.Println("plex.token updated")
	case "plex.map":
		parts := strings.SplitN(value, ":", 2)
		if len(parts) != 2 {
			return fmt.Errorf("plex.map value must be remote:local, got %q", value)
		}
		if cfg.PathMaps == nil {
			cfg.PathMaps = make(map[string]string)
		}
		cfg.PathMaps[parts[0]] = parts[1]
		fmt.Printf("plex.map %s -> %s\n", parts[0], parts[1])
	case "rekordbox.file-path":
		cfg.PrimaryFilePath = value
		fmt.Printf("rekordbox.file-path = %s\n", value)
	default:
		return fmt.Errorf("unknown config key %q; run 'djlt config --help' for valid keys", key)
	}
	return config.SaveAppConfig(cfg)
}

func runConfigGet(cfg *config.AppConfig, key string) error {
	switch key {
	case "plex.host":
		fmt.Println(cfg.PlexHost)
	case "plex.port":
		fmt.Println(cfg.PlexPort)
	case "plex.token":
		fmt.Println(maskToken(cfg.PlexToken))
	case "plex.map":
		for remote, local := range cfg.PathMaps {
			fmt.Printf("%s:%s\n", remote, local)
		}
	case "rekordbox.file-path":
		fmt.Println(cfg.PrimaryFilePath)
	default:
		return fmt.Errorf("unknown config key %q; run 'djlt config --help' for valid keys", key)
	}
	return nil
}

func runConfigUnset(cfg *config.AppConfig, key string, rest []string) error {
	switch key {
	case "plex.host":
		cfg.PlexHost = ""
		fmt.Println("unset plex.host")
	case "plex.port":
		cfg.PlexPort = 0
		fmt.Println("unset plex.port")
	case "plex.token":
		cfg.PlexToken = ""
		fmt.Println("unset plex.token")
	case "plex.map":
		if len(rest) == 0 {
			return fmt.Errorf("--unset plex.map requires the remote path to remove")
		}
		remote := rest[0]
		if cfg.PathMaps == nil {
			return fmt.Errorf("no path maps configured")
		}
		if _, exists := cfg.PathMaps[remote]; !exists {
			return fmt.Errorf("no path map found for %q", remote)
		}
		delete(cfg.PathMaps, remote)
		fmt.Printf("removed plex.map %s\n", remote)
	case "rekordbox.file-path":
		cfg.PrimaryFilePath = ""
		fmt.Println("unset rekordbox.file-path")
	default:
		return fmt.Errorf("unknown config key %q; run 'djlt config --help' for valid keys", key)
	}
	return config.SaveAppConfig(cfg)
}

func printConfig(cfg *config.AppConfig) {
	fmt.Printf("plex.host = %s\n", cfg.PlexHost)
	fmt.Printf("plex.port = %d\n", cfg.PlexPort)
	fmt.Printf("plex.token = %s\n", maskToken(cfg.PlexToken))
	fmt.Printf("rekordbox.file-path = %s\n", cfg.PrimaryFilePath)
	for remote, local := range cfg.PathMaps {
		fmt.Printf("plex.map = %s:%s\n", remote, local)
	}
}

func maskToken(t string) string {
	if len(t) < 8 {
		return "****"
	}
	return t[:4] + "...." + t[len(t)-4:]
}


