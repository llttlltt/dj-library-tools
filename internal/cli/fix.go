package cli

import (
	"fmt"
	"os"

	"github.com/llttlltt/dj-library-tools/internal/playlist"
	"github.com/spf13/cobra"
)

func newFixCmd() *cobra.Command {
	var extsFlag []string
	var m3u8Flag, removeOriginal, forceOverwrite bool
	var outputFileFlag string

	fix := &cobra.Command{
		Use:   "fix",
		Short: "Fix library issues or resource metadata",
	}

	fixPlaylist := &cobra.Command{
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
				Verbose:        verbose,
				DryRun:         dryRun,
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

				if dryRun {
					fmt.Printf("DRY RUN: Would process '%s' -> '%s'\n", inputPath, result.OutputPath)
				} else {
					fmt.Printf("Successfully processed '%s' -> '%s'\n", inputPath, result.OutputPath)
				}
				fmt.Printf("Total tracks found: %d\n", result.TotalTracks-len(result.SkippedTracks))
				if len(result.SkippedTracks) > 0 {
					fmt.Printf("Skipped tracks (not found): %d\n", len(result.SkippedTracks))
					if verbose {
						for _, p := range result.SkippedTracks {
							fmt.Printf("  - %s\n", p)
						}
					} else {
						fmt.Println("  (Use -v to see full list of skipped tracks)")
					}
				}

				if !dryRun && removeOriginal && inputPath != result.OutputPath {
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
	fixPlaylist.Flags().StringSliceVarP(&extsFlag, "ext", "e", []string{}, "Priority list of file extensions (comma-separated)")
	fixPlaylist.Flags().BoolVar(&m3u8Flag, "m3u8", false, "Enrich playlist with M3U8 #EXTINF tags")
	fixPlaylist.Flags().BoolVarP(&removeOriginal, "remove-original", "r", false, "Remove the original playlist file after processing")
	fixPlaylist.Flags().BoolVarP(&forceOverwrite, "force", "f", false, "Force overwrite if output file exists")
	fixPlaylist.Flags().StringVarP(&outputFileFlag, "output", "o", "", "Specific output path (optional)")
	fix.AddCommand(fixPlaylist)
	return fix
}
