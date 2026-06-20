package main

import (
	"fmt"
	"strings"

	"github.com/fatih/color"
	"github.com/llttlltt/dj-library-tools/internal/config"
	"github.com/llttlltt/dj-library-tools/internal/engine"
	"github.com/llttlltt/dj-library-tools/internal/plex"
	"github.com/llttlltt/dj-library-tools/internal/utils"
	"github.com/llttlltt/dj-library-tools/pkg/rekordbox"
	"github.com/spf13/cobra"
	"os"
)

var listCmd = &cobra.Command{
	Use:   "list [location]",
	Short: "List items from a location (e.g. plex/playlists:Summer)",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		loc := utils.ParseLocation(args[0])

		switch loc.Provider {
		case "rb", "rekordbox":
			return listRekordbox(loc)
		case "plex":
			return listPlex(loc)
		default:
			return fmt.Errorf("unknown provider: %s", loc.Provider)
		}
	},
}

func listRekordbox(loc utils.Location) error {
	if xmlPath == "" {
		return fmt.Errorf("XML path must be provided via --xml or -x")
	}

	lib, err := rekordbox.ReadRekordboxLibrary(xmlPath)
	if err != nil {
		return fmt.Errorf("failed to read XML: %w", err)
	}

	eng := engine.NewEngine(lib)
	
	if loc.Resource == "playlists" {
		// Logic to list RB playlists
		for _, node := range lib.Playlists.Node.Node {
			if loc.Query == "" || strings.Contains(strings.ToLower(node.Name), strings.ToLower(loc.Query)) {
				fmt.Printf("Playlist: %s (%d entries)\n", node.Name, node.Entries)
			}
		}
		return nil
	}

	tracks, err := eng.Ls(loc.Query)
	if err != nil {
		return fmt.Errorf("ls failed: %w", err)
	}

	if len(tracks) == 0 {
		color.Yellow("No tracks matched the query.")
		return nil
	}

	headerFmt := color.New(color.FgCyan, color.Bold, color.Underline).SprintfFunc()
	dimFmt := color.New(color.FgHiBlack).SprintfFunc()
	artistFmt := color.New(color.FgHiMagenta).SprintfFunc()
	titleFmt := color.New(color.FgHiWhite).SprintfFunc()
	bpmFmt := color.New(color.FgHiGreen).SprintfFunc()
	keyFmt := color.New(color.FgHiYellow).SprintfFunc()

	fmt.Printf("%s%s %s%s %s\n", "   ", headerFmt("BPM"), " ", headerFmt("Key"), headerFmt("Artist - Title"))
	for _, t := range tracks {
		bpm := 0.0
		if len(t.Tempo) > 0 {
			bpm = t.Tempo[0].Bpm
		}
		fmt.Printf("%s %s %s %s %s\n", 
			bpmFmt("%6.2f", bpm), 
			keyFmt("%4s", t.Tonality),
			artistFmt(t.Artist), 
			dimFmt("-"), 
			titleFmt(t.Name))
	}

	fmt.Printf("\n%s\n", color.HiGreenString("Matched %d tracks.", len(tracks)))
	return nil
}

func listPlex(loc utils.Location) error {
	token := os.Getenv("PLEX_TOKEN")
	if token == "" {
		cfg, _ := config.LoadAppConfig()
		token = cfg.PlexToken
	}
	if token == "" {
		return fmt.Errorf("Plex token not found. Run 'djlt auth plex' or set PLEX_TOKEN env var")
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
		if len(res.Connections) == 0 {
			continue
		}

		baseURL := res.Connections[0].URI
		serverClient := plex.NewClient(res.AccessToken)

		if loc.Resource == "playlists" {
			playlists, err := serverClient.GetPlaylists(baseURL)
			if err != nil {
				fmt.Printf("  Failed to get playlists: %v\n", err)
				continue
			}

			for _, pl := range playlists {
				if loc.Query != "" && !strings.Contains(strings.ToLower(pl.Title), strings.ToLower(loc.Query)) {
					continue
				}
				fmt.Printf("  - [%s] %s (%d tracks)\n", pl.RatingKey, pl.Title, pl.LeafCount)
			}
		} else if loc.Resource == "tracks" {
			// This would ideally fetch all tracks and apply the query
			fmt.Printf("  (Listing all tracks not yet implemented, use playlists)\n")
		}
	}
	return nil
}

func init() {
	rootCmd.AddCommand(listCmd)
}
