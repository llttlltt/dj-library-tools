package cli

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	"github.com/fatih/color"
	"github.com/llttlltt/dj-library-tools/internal/engine"
	"github.com/llttlltt/dj-library-tools/internal/utils"
	"github.com/spf13/cobra"
)

var statCmd = &cobra.Command{
	Use:   "stat [resource] [query]",
	Short: "Show statistics for tracks matching the query",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		rbXML, _, err := loadXMLFunc()
		if err != nil {
			return err
		}

		eng := engine.NewEngine(rbXML)
		
		queryStr := ""
		if len(args) > 1 {
			queryStr = strings.Join(args[1:], " ")
		}
		loc := utils.ParseLocation(args[0], queryStr)

		if loc.Provider != "rb" || loc.Resource != "tracks" {
			return fmt.Errorf("stat currently only supports rb/tracks")
		}

		res, err := eng.Stat(loc.Query)
		if err != nil {
			return fmt.Errorf("stat failed: %w", err)
		}

		if jsonOutput {
			data, _ := json.MarshalIndent(res, "", "  ")
			fmt.Println(string(data))
			return nil
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
	RootCmd.AddCommand(statCmd)
}
