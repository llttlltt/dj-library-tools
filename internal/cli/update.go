package cli

import (
	"fmt"

	"github.com/llttlltt/dj-library-tools/internal/library"
	"github.com/llttlltt/dj-library-tools/internal/utils"
	"github.com/llttlltt/dj-library-tools/internal/rekordbox"
	"github.com/spf13/cobra"
)

func newUpdateCmd() *cobra.Command {
	var updateFrom, updateTo, updateOutput string
	var updateMerge, updateForce bool

	cmd := &cobra.Command{
		Use:   "update [resource] [query] --from [source-file]",
		Short: "Update track metadata or merge markers between libraries",
		Long: `Update metadata for tracks in the library using another Rekordbox XML as a source.
Currently supports updating/merging Tempo markers (Beatgrids).

Example:
  djlt update rb/tracks --from other_library.xml`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if updateFrom == "" {
				return fmt.Errorf("--from source-file is required")
			}

			destLibrary, destPath, err := loadXMLFunc()
			if err != nil {
				return err
			}

			if updateTo != "" {
				destLibrary, err = rekordbox.ReadRekordboxLibrary(updateTo)
				if err != nil {
					return fmt.Errorf("error reading destination library %q: %w", updateTo, err)
				}
				destPath = updateTo
			}

			outputPath := updateOutput
			if outputPath == "" {
				outputPath = destPath
			}

			if err := utils.CheckFileOverwrite(outputPath, updateForce); err != nil {
				return fmt.Errorf("output file validation failed: %w", err)
			}

			sourceLibrary, err := rekordbox.ReadRekordboxLibrary(updateFrom)
			if err != nil {
				return fmt.Errorf("error reading source library %q: %w", updateFrom, err)
			}

			mergedLibrary := mergeLibraryData(sourceLibrary, destLibrary)

			if dryRun {
				fmt.Printf("[Dry Run] Would write updated library to %q\n", outputPath)
				return nil
			}

			destLibWrapper := library.NewRekordboxLibrary(mergedLibrary)
			if err := destLibWrapper.Save(outputPath); err != nil {
				return fmt.Errorf("error writing merged library to %q: %w", outputPath, err)
			}
			fmt.Printf("Successfully updated library: %s\n", outputPath)
			return nil
		},
	}
	cmd.Flags().StringVarP(&updateFrom, "from", "f", "", "Source Rekordbox XML to read metadata from")
	cmd.Flags().StringVarP(&updateTo, "to", "t", "", "Destination Rekordbox XML to update (defaults to primary library)")
	cmd.Flags().StringVarP(&updateOutput, "output", "o", "", "Output path for the updated Rekordbox XML")
	cmd.Flags().BoolVar(&updateMerge, "merge", false, "Merge metadata instead of overwriting")
	cmd.Flags().BoolVar(&updateForce, "force", false, "Overwrite output file if it already exists")
	return cmd
}

// ── matching & merging helpers ─────────────────────────────────────────────

type TrackMatch struct {
	SourceTrack      rekordbox.Track
	DestinationTrack rekordbox.Track
	FoundMatch       bool
}

type TrackMatches struct {
	tracks         []*TrackMatch
	matchedCount   int
	unmatchedCount int
}

func matchTrackInLibrary(sourceTrack rekordbox.Track, destinationLibrary *rekordbox.RekordboxLibraryXML) (*rekordbox.Track, error) {
	for _, destinationTrack := range destinationLibrary.Collection.TRACK {
		if tracksMatchMetadata(sourceTrack, destinationTrack) {
			return &destinationTrack, nil
		}
	}
	return nil, fmt.Errorf("no match found for %q", sourceTrack.Name)
}

func tracksMatchMetadata(t1, t2 rekordbox.Track) bool {
	return t1.Name == t2.Name &&
		t1.Artist == t2.Artist &&
		t1.Composer == t2.Composer &&
		t1.Album == t2.Album &&
		t1.Comments == t2.Comments &&
		t1.DiscNumber == t2.DiscNumber &&
		t1.TrackNumber == t2.TrackNumber &&
		t1.Year == t2.Year
}

func getLibraryMatches(sourceLibrary, destinationLibrary *rekordbox.RekordboxLibraryXML) *TrackMatches {
	trackMatches := &TrackMatches{}

	for _, sourceTrack := range sourceLibrary.Collection.TRACK {
		matchedTrack, err := matchTrackInLibrary(sourceTrack, destinationLibrary)
		if err != nil {
			trackMatches.tracks = append(trackMatches.tracks, &TrackMatch{
				SourceTrack: sourceTrack,
				FoundMatch:  false,
			})
			trackMatches.unmatchedCount++
		} else {
			trackMatches.tracks = append(trackMatches.tracks, &TrackMatch{
				SourceTrack:      sourceTrack,
				DestinationTrack: *matchedTrack,
				FoundMatch:       true,
			})
			trackMatches.matchedCount++
		}
	}
	return trackMatches
}

func mergeLibraryData(sourceLibrary, destinationLibrary *rekordbox.RekordboxLibraryXML) *rekordbox.RekordboxLibraryXML {
	trackMatches := getLibraryMatches(sourceLibrary, destinationLibrary)

	fmt.Println("\n--- Summary ---")
	fmt.Printf("Source Tracks:      %d\n", len(sourceLibrary.Collection.TRACK))
	fmt.Printf("Destination Tracks: %d\n", len(destinationLibrary.Collection.TRACK))
	fmt.Printf("Matched:            %d\n", trackMatches.matchedCount)
	fmt.Printf("Unmatched:          %d\n", trackMatches.unmatchedCount)

	var mergedTracks []rekordbox.Track
	for _, trackMatch := range trackMatches.tracks {
		if trackMatch.FoundMatch {
			mergedTrack := trackMatch.DestinationTrack
			mergedTrack.Tempo = trackMatch.SourceTrack.Tempo
			mergedTracks = append(mergedTracks, mergedTrack)
		}
	}

	return &rekordbox.RekordboxLibraryXML{
		Version: destinationLibrary.Version,
		Product: destinationLibrary.Product,
		Collection: rekordbox.Collection{
			Entries: int32(len(mergedTracks)),
			TRACK:   mergedTracks,
		},
		Playlists: destinationLibrary.Playlists,
	}
}
