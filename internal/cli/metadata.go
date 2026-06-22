package cli

import (
	"fmt"

	"github.com/llttlltt/dj-library-tools/internal/utils"
	"github.com/llttlltt/dj-library-tools/pkg/rekordbox"
	"github.com/spf13/cobra"
)

var (
	sourceXML      string
	destinationXML string
	outputXML      string
	forceMetadata  bool
)

var metadataCmd = &cobra.Command{
	Use:   "metadata [flags]",
	Short: "Manage track metadata between Rekordbox XML libraries",
	Long: `Reads two Rekordbox XML libraries: a source library from which Tempo
markers are copied, and a destination library whose tracks receive them.
Tracks are matched by strict metadata equality (Name, Artist, Album, etc.).
A merged library is written to the output path.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if sourceXML == "" && destinationXML == "" {
			return cmd.Help()
		}
		if sourceXML == "" {
			return fmt.Errorf("--source is required")
		}
		if destinationXML == "" {
			return fmt.Errorf("--destination is required")
		}
		if outputXML == "" {
			return fmt.Errorf("--output is required")
		}

		if err := utils.CheckFileOverwrite(outputXML, forceMetadata); err != nil {
			return fmt.Errorf("output file validation failed: %w", err)
		}

		sourceLibrary, err := rekordbox.ReadRekordboxLibrary(sourceXML)
		if err != nil {
			return fmt.Errorf("error reading source library %q: %w", sourceXML, err)
		}

		destinationLibrary, err := rekordbox.ReadRekordboxLibrary(destinationXML)
		if err != nil {
			return fmt.Errorf("error reading destination library %q: %w", destinationXML, err)
		}

		mergedLibrary := mergeLibraryData(sourceLibrary, destinationLibrary)

		if err := rekordbox.WriteRekordboxLibrary(outputXML, mergedLibrary); err != nil {
			return fmt.Errorf("error writing merged library to %q: %w", outputXML, err)
		}
		fmt.Printf("Successfully wrote merged library to %q\n", outputXML)
		return nil
	},
}

func init() {
	metadataCmd.Flags().StringVarP(&sourceXML, "source", "s", "", "Source Rekordbox XML (Tempo markers are read from here)")
	metadataCmd.Flags().StringVarP(&destinationXML, "destination", "d", "", "Destination Rekordbox XML (tracks receive the Tempo markers)")
	metadataCmd.Flags().StringVarP(&outputXML, "output", "o", "", "Output path for the merged Rekordbox XML")
	metadataCmd.Flags().BoolVarP(&forceMetadata, "force", "f", false, "Overwrite output file if it already exists")
	RootCmd.AddCommand(metadataCmd)
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
