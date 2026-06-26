package cli

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/fatih/color"
	"github.com/llttlltt/dj-library-tools/internal/provider"
	"github.com/llttlltt/dj-library-tools/pkg/rekordbox"
	"golang.org/x/term"
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

	// Get terminal width
	termWidth, _, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil || termWidth <= 0 {
		termWidth = 80 // fallback
	}

	// Calculate initial column widths
	numCols := len(t.Headers)
	widths := make([]int, numCols)
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

	// Dynamic scaling logic
	// We want to preserve BPM (6) and Key (4) widths as they are fixed
	// We allocate the rest of the terminal width to Artist and Title/Name
	totalFixed := 0
	flexibleCols := 0
	for i, h := range t.Headers {
		hLower := strings.ToLower(h)
		if hLower == "bpm" || hLower == "key" || hLower == "entries" {
			totalFixed += widths[i] + 1 // +1 for gutter
		} else {
			flexibleCols++
		}
	}

	availableWidth := termWidth - totalFixed - (numCols - 1)
	if flexibleCols > 0 {
		maxFlexWidth := availableWidth / flexibleCols
		if maxFlexWidth < 10 {
			maxFlexWidth = 10 // Minimum sensible width
		}

		for i, h := range t.Headers {
			hLower := strings.ToLower(h)
			if hLower != "bpm" && hLower != "key" && hLower != "entries" {
				if widths[i] > maxFlexWidth {
					widths[i] = maxFlexWidth
				}
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
			displayVal := val
			if len(displayVal) > widths[i] {
				displayVal = displayVal[:widths[i]-3] + "..."
			}

			// Apply specific colors based on header name
			rendered := displayVal
			header := strings.ToLower(t.Headers[i])
			switch {
			case header == "bpm":
				rendered = color.HiGreenString("%*s", widths[i], displayVal)
			case header == "key":
				rendered = color.HiYellowString("%*s", widths[i], displayVal)
			case header == "artist":
				rendered = color.HiMagentaString("%-*s", widths[i], displayVal)
			case header == "title" || header == "name":
				rendered = color.HiWhiteString("%-*s", widths[i], displayVal)
			case header == "entries" || header == "count":
				rendered = color.CyanString("%*s", widths[i], displayVal)
			default:
				rendered = fmt.Sprintf("%-*s", widths[i], displayVal)
			}

			fmt.Print(rendered)
			if i < len(row)-1 {
				fmt.Print(" ")
			}
		}
		fmt.Println()
	}
}

func renderTrackTable(tracks []rekordbox.Track) {
	table := Table{
		Headers: []string{"BPM", "Key", "Artist", "Title"},
	}

	for _, t := range tracks {
		bpm := "0.00"
		if len(t.Tempo) > 0 {
			bpm = t.Tempo[0].Bpm
		}
		table.Rows = append(table.Rows, []string{
			fmt.Sprintf("%6s", bpm),
			fmt.Sprintf("%4s", t.Tonality),
			t.Artist,
			t.Name,
		})
	}

	table.Render()
	fmt.Printf("\n%s\n", color.HiGreenString("Matched %d tracks.", len(tracks)))
}

func renderNodeTable(results []provider.NodeResult, resourceType string) {
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
