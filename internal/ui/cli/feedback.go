package cli

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/llttlltt/dj-library-tools/internal/ui/cli/render"
)

type TerminalFeedback struct{}

func (f *TerminalFeedback) OnPreview(msg string) {
	fmt.Printf("[%s] %s\n", color.YellowString("Preview"), msg)
}

func (f *TerminalFeedback) OnSuccess(msg string) {
	fmt.Printf("%s %s\n", color.GreenString("✔"), msg)
}

func (f *TerminalFeedback) OnWarning(msg string) {
	fmt.Fprintf(os.Stderr, "⚠️  %s\n", color.YellowString(msg))
}

func (f *TerminalFeedback) OnStatus(msg string) {
	fmt.Println(msg)
}

func (f *TerminalFeedback) OnProgress(done, total int) {
	fmt.Printf("\rProcessing: [%d/%d]", done, total)
	if done == total {
		fmt.Println()
	}
}
func (f *TerminalFeedback) OnTable(headers []string, rows [][]string) {
	headerFmt := color.New(color.FgCyan, color.Bold, color.Underline).SprintFunc()
	t := render.Table{
		Headers: headers,
		Rows:    rows,
		HeaderFormat: func(s string) string {
			return headerFmt(s)
		},
	}
	t.Render()
}
