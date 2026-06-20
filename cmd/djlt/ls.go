package main

import (
	"fmt"
	"strings"

	"github.com/fatih/color"
	"github.com/llttlltt/dj-library-tools/internal/engine"
	"github.com/llttlltt/dj-library-tools/pkg/rekordbox"
	"github.com/spf13/cobra"
)

var lsCmd = &cobra.Command{
	Use:   "ls [query]",
	Short: "List tracks matching the query",
	Long: `List tracks from the Rekordbox XML library that match the specified query.
Example: djlt ls "artist:Four Tet bpm:120..128"`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if xmlPath == "" {
			return fmt.Errorf("XML path must be provided via --xml or -x")
		}

		lib, err := rekordbox.ReadRekordboxLibrary(xmlPath)
		if err != nil {
			return fmt.Errorf("failed to read XML: %w", err)
		}

		eng := engine.NewEngine(lib)
		queryStr := ""
		if len(args) > 0 {
			queryStr = strings.Join(args, " ")
		}

		tracks, err := eng.Ls(queryStr)
		if err != nil {
			return fmt.Errorf("ls failed: %w", err)
		}

		if len(tracks) == 0 {
			color.Yellow("No tracks matched the query.")
			return nil
		}

		headerFmt := color.New(color.FgCyan, color.Bold, color.Underline).SprintfFunc()
		dimFmt := color.New(color.FgHiBlack).SprintfFunc()
		artistFmt := color.New(color.FgHiMagenta).SprintfFunc()
		titleFmt := color.New(color.FgHiWhite).SprintfFunc()
		bpmFmt := color.New(color.FgHiGreen).SprintfFunc()
		keyFmt := color.New(color.FgHiYellow).SprintfFunc()

		fmt.Printf("%s %s %s\n", headerFmt("%6s", "BPM"), headerFmt("%4s", "Key"), headerFmt("Artist - Title"))
		for _, t := range tracks {
			bpm := 0.0
			if len(t.Tempo) > 0 {
				bpm = t.Tempo[0].Bpm
			}
			fmt.Printf("%s %s %s %s %s\n", 
				bpmFmt("%6.2f", bpm), 
				keyFmt("%4s", t.Tonality),
				artistFmt(t.Artist), 
				dimFmt("-"), 
				titleFmt(t.Name))
		}

		fmt.Printf("\n%s\n", color.HiGreenString("Matched %d tracks.", len(tracks)))
		return nil
	},
}

func init() {
	rootCmd.AddCommand(lsCmd)
}
