package main

import (
	"fmt"
	"strings"

	"github.com/fatih/color"
	"github.com/llttlltt/dj-library-tools/internal/engine"
	"github.com/llttlltt/dj-library-tools/internal/plex"
	"github.com/llttlltt/dj-library-tools/pkg/rekordbox"
	"github.com/spf13/cobra"
	"os"
)

var listCmd = &cobra.Command{
	Use:   "list [source] [query]",
	Short: "List items from a source matching a query",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		source := args[0]
		queryStr := ""
		if len(args) > 1 {
			queryStr = strings.Join(args[1:], " ")
		}

		switch source {
		case "rekordbox":
			return listRekordbox(queryStr)
		case "plex":
			return listPlex(queryStr)
		default:
			return fmt.Errorf("unknown source: %s (supported: rekordbox, plex)", source)
		}
	},
}

func listRekordbox(queryStr string) error {
	if xmlPath == "" {
		return fmt.Errorf("XML path must be provided via --xml or -x")
	}

	lib, err := rekordbox.ReadRekordboxLibrary(xmlPath)
	if err != nil {
		return fmt.Errorf("failed to read XML: %w", err)
	}

	eng := engine.NewEngine(lib)
	tracks, err := eng.Ls(queryStr)
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

func listPlex(queryStr string) error {
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
		if len(res.Connections) == 0 {
			continue
		}

		baseURL := res.Connections[0].URI
		serverClient := plex.NewClient(res.AccessToken)

		playlists, err := serverClient.GetPlaylists(baseURL)
		if err != nil {
			fmt.Printf("  Failed to get playlists: %v\n", err)
			continue
		}

		for _, pl := range playlists {
			// Simple substring match for now, will integrate Query Engine later
			if queryStr != "" && !strings.Contains(strings.ToLower(pl.Title), strings.ToLower(queryStr)) {
				continue
			}
			fmt.Printf("  - [%s] %s (%d tracks)\n", pl.RatingKey, pl.Title, pl.LeafCount)
		}
	}
	return nil
}

func init() {
	rootCmd.AddCommand(listCmd)
}
