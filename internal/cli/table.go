package cli

import (
	"fmt"

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

	colWidths := make([]int, len(t.Headers))
	for i, h := range t.Headers {
		colWidths[i] = len(h)
	}

	for _, row := range t.Rows {
		for i, val := range row {
			if len(val) > colWidths[i] {
				colWidths[i] = len(val)
			}
		}
	}

	headerFmt := color.New(color.FgCyan, color.Bold, color.Underline).SprintFunc()
	for i, h := range t.Headers {
		fmt.Printf("%-*s ", colWidths[i], headerFmt(h))
	}
	fmt.Println()

	for _, row := range t.Rows {
		for i, val := range row {
			fmt.Printf("%-*s ", colWidths[i], val)
		}
		fmt.Println()
	}
}

func renderTrackTable(providerName string, tracks []models.Track) {
	// Standard headers for all providers
	headers := []string{"BPM", "Key", "Artist", "Title", "Rating"}
	if providerName == "m3u" || providerName == "m3u8" {
		headers = []string{"Duration", "Display Name", "Location"}
	}

	t := &Table{Headers: headers}
	for _, tr := range tracks {
		var row []string
		if providerName == "m3u" || providerName == "m3u8" {
			row = []string{
				fmt.Sprintf("%d", tr.Duration),
				tr.Display,
				tr.Location,
			}
		} else {
			row = []string{
				fmt.Sprintf("%.2f", tr.BPM),
				tr.Key,
				tr.Artist,
				tr.Title,
				fmt.Sprintf("%d", tr.Rating),
			}
		}
		t.Rows = append(t.Rows, row)
	}
	t.Render()
	fmt.Printf("\nMatched %d tracks.\n", len(tracks))
}

func renderGroupTable(groups []models.ResourceGroup, label string) {
	t := &Table{
		Headers: []string{"Items", label},
	}
	for _, n := range groups {
		path := n.Name
		if n.ParentFolder != "" {
			path = n.ParentFolder + "/" + n.Name
		}
		t.Rows = append(t.Rows, []string{
			fmt.Sprintf("%d", n.Items),
			path,
		})
	}
	t.Render()
	fmt.Printf("\nMatched %d %s.\n", len(groups), label)
}
