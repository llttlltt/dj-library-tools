package cli

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/llttlltt/dj-library-tools/internal/config"
	"github.com/llttlltt/dj-library-tools/internal/providers/plex"
	"github.com/spf13/cobra"
)

func newConfigCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "config",
		Short: "Manage application configuration",
	}

	cmd.AddCommand(
		newConfigListCmd(),
		newConfigPlexCmd(),
		newConfigRekordboxCmd(),
	)

	return cmd
}

func newConfigListCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "Show all configuration values",
		Run: func(cmd *cobra.Command, args []string) {
			cfg, _ := config.LoadAppConfig()
			fmt.Printf("plex.host = %s\n", cfg.Plex.Host)
			fmt.Printf("plex.port = %d\n", cfg.Plex.Port)
			fmt.Printf("plex.token = %s\n", maskToken(cfg.Plex.Token))
			fmt.Printf("rekordbox.primary_file_path = %s\n", cfg.Rekordbox.PrimaryFilePath)
			for remote, local := range cfg.PathMaps {
				fmt.Printf("path_map = %s:%s\n", remote, local)
			}
		},
	}
}

func newConfigPlexCmd() *cobra.Command {
	plexCmd := &cobra.Command{
		Use:   "plex",
		Short: "Configure Plex settings",
	}

	plexCmd.AddCommand(
		&cobra.Command{
			Use:   "host [value]",
			Short: "Set or get Plex host",
			Args:  cobra.MaximumNArgs(1),
			RunE: func(cmd *cobra.Command, args []string) error {
				cfg, _ := config.LoadAppConfig()
				if len(args) == 0 {
					fmt.Println(cfg.Plex.Host)
					return nil
				}
				cfg.Plex.Host = args[0]
				return config.SaveAppConfig(cfg)
			},
		},
		&cobra.Command{
			Use:   "port [value]",
			Short: "Set or get Plex port",
			Args:  cobra.MaximumNArgs(1),
			RunE: func(cmd *cobra.Command, args []string) error {
				cfg, _ := config.LoadAppConfig()
				if len(args) == 0 {
					fmt.Println(cfg.Plex.Port)
					return nil
				}
				port, err := strconv.Atoi(args[0])
				if err != nil {
					return fmt.Errorf("port must be an integer")
				}
				cfg.Plex.Port = port
				return config.SaveAppConfig(cfg)
			},
		},
		&cobra.Command{
			Use:   "auth",
			Short: "Interactive Plex authentication (PIN flow)",
			RunE: func(cmd *cobra.Command, args []string) error {
				return runPlexAuth()
			},
		},
		&cobra.Command{
			Use:   "map [remote:local]",
			Short: "Add a path mapping",
			Args:  cobra.ExactArgs(1),
			RunE: func(cmd *cobra.Command, args []string) error {
				parts := strings.SplitN(args[0], ":", 2)
				if len(parts) != 2 {
					return fmt.Errorf("format must be remote:local")
				}
				cfg, _ := config.LoadAppConfig()
				if cfg.PathMaps == nil {
					cfg.PathMaps = make(map[string]string)
				}
				cfg.PathMaps[parts[0]] = parts[1]
				return config.SaveAppConfig(cfg)
			},
		},
	)

	return plexCmd
}

func newConfigRekordboxCmd() *cobra.Command {
	rbCmd := &cobra.Command{
		Use:   "rb",
		Short: "Configure Rekordbox settings",
	}

	rbCmd.AddCommand(
		&cobra.Command{
			Use:   "file [path]",
			Short: "Set or get primary Rekordbox XML path",
			Args:  cobra.MaximumNArgs(1),
			RunE: func(cmd *cobra.Command, args []string) error {
				cfg, _ := config.LoadAppConfig()
				if len(args) == 0 {
					fmt.Println(cfg.Rekordbox.PrimaryFilePath)
					return nil
				}
				cfg.Rekordbox.PrimaryFilePath = args[0]
				return config.SaveAppConfig(cfg)
			},
		},
	)

	return rbCmd
}

func runPlexAuth() error {
	client := plex.NewClient("")
	pin, err := client.RequestPin(context.Background())
	if err != nil {
		return fmt.Errorf("failed to request pin: %w", err)
	}

	fmt.Printf("Please visit: https://app.plex.tv/auth/#!?code=%s&context%%5Bdevice%%5D%%5Bproduct%%5D=%s&clientID=%s\n", pin.Code, plex.ClientName, plex.AppID)
	fmt.Printf("Waiting for authentication...\n")

	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()
	timeout := time.After(5 * time.Minute)

	for {
		select {
		case <-timeout:
			return fmt.Errorf("authentication timed out")
		case <-ticker.C:
			status, err := client.CheckPin(context.Background(), pin.ID)
			if err != nil {
				return fmt.Errorf("failed to check pin status: %w", err)
			}

			if status.AuthToken != "" {
				fmt.Printf("Successfully authenticated!\n")
				cfg, _ := config.LoadAppConfig()
				cfg.Plex.Token = status.AuthToken
				return config.SaveAppConfig(cfg)
			}
		}
	}
}

func maskToken(t string) string {
	if len(t) < 8 {
		return "****"
	}
	return t[:4] + "...." + t[len(t)-4:]
}
