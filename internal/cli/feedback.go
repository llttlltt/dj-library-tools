package cli

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/llttlltt/dj-library-tools/internal/utils"
)

type TerminalFeedback struct{}

func (f *TerminalFeedback) OnPreview(msg string) {
	fmt.Printf("[%s] %s\n", color.YellowString("Preview"), msg)
}

func (f *TerminalFeedback) OnSuccess(msg string) {
	fmt.Printf("%s %s\n", color.GreenString("✔"), msg)
}

func (f *TerminalFeedback) OnTable(headers []string, rows [][]string) {
	headerFmt := color.New(color.FgCyan, color.Bold, color.Underline).SprintFunc()
	t := utils.Table{
		Headers: headers,
		Rows:    rows,
		HeaderFormat: func(s string) string {
			return headerFmt(s)
		},
	}
	t.Render()
}
