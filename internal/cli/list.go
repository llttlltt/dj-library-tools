package cli

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

func newListCmd() *cobra.Command {
	var listSort string
	var listStats bool

	cmd := &cobra.Command{
		Use:     "ls [resource] [query]",
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

			if listStats {
				return listProviderStats(sel)
			}

			return listProvider(sel, listSort)
		},
	}
	cmd.Flags().StringVar(&listSort, "sort", "", "Sort results by field (e.g. bpm, artist, title)")
	cmd.Flags().BoolVar(&listStats, "stats", false, "Show summary statistics for the selection")
	return cmd
}

type StatResult struct {
	Count      int
	AvgBPM     float64
	Genres     map[string]int
	Labels     map[string]int
	Keys       map[string]int
	Artists    map[string]int
	TotalTempo float64
}

func listProviderStats(sel *Selection) error {
	if sel.Location.Resource != "tracks" {
		return fmt.Errorf("stats only available for track resources")
	}

	// Calculate stats from the neutral tracks
	res := StatResult{
		Count:   len(sel.Tracks),
		Genres:  make(map[string]int),
		Labels:  make(map[string]int),
		Keys:    make(map[string]int),
		Artists: make(map[string]int),
	}

	if len(sel.Tracks) == 0 {
		color.Yellow("No tracks matched the query.")
		return nil
	}

	totalBPM := 0.0
	for _, t := range sel.Tracks {
		if t.Genre != "" {
			res.Genres[t.Genre]++
		}
		if t.Label != "" {
			res.Labels[t.Label]++
		}
		if t.Key != "" {
			res.Keys[t.Key]++
		}
		if t.Artist != "" {
			res.Artists[t.Artist]++
		}
		totalBPM += t.BPM
	}
	res.AvgBPM = totalBPM / float64(len(sel.Tracks))

	if jsonOutput {
		data, _ := json.MarshalIndent(res, "", "  ")
		fmt.Println(string(data))
		return nil
	}

	titleFmt := color.New(color.FgCyan, color.Bold, color.Underline).SprintFunc()
	labelFmt := color.New(color.FgHiWhite).SprintFunc()
	valFmt := color.New(color.FgHiGreen).SprintFunc()

	fmt.Printf("%s\n", titleFmt("Selection Summary"))
	fmt.Printf("%-15s %s\n", labelFmt("Total Tracks:"), valFmt(fmt.Sprintf("%d", res.Count)))
	fmt.Printf("%-15s %s\n", labelFmt("Average BPM:"), valFmt(fmt.Sprintf("%.2f", res.AvgBPM)))

	printTop(res.Genres, "Top Genres", 5)
	printTop(res.Artists, "Top Artists", 5)
	printTop(res.Keys, "Top Keys", 5)

	return nil
}

func listProvider(sel *Selection, listSort string) error {
	if sel.Location.Resource == "playlists" || sel.Location.Resource == "folders" {
		if jsonOutput {
			data, _ := json.MarshalIndent(sel.Groups, "", "  ")
			fmt.Println(string(data))
			return nil
		}
		if len(sel.Groups) == 0 {
			color.Yellow("No %s matched the query.", sel.Location.Resource)
			return nil
		}

		sortGroups(sel, sel.Groups, listSort)
		renderGroupTable(sel.Groups, sel.Location.Resource[:len(sel.Location.Resource)-1])
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

	sortTracks(sel, sel.Tracks, listSort) //nolint:staticcheck
	renderTrackTable(sel.Tracks)
	return nil
}

func printTop(m map[string]int, title string, limit int) {
	if len(m) == 0 {
		return
	}

	type kv struct {
		Key   string
		Value int
	}

	var ss []kv
	for k, v := range m {
		ss = append(ss, kv{k, v})
	}

	sort.Slice(ss, func(i, j int) bool {
		return ss[i].Value > ss[j].Value
	})

	headerFmt := color.New(color.FgCyan, color.Bold, color.Underline).SprintFunc()
	labelFmt := color.New(color.FgHiWhite).SprintFunc()
	valFmt := color.New(color.FgHiMagenta).SprintFunc()

	fmt.Printf("\n%s\n", headerFmt(title))
	for i, kv := range ss {
		if i >= limit {
			break
		}
		fmt.Printf("%-20s %s\n", labelFmt(kv.Key), valFmt(fmt.Sprintf("%d", kv.Value)))
	}
}


