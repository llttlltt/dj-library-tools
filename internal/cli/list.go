package cli

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/fatih/color"
	"github.com/llttlltt/dj-library-tools/internal/config"
	"github.com/llttlltt/dj-library-tools/internal/provider"
	"github.com/llttlltt/dj-library-tools/internal/utils"
	"github.com/llttlltt/dj-library-tools/pkg/rekordbox"
	"github.com/spf13/cobra"
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

		if loc.Resource == "" {
			return fmt.Errorf("resource must be specified (e.g. %s/tracks or %s/playlists)", loc.Provider, loc.Provider)
		}

		cfg, _ := config.LoadAppConfig()
		var rbXML *rekordbox.RekordboxLibraryXML
		if loc.Provider == "rb" || loc.Provider == "rekordbox" {
			var err error
			rbXML, _, err = loadXMLFunc()
			if err != nil {
				return err
			}
		}

		prov, err := provider.NewProvider(loc.Provider, rbXML, cfg)
		if err != nil {
			return err
		}

		return listProvider(prov, loc)
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
			fmt.Printf("%s: %s (%d entries)\n", stringsTitle(loc.Resource[:len(loc.Resource)-1]), name, res.Entries)
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

func init() {
	RootCmd.AddCommand(listCmd)
}
