package main

import (
	"fmt"
	"sort"
	"strings"

	"github.com/fatih/color"
	"github.com/llttlltt/dj-library-tools/internal/engine"
	"github.com/llttlltt/dj-library-tools/pkg/rekordbox"
	"github.com/spf13/cobra"
)

var statCmd = &cobra.Command{
	Use:   "stat [query]",
	Short: "Show statistics for tracks matching the query",
	RunE: func(cmd *cobra.Command, args []string) error {
		if xmlPath == "" {
			return fmt.Errorf("XML path must be provided via --xml or -x")
		}

		lib, err := rekordbox.ReadRekordboxLibrary(xmlPath)
		if err != nil {
			return fmt.Errorf("failed to read XML: %w", err)
		}

		eng := engine.NewEngine(lib)
		queryStr := ""
		if len(args) > 0 {
			queryStr = strings.Join(args, " ")
		}

		res, err := eng.Stat(queryStr)
		if err != nil {
			return fmt.Errorf("stat failed: %w", err)
		}

		titleFmt := color.New(color.FgCyan, color.Bold, color.Underline).SprintFunc()
		labelFmt := color.New(color.FgHiWhite).SprintFunc()
		valFmt := color.New(color.FgHiGreen).SprintFunc()

		fmt.Printf("%s\n", titleFmt("Library Summary"))
		fmt.Printf("%-15s %s\n", labelFmt("Total Tracks:"), valFmt(fmt.Sprintf("%d", res.Count)))
		if res.Count > 0 {
			fmt.Printf("%-15s %s\n", labelFmt("Average BPM:"), valFmt(fmt.Sprintf("%.2f", res.AvgBPM)))
		}

		printTop(res.Genres, "Top Genres", 5)
		printTop(res.Artists, "Top Artists", 5)
		printTop(res.Keys, "Top Keys", 5)

		return nil
	},
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

	titleFmt := color.New(color.FgCyan, color.Bold, color.Underline).SprintFunc()
	labelFmt := color.New(color.FgHiWhite).SprintFunc()
	valFmt := color.New(color.FgHiMagenta).SprintFunc()

	fmt.Printf("\n%s\n", titleFmt(title))
	for i, kv := range ss {
		if i >= limit {
			break
		}
		fmt.Printf("%-20s %s\n", labelFmt(kv.Key), valFmt(fmt.Sprintf("%d", kv.Value)))
	}
}

func init() {
	rootCmd.AddCommand(statCmd)
}
