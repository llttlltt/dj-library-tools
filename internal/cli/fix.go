package cli

import (
	"fmt"
	"os"

	"github.com/llttlltt/dj-library-tools/internal/playlist"
	"github.com/spf13/cobra"
)

var (
	// fix flags
	extsFlag       []string
	m3u8Flag       bool
	removeOriginal bool
	forceOverwrite bool
	outputFileFlag string
	verboseFixFlag bool
	dryRunFixFlag  bool
)

var fixCmd = &cobra.Command{
	Use:   "fix",
	Short: "Fix library issues or resource metadata",
}

var fixPlaylistCmd = &cobra.Command{
	Use:   "playlist [file...]",
	Short: "Fix playlist extensions and/or enrich with M3U8 metadata",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		for _, inputPath := range args {
			opts := playlist.FixOptions{
				Exts:           extsFlag,
				M3U8:           m3u8Flag,
				RemoveOriginal: removeOriginal,
				Force:          forceOverwrite,
				OutputPath:     outputFileFlag,
				Verbose:        verboseFixFlag,
				DryRun:         dryRunFixFlag,
			}

			if len(args) > 1 && outputFileFlag != "" {
				fmt.Printf("Warning: --output ignored when processing multiple files. Using default names for %s\n", inputPath)
				opts.OutputPath = ""
			}

			result, err := playlist.FixPlaylist(inputPath, opts)
			if err != nil {
				fmt.Printf("Error processing %s: %v\n", inputPath, err)
				continue
			}

			if opts.DryRun {
				fmt.Printf("DRY RUN: Would process '%s' -> '%s'\n", inputPath, result.OutputPath)
			} else {
				fmt.Printf("Successfully processed '%s' -> '%s'\n", inputPath, result.OutputPath)
			}
			fmt.Printf("Total tracks found: %d\n", result.TotalTracks-len(result.SkippedTracks))
			if len(result.SkippedTracks) > 0 {
				fmt.Printf("Skipped tracks (not found): %d\n", len(result.SkippedTracks))
				if verboseFixFlag {
					for _, p := range result.SkippedTracks {
						fmt.Printf("  - %s\n", p)
					}
				} else {
					fmt.Println("  (Use -v to see full list of skipped tracks)")
				}
			}

			if !opts.DryRun && removeOriginal && inputPath != result.OutputPath {
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
			fmt.Println("---")
		}
		return nil
	},
}

func init() {
	fixPlaylistCmd.Flags().StringSliceVarP(&extsFlag, "ext", "e", []string{}, "Priority list of file extensions (comma-separated)")
	fixPlaylistCmd.Flags().BoolVar(&m3u8Flag, "m3u8", false, "Enrich playlist with M3U8 #EXTINF tags")
	fixPlaylistCmd.Flags().BoolVarP(&removeOriginal, "remove-original", "r", false, "Remove the original playlist file after processing")
	fixPlaylistCmd.Flags().BoolVarP(&forceOverwrite, "force", "f", false, "Force overwrite if output file exists")
	fixPlaylistCmd.Flags().StringVarP(&outputFileFlag, "output", "o", "", "Specific output path (optional)")
	fixPlaylistCmd.Flags().BoolVarP(&verboseFixFlag, "verbose", "v", false, "Enable verbose logging")
	fixPlaylistCmd.Flags().BoolVar(&dryRunFixFlag, "dry-run", false, "Show what would be done without modifying files")

	fixCmd.AddCommand(fixPlaylistCmd)
	RootCmd.AddCommand(fixCmd)
}
