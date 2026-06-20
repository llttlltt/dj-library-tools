package main

import (
	"fmt"
	"os"

	"github.com/llttlltt/dj-library-tools/internal/playlist"
	"github.com/spf13/cobra"
)

var (
	extFlag           string
	m3u8Flag          bool
	removeOriginal    bool
	forceOverwrite    bool
	outputFileFlag    string
	verboseFlag       bool
)

var playlistCmd = &cobra.Command{
	Use:   "playlist",
	Short: "Manage playlists",
}

var fixCmd = &cobra.Command{
	Use:   "fix [file]",
	Short: "Fix playlist extensions and/or enrich with M3U8 metadata",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		inputPath := args[0]
		opts := playlist.FixOptions{
			Ext:            extFlag,
			M3U8:           m3u8Flag,
			RemoveOriginal: removeOriginal,
			Force:          forceOverwrite,
			OutputPath:     outputFileFlag,
			Verbose:        verboseFlag,
		}
		result, err := playlist.FixPlaylist(inputPath, opts)
		if err != nil {
			return err
		}

		fmt.Printf("Successfully processed '%s' -> '%s'\n", inputPath, result.OutputPath)
		fmt.Printf("Total tracks: %d\n", result.TotalTracks)
		if len(result.MissingTracks) > 0 {
			fmt.Printf("⚠️ Missing tracks: %d\n", len(result.MissingTracks))
			if verboseFlag {
				for _, path := range result.MissingTracks {
					fmt.Printf("  - %s\n", path)
				}
			} else {
				fmt.Println("  (Use -v to see full list of missing tracks)")
			}
		}

		// Prompt for removal if requested and output is different
		if removeOriginal && inputPath != result.OutputPath {
			fmt.Printf("\nRemove original file '%s'? (y/N): ", inputPath)
			var response string
			fmt.Scanln(&response)
			if response == "y" || response == "Y" {
				if err := os.Remove(inputPath); err != nil {
					return fmt.Errorf("failed to remove original file: %w", err)
				}
				fmt.Println("Original file removed.")
			} else {
				fmt.Println("Original file retained.")
			}
		}

		return nil
	},
}

func init() {
	fixCmd.Flags().StringVarP(&extFlag, "ext", "e", "", "New file extension (e.g., .mp3)")
	fixCmd.Flags().BoolVar(&m3u8Flag, "m3u8", false, "Enrich playlist with M3U8 #EXTINF tags")
	fixCmd.Flags().BoolVarP(&removeOriginal, "remove-original", "r", false, "Remove the original playlist file after processing")
	fixCmd.Flags().BoolVarP(&forceOverwrite, "force", "f", false, "Force overwrite if output file exists")
	fixCmd.Flags().StringVarP(&outputFileFlag, "output", "o", "", "Specific path for the output file (optional)")
	fixCmd.Flags().BoolVarP(&verboseFlag, "verbose", "v", false, "Enable verbose logging")

	playlistCmd.AddCommand(fixCmd)
	rootCmd.AddCommand(playlistCmd)
}
