package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func EnforceExtension(path, extension string) string {
	base := strings.TrimSuffix(path, filepath.Ext(path))
	if extension != "" && !strings.HasPrefix(extension, ".") {
		extension = "." + extension
	}
	return base + extension
}

func CheckFileOverwrite(path string, forceOverwrite bool) error {
	if _, err := os.Stat(path); err == nil {
		if !forceOverwrite {
			return fmt.Errorf("file '%s' already exists. Use the --force (-f) flag to overwrite", path)
		}
	} else if !os.IsNotExist(err) {
		return fmt.Errorf("failed to check file status '%s': %w", path, err)
	}
	return nil
}

func ExpandPath(path string) string {
	if strings.HasPrefix(path, "~/") {
		home, err := os.UserHomeDir()
		if err != nil {
			return path
		}
		return filepath.Join(home, path[2:])
	}
	return path
}

// Table represents a simple console table.
type Table struct {
	Headers []string
	Rows    [][]string
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

	// Limit column widths for certain columns to keep table compact
	for i, h := range t.Headers {
		if h == "ARTIST" || h == "TITLE" {
			if colWidths[i] > 30 {
				colWidths[i] = 30
			}
		}
	}

	for i, h := range t.Headers {
		fmt.Printf("%-*s ", colWidths[i], strings.ToUpper(h))
	}
	fmt.Println()

	for _, row := range t.Rows {
		for i, val := range row {
			text := val
			if (t.Headers[i] == "ARTIST" || t.Headers[i] == "TITLE") && len(text) > colWidths[i] {
				text = text[:colWidths[i]-3] + "..."
			}
			fmt.Printf("%-*s ", colWidths[i], text)
		}
		fmt.Println()
	}
}
