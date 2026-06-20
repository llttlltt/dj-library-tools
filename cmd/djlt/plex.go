package main

import (
	"fmt"
	"os"
	"time"

	"github.com/llttlltt/dj-library-tools/internal/plex"
	"github.com/spf13/cobra"
)

var plexCmd = &cobra.Command{
	Use:   "plex",
	Short: "Plex integration commands",
}

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Authenticate with Plex using PIN flow",
	RunE: func(cmd *cobra.Command, args []string) error {
		client := plex.NewClient("")
		pin, err := client.RequestPin()
		if err != nil {
			return fmt.Errorf("failed to request pin: %w", err)
		}

		fmt.Printf("Please visit: https://app.plex.tv/auth/#!?pin=%s&clientID=%s\n", pin.Code, plex.AppID)
		fmt.Printf("Waiting for authentication...\n")

		ticker := time.NewTicker(2 * time.Second)
		defer ticker.Stop()

		timeout := time.After(5 * time.Minute)

		for {
			select {
			case <-timeout:
				return fmt.Errorf("authentication timed out")
			case <-ticker.C:
				status, err := client.CheckPin(pin.ID)
				if err != nil {
					return fmt.Errorf("failed to check pin status: %w", err)
				}

				if status.AuthToken != "" {
					fmt.Printf("Successfully authenticated!\n")
					fmt.Printf("Your Plex Token: %s\n", status.AuthToken)
					fmt.Printf("Save this token in your environment as PLEX_TOKEN or use it in subsequent commands.\n")
					return nil
				}
			}
		}
	},
}

var plexLsCmd = &cobra.Command{
	Use:   "ls",
	Short: "List Plex playlists",
	RunE: func(cmd *cobra.Command, args []string) error {
		token := os.Getenv("PLEX_TOKEN")
		if token == "" {
			return fmt.Errorf("PLEX_TOKEN environment variable not set")
		}

		client := plex.NewClient(token)
		resources, err := client.GetResources()
		if err != nil {
			return fmt.Errorf("failed to get resources: %w", err)
		}

		for _, res := range resources {
			if res.Provides != "server" {
				continue
			}

			fmt.Printf("Server: %s\n", res.Name)
			// Use the first connection for simplicity in CLI list
			if len(res.Connections) == 0 {
				continue
			}

			baseURL := res.Connections[0].URI
			// If we are using a resource-specific token
			serverClient := plex.NewClient(res.AccessToken)

			playlists, err := serverClient.GetPlaylists(baseURL)
			if err != nil {
				fmt.Printf("  Failed to get playlists: %v\n", err)
				continue
			}

			for _, pl := range playlists {
				fmt.Printf("  - [%s] %s (%d tracks)\n", pl.RatingKey, pl.Title, pl.LeafCount)
			}
		}

		return nil
	},
}

func init() {
	plexCmd.AddCommand(loginCmd)
	plexCmd.AddCommand(plexLsCmd)
	rootCmd.AddCommand(plexCmd)
}
