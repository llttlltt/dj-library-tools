package cli

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/fatih/color"
	"github.com/llttlltt/dj-library-tools/internal/config"
	"github.com/llttlltt/dj-library-tools/internal/engine"
	"github.com/llttlltt/dj-library-tools/internal/plex"
	"github.com/llttlltt/dj-library-tools/internal/provider"
	"github.com/llttlltt/dj-library-tools/internal/utils"
	"github.com/llttlltt/dj-library-tools/pkg/rekordbox"
	"github.com/spf13/cobra"
	"os"
)

var listCmd = &cobra.Command{
	Use:     "list [resource] [query]",
	Aliases: []string{"ls"},
	Short:   "List items from a location (e.g. rb/tracks title:Oceans)",
	Args:    cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		var loc utils.Location
		if len(args) > 1 {
			loc = utils.ParseLocation(args[0], strings.Join(args[1:], " "))
		} else {
			loc = utils.ParseLocation(args[0], "")
		}

		switch loc.Provider {
		case "rb", "rekordbox":
			if loc.Resource == "" {
				return fmt.Errorf("resource must be specified (e.g. rb/tracks or rb/playlists)")
			}
			lib, _, err := loadXMLFunc()
			if err != nil {
				return err
			}
			eng := engine.NewEngine(engine.NewRekordboxLibrary(lib))
			prov := provider.NewRekordboxProvider(eng)
			return listProvider(prov, loc)
		case "plex":
			if loc.Resource == "" {
				return fmt.Errorf("resource must be specified (e.g. plex/playlists)")
			}
			return listPlex(loc)
		default:
			return fmt.Errorf("unknown provider: %s", loc.Provider)
		}
	},
}

func listProvider(p provider.Provider, loc utils.Location) error {
	if loc.Resource == "playlists" || loc.Resource == "folders" {
		results, err := p.GetPlaylists(loc.Query)
		if err != nil {
			return fmt.Errorf("ls failed: %w", err)
		}
		if jsonOutput {
			data, _ := json.MarshalIndent(results, "", "  ")
			fmt.Println(string(data))
			return nil
		}
		if verbose {
			fmt.Printf("Query %q matched %d %s\n", loc.Query, len(results), loc.Resource)
		}
		if len(results) == 0 {
			color.Yellow("No %s matched the query.", loc.Resource)
			return nil
		}
		for _, res := range results {
			name := res.Name
			if res.ParentFolder != "" {
				name = res.ParentFolder + "/" + name
			}
			fmt.Printf("%s: %s (%d entries)\n", strings.Title(loc.Resource[:len(loc.Resource)-1]), name, res.Entries)
		}
		return nil
	}

	tracks, err := p.GetTracks(loc.Query)
	if err != nil {
		return fmt.Errorf("ls failed: %w", err)
	}

	if jsonOutput {
		data, _ := json.MarshalIndent(tracks, "", "  ")
		fmt.Println(string(data))
		return nil
	}

	if verbose {
		fmt.Printf("Query %q matched %d tracks\n", loc.Query, len(tracks))
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

func listRekordbox(loc utils.Location) error {
	// listRekordbox is now absorbed by listProvider
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
