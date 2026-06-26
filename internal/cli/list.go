package cli

import (
	"context"
	"fmt"
	"strconv"
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
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		var loc utils.Location
		if len(args) > 1 {
			loc = utils.ParseLocation(args[0], strings.Join(args[1:], " "))
		} else {
			loc = utils.ParseLocation(args[0], "")
		}

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
	cfg, _ := config.LoadAppConfig()
	path := utils.ExpandPath(xmlPath)
	if path == "" {
		path = utils.ExpandPath(cfg.RekordboxXMLPath)
	}
	if path == "" {
		return fmt.Errorf("Rekordbox XML path not found. Use --xml or run 'djlt config rekordbox --xml PATH'")
	}

	lib, err := rekordbox.ReadRekordboxLibrary(path)
	if err != nil {
		return fmt.Errorf("failed to read XML: %w", err)
	}

	eng := engine.NewEngine(lib)
	
	if loc.Resource == "playlists" {
		results, err := eng.LsPlaylists(loc.Query)
		if err != nil {
			return fmt.Errorf("ls failed: %w", err)
		}
		if len(results) == 0 {
			color.Yellow("No playlists matched the query.")
			return nil
		}
		for _, res := range results {
			name := res.Node.Name
			if res.ParentFolder != "" {
				name = res.ParentFolder + "/" + name
			}
			fmt.Printf("Playlist: %s (%d entries)\n", name, rekordbox.DerefInt32(res.Node.Entries))
		}
		return nil
	}

	if loc.Resource == "folders" {
		results, err := eng.LsFolders(loc.Query)
		if err != nil {
			return fmt.Errorf("ls failed: %w", err)
		}
		if len(results) == 0 {
			color.Yellow("No folders matched the query.")
			return nil
		}
		for _, res := range results {
			name := res.Node.Name
			if res.ParentFolder != "" {
				name = res.ParentFolder + "/" + name
			}
			fmt.Printf("Folder: %s (%d entries)\n", name, rekordbox.DerefInt32(res.Node.Count))
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
			bpm, _ = strconv.ParseFloat(t.Tempo[0].Bpm, 64)
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

// printPlaylists recursively walks the node tree and prints matching playlists.
func printPlaylists(nodes []rekordbox.Node, query, folderPath string) {
	for _, node := range nodes {
		if node.Type == 1 { // playlist
			name := node.Name
			if folderPath != "" {
				name = folderPath + "/" + name
			}
			if query == "" || strings.Contains(strings.ToLower(node.Name), strings.ToLower(query)) {
				fmt.Printf("Playlist: %s (%d entries)\n", name, rekordbox.DerefInt32(node.Entries))
			}
		} else if node.Type == 0 { // folder
			nextPath := node.Name
			if folderPath != "" {
				nextPath = folderPath + "/" + node.Name
			}
			if len(node.Node) > 0 {
				printPlaylists(node.Node, query, nextPath)
			}
		}
	}
}

func listPlex(loc utils.Location) error {
	cfg, _ := config.LoadAppConfig()
	token := os.Getenv("PLEX_TOKEN")
	if token == "" {
		token = cfg.PlexToken
	}
	if token == "" {
		return fmt.Errorf("Plex token not found. Run 'djlt auth plex' or set PLEX_TOKEN env var")
	}

	client := plex.NewClient(token)
	resources, err := client.GetResources(context.Background())
	if err != nil {
		return fmt.Errorf("failed to get resources: %w", err)
	}

	for _, res := range resources {
		if res.Provides != "server" {
			continue
		}

		fmt.Printf("Server: %s [%s]\n", res.Name, res.ClientIdentifier)
		serverClient := plex.NewClient(res.AccessToken)

		if loc.Resource == "playlists" {
			var probe *plex.ConnectionResult
			var lastErr error

			if cfg.PlexHost != "" {
				port := cfg.PlexPort
				if port == 0 {
					port = 32400
				}
				baseURL := fmt.Sprintf("http://%s:%d", cfg.PlexHost, port)
				playlists, err := serverClient.GetPlaylists(context.Background(), baseURL)
				if err == nil {
					probe = &plex.ConnectionResult{BaseURL: baseURL, Playlists: playlists}
				} else {
					lastErr = err
				}
			} else {
				probe, lastErr = serverClient.ProbeBestConnection(res)
			}

			if lastErr != nil {
				fmt.Printf("  Failed to connect: %v\n", lastErr)
				continue
			}

			fmt.Printf("  Connected via: %s\n", probe.BaseURL)
			for _, pl := range probe.Playlists {
				if loc.Query != "" && !strings.Contains(strings.ToLower(pl.Title), strings.ToLower(loc.Query)) {
					continue
				}
				fmt.Printf("  - [%s] %s (%d tracks)\n", pl.RatingKey, pl.Title, pl.LeafCount)
			}
		} else if loc.Resource == "tracks" {
			var probe *plex.ConnectionResult
			var lastErr error

			if cfg.PlexHost != "" {
				port := cfg.PlexPort
				if port == 0 {
					port = 32400
				}
				baseURL := fmt.Sprintf("http://%s:%d", cfg.PlexHost, port)
				// We need to find the playlist to get its tracks, or just fetch library
				// For now, let's assume query is a playlist ID
				tracks, err := serverClient.GetPlaylistTracks(context.Background(), baseURL, "/playlists/"+loc.Query+"/items")
				if err == nil {
					probe = &plex.ConnectionResult{BaseURL: baseURL, Tracks: tracks}
				} else {
					lastErr = err
				}
			} else {
				// Probing for tracks is complex, for now we require a host or use first connection
				if len(res.Connections) > 0 {
					baseURL := res.Connections[0].URI
					tracks, err := serverClient.GetPlaylistTracks(context.Background(), baseURL, "/playlists/"+loc.Query+"/items")
					if err == nil {
						probe = &plex.ConnectionResult{BaseURL: baseURL, Tracks: tracks}
					} else {
						lastErr = err
					}
				}
			}

			if lastErr != nil {
				fmt.Printf("  Failed to get tracks: %v\n", lastErr)
				continue
			}

			for _, t := range probe.Tracks {
				path := "No Path Found"
				if len(t.Media) > 0 && len(t.Media[0].Part) > 0 {
					path = t.Media[0].Part[0].File
				}
				fmt.Printf("  - %s - %s\n", t.Artist, t.Title)
				fmt.Printf("    Path: %s\n", path)
			}
		}
	}
	return nil
}

func init() {
	RootCmd.AddCommand(listCmd)
}
