package cli

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var (
	listSort string
)

var listCmd = &cobra.Command{
	Use:     "list [resource] [query]",
	Aliases: []string{"ls"},
	Short:   "List items from a location (e.g. rb/tracks title:Oceans)",
	Args:    cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		queryOverride := ""
		if len(args) > 1 {
			queryOverride = strings.Join(args[1:], " ")
		}

		sel, err := ResolveSelection(args[0], queryOverride)
		if err != nil {
			return err
		}

		return listProvider(sel)
	},
}

func listProvider(sel *Selection) error {
	if sel.Location.Resource == "playlists" || sel.Location.Resource == "folders" {
		if jsonOutput {
			data, _ := json.MarshalIndent(sel.Nodes, "", "  ")
			fmt.Println(string(data))
			return nil
		}
		if len(sel.Nodes) == 0 {
			color.Yellow("No %s matched the query.", sel.Location.Resource)
			return nil
		}

		sortNodes(sel.Nodes, listSort)
		renderNodeTable(sel.Nodes, sel.Location.Resource[:len(sel.Location.Resource)-1])
		return nil
	}

	if jsonOutput {
		data, _ := json.MarshalIndent(sel.Tracks, "", "  ")
		fmt.Println(string(data))
		return nil
	}

	if len(sel.Tracks) == 0 {
		color.Yellow("No tracks matched the query.")
		return nil
	}

	sortTracks(sel.Tracks, listSort)
	renderTrackTable(sel.Tracks)
	return nil
}

func init() {
	listCmd.Flags().StringVar(&listSort, "sort", "", "Sort results by field (e.g. bpm, artist, title)")
	RootCmd.AddCommand(listCmd)
}
