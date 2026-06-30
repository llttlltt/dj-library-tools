package cli

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	"github.com/fatih/color"
	"github.com/llttlltt/dj-library-tools/internal/core/models"
	"github.com/llttlltt/dj-library-tools/internal/core/query"
	"github.com/llttlltt/dj-library-tools/internal/providers"
	"github.com/spf13/cobra"
)

func newListCmd() *cobra.Command {
	var listSort string
	var listStats bool
	var jsonOutput bool
	var filterMissing bool
	var filterExists bool
	var columns []string

	cmd := &cobra.Command{
		Use:     "ls [resource] [query]",
		Short:   "List items from a location (e.g. rb/tracks title:Oceans)",
		Args:    cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			queryOverride := ""
			if len(args) > 1 {
				queryOverride = strings.Join(args[1:], " ")
			}

			orch := getOrchestrator()
			runOpts := getRunOptions()
			res, err := orch.List(args[0], queryOverride, runOpts)
			if err != nil {
				return HandleError(err)
			}

			if listStats {
				return listProviderStats(res.Tracks, res.Groups, jsonOutput)
			}

			return listProvider(res.Tracks, res.Groups, res.Provider, listSort, jsonOutput, columns)
		},
	}
	cmd.Flags().StringVar(&listSort, "sort", "", "Sort results by any available field (e.g. artist, title, bpm, etc.)")
	cmd.Flags().BoolVar(&listStats, "stats", false, "Show summary statistics for the selection")
	cmd.Flags().BoolVar(&jsonOutput, "json", false, "Output results in JSON format")
	cmd.Flags().BoolVar(&filterMissing, "missing", false, "Filter for tracks where the physical file is missing")
	cmd.Flags().BoolVar(&filterExists, "exists", false, "Filter for tracks where the physical file exists")
	cmd.Flags().StringSliceVar(&columns, "columns", []string{}, "Comma-separated list of columns to display")
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

func listProviderStats(tracks []models.Track, groups []models.ResourceGroup, jsonOutput bool) error {
	// Calculate stats from the neutral tracks
	res := StatResult{
		Count:   len(tracks),
		Genres:  make(map[string]int),
		Labels:  make(map[string]int),
		Keys:    make(map[string]int),
		Artists: make(map[string]int),
	}

	if len(tracks) == 0 {
		color.Yellow("No tracks matched the query.")
		return nil
	}

	totalBPM := 0.0
	for _, t := range tracks {
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
	res.AvgBPM = totalBPM / float64(len(tracks))

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

func listProvider(tracks []models.Track, groups []models.ResourceGroup, prov provider.Provider, listSort string, jsonOutput bool, columns []string) error {
	if len(groups) > 0 {
		// Group Logic...
		if listSort != "" {
			if _, ok := models.GroupFields[listSort]; !ok {
				return fmt.Errorf("invalid sort field %q; valid fields are: %v", listSort, strings.Join(query.AllowedGroupFields, ", "))
			}
			prov.Groups().Sort(getExecContext(), groups, listSort)
		}
		renderGroupTable(groups, "Group")
		return nil
	}

	// Track Logic...
	if len(tracks) == 0 {
		color.Yellow("No tracks matched the query.")
		return nil
	}

	if listSort != "" {
		if _, ok := models.TrackFields[listSort]; !ok {
			return fmt.Errorf("invalid sort field %q; valid fields are: %v", listSort, strings.Join(query.AllowedTrackFields, ", "))
		}
		prov.Tracks().Sort(getExecContext(), tracks, listSort)
	}
	renderTrackTable(prov, tracks, columns)
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
