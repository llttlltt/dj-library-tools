package main

import (
	"fmt"
	"time"

	"github.com/llttlltt/dj-library-tools/internal/plex"
	"github.com/spf13/cobra"
)

var authCmd = &cobra.Command{
	Use:   "auth",
	Short: "Authenticate with third-party providers",
}

var plexAuthCmd = &cobra.Command{
	Use:   "plex",
	Short: "Authenticate with Plex using PIN flow",
	RunE: func(cmd *cobra.Command, args []string) error {
		client := plex.NewClient("")
		pin, err := client.RequestPin()
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

func init() {
	authCmd.AddCommand(plexAuthCmd)
	rootCmd.AddCommand(authCmd)
}
