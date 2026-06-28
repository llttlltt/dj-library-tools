package cli

import (
	"fmt"
	"github.com/fatih/color"
)

type TerminalFeedback struct{}

func (f *TerminalFeedback) OnPreview(msg string) {
	fmt.Printf("[%s] %s\n", color.YellowString("Preview"), msg)
}

func (f *TerminalFeedback) OnSuccess(msg string) {
	fmt.Printf("%s %s\n", color.GreenString("✔"), msg)
}
