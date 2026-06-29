package cli

import (
	"fmt"
	"strings"

	"github.com/fatih/color"
	"github.com/llttlltt/dj-library-tools/internal/models"
	"github.com/llttlltt/dj-library-tools/internal/provider"
	"github.com/llttlltt/dj-library-tools/internal/query"
	"github.com/llttlltt/dj-library-tools/internal/utils"
)

type Table struct {
	Headers []string
	Rows    [][]string
}

func (f *Table) Render() {
	if len(f.Rows) == 0 {
		return
	}

	headerFmt := color.New(color.FgCyan, color.Bold, color.Underline).SprintFunc()
	t := utils.Table{
		Headers: f.Headers,
		Rows:    f.Rows,
		HeaderFormat: func(s string) string {
			return headerFmt(s)
		},
	}
	t.Render()
}

func renderTrackTable(prov provider.Provider, tracks []models.Track, columns []string) {
	if len(columns) == 0 {
		columns = prov.System().TableHeaders()
	}

	t := &Table{Headers: columns}
	for _, tr := range tracks {
		row := make([]string, len(columns))
		for i, col := range columns {
			field := strings.ToLower(col)
			val := tr.Value(field)
			
			// If not a standard field, try path resolution
			if val == "" && (strings.ContainsAny(field, "./-")) {
				if pVal, ok := query.ResolvePath(tr, field); ok {
					val = pVal
				}
			}
			row[i] = val
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
