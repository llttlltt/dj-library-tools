package cli

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	"github.com/fatih/color"
	"github.com/llttlltt/dj-library-tools/internal/services/orchestrator"
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
		Use:   "ls [resource] [query]",
		Short: "List items from a location (e.g. rb/tracks title:Oceans)",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			queryOverride := ""
			if len(args) > 1 {
				queryOverride = strings.Join(args[1:], " ")
			}

			orch := getOrchestrator()
			runOpts := getRunOptions()
			runOpts.FilterMissing = filterMissing
			runOpts.FilterExists = filterExists

			if listStats {
				res, err := orch.Stats(cmd.Context(), args[0], queryOverride, runOpts)
				if err != nil {
					return HandleError(err)
				}
				return renderStats(res, jsonOutput)
			}

			res, err := orch.List(cmd.Context(), args[0], queryOverride, runOpts, listSort)
			if err != nil {
				return HandleError(err)
			}

			return listProvider(cmd.Context(), res, jsonOutput, columns)
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

func renderStats(res *orchestrator.StatResult, jsonOutput bool) error {
	if jsonOutput {
		data, _ := json.MarshalIndent(res, "", "  ")
		fmt.Println(string(data))
		return nil
	}

	if res.Count == 0 {
		color.Yellow("No tracks matched the query.")
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

func listProvider(ctx context.Context, res *orchestrator.ListResult, jsonOutput bool, columns []string) error {
	if len(res.Groups) > 0 {
		label := res.Resource
		if len(label) > 0 && label[len(label)-1] == 's' {
			label = label[:len(label)-1]
		}
		renderGroupTable(res.Groups, label)
		return nil
	}

	// Track Logic...
	if len(res.Tracks) == 0 {
		color.Yellow("No tracks matched the query.")
		return nil
	}

	renderTrackTable(res.Tracks, columns, res.DefaultColumns)
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
