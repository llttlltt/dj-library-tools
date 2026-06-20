package main

import (
	"fmt"
	"strings"

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
			fmt.Println("No tracks matched the query.")
			return nil
		}

		for _, t := range tracks {
			bpm := 0.0
			if len(t.Tempo) > 0 {
				bpm = t.Tempo[0].Bpm
			}
			fmt.Printf("[%6.2f] %s - %s\n", bpm, t.Artist, t.Name)
		}

		fmt.Printf("\nMatched %d tracks.\n", len(tracks))
		return nil
	},
}

func init() {
	rootCmd.AddCommand(lsCmd)
}
