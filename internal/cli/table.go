package cli

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/fatih/color"
	"github.com/llttlltt/dj-library-tools/internal/models"
)

type Table struct {
	Headers []string
	Rows    [][]string
}

func (t *Table) Render() {
	if len(t.Rows) == 0 {
		return
	}

	headerFmt := color.New(color.FgCyan, color.Bold, color.Underline).SprintfFunc()

	// Calculate column widths
	widths := make([]int, len(t.Headers))
	for i, h := range t.Headers {
		widths[i] = len(h)
	}
	for _, row := range t.Rows {
		for i, val := range row {
			if len(val) > widths[i] {
				widths[i] = len(val)
			}
		}
	}

	// Print Headers
	for i, h := range t.Headers {
		fmt.Print(headerFmt("%-*s", widths[i], h))
		if i < len(t.Headers)-1 {
			fmt.Print(" ")
		}
	}
	fmt.Println()

	// Print Rows
	for _, row := range t.Rows {
		for i, val := range row {
			// Apply specific colors based on header name
			rendered := val
			header := strings.ToLower(t.Headers[i])
			switch {
			case header == "bpm":
				rendered = color.HiGreenString("%*s", widths[i], val)
			case header == "key":
				rendered = color.HiYellowString("%*s", widths[i], val)
			case header == "artist":
				rendered = color.HiMagentaString("%-*s", widths[i], val)
			case header == "title" || header == "name":
				rendered = color.HiWhiteString("%-*s", widths[i], val)
			case header == "entries" || header == "count":
				rendered = color.CyanString("%*s", widths[i], val)
			default:
				rendered = fmt.Sprintf("%-*s", widths[i], val)
			}

			fmt.Print(rendered)
			if i < len(row)-1 {
				fmt.Print(" ")
			}
		}
		fmt.Println()
	}
}

func renderTrackTable(tracks []models.Track) {
	table := Table{
		Headers: []string{"BPM", "Key", "Artist", "Title"},
	}

	for _, t := range tracks {
		table.Rows = append(table.Rows, []string{
			fmt.Sprintf("%6.2f", t.BPM),
			fmt.Sprintf("%4s", t.Key),
			t.Artist,
			t.Title,
		})
	}

	table.Render()
	fmt.Printf("\n%s\n", color.HiGreenString("Matched %d tracks.", len(tracks)))
}

func renderNodeTable(results []models.Node, resourceType string) {
	table := Table{
		Headers: []string{"Entries", stringsTitle(resourceType)},
	}

	for _, res := range results {
		name := res.Name
		if res.ParentFolder != "" {
			name = res.ParentFolder + "/" + name
		}
		table.Rows = append(table.Rows, []string{
			strconv.Itoa(res.Entries),
			name,
		})
	}

	table.Render()
	fmt.Printf("\n%s\n", color.HiGreenString("Matched %d %s.", len(results), resourceType))
}
