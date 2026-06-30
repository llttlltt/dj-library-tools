package render

import (
	"fmt"
	"strings"
)

// Table represents a simple console table with optional truncation.
type Table struct {
	Headers      []string
	Rows         [][]string
	HeaderFormat func(string) string
}

// Render writes the table to standard out.
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

	// Smart truncation for standard music metadata fields
	for i, h := range t.Headers {
		lowerH := strings.ToLower(h)
		if lowerH == "artist" || lowerH == "title" || lowerH == "album" {
			if colWidths[i] > 35 {
				colWidths[i] = 35
			}
		}
	}

	// Render headers
	for i, h := range t.Headers {
		text := h
		if t.HeaderFormat != nil {
			text = t.HeaderFormat(text)
		}
		fmt.Print(text)
		
		padding := colWidths[i] - len(h)
		if padding < 0 { padding = 0 }
		fmt.Print(strings.Repeat(" ", padding + 1))
	}
	fmt.Println()

	// Render rows
	for _, row := range t.Rows {
		for i, val := range row {
			text := val
			if len(text) > colWidths[i] {
				text = text[:colWidths[i]-3] + "..."
			}
			fmt.Printf("%-*s ", colWidths[i], text)
		}
		fmt.Println()
	}
}
