package main

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
	forceMove      bool
)

var metadataCmd = &cobra.Command{
	Use:   "metadata",
	Short: "Manage track metadata operations between Rekordbox XML libraries",
}

var moveCmd = &cobra.Command{
	Use:   "move",
	Short: "Move metadata (Tempo) from source to destination XML",
	Long: `This command reads two Rekordbox XML libraries: a 'source' library 
from which specific metadata (currently only Tempo) will be taken, and a 
'destination' library whose tracks will receive this metadata if a match 
is found. A new merged library is created with the updated tracks. 
Tracks are matched using strict metadata equality (Name, Artist, Album, etc.).`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := utils.CheckFileOverwrite(outputXML, forceMove); err != nil {
			return fmt.Errorf("output file validation failed: %w", err)
		}

		sourceLibrary, err := rekordbox.ReadRekordboxLibrary(sourceXML)
		if err != nil {
			return fmt.Errorf("error getting source library '%s': %w", sourceXML, err)
		}

		destinationLibrary, err := rekordbox.ReadRekordboxLibrary(destinationXML)
		if err != nil {
			return fmt.Errorf("error getting destination library '%s': %w", destinationXML, err)
		}

		mergedLibrary := mergeLibraryData(sourceLibrary, destinationLibrary)

		if err := rekordbox.WriteRekordboxLibrary(outputXML, mergedLibrary); err != nil {
			return fmt.Errorf("error writing merged library to '%s': %w", outputXML, err)
		}
		fmt.Printf("Successfully wrote merged library to '%s'\n", outputXML)
		return nil
	},
}

func init() {
	moveCmd.Flags().StringVarP(&sourceXML, "source", "s", "", "Path to the source Rekordbox XML library (required)")
	moveCmd.Flags().StringVarP(&destinationXML, "destination", "d", "", "Path to the destination Rekordbox XML library (required)")
	moveCmd.Flags().StringVarP(&outputXML, "output", "o", "", "Path where the merged Rekordbox XML library will be saved (required)")
	moveCmd.Flags().BoolVarP(&forceMove, "force", "f", false, "Force overwrite of the output file if it already exists")

	moveCmd.MarkFlagRequired("source")
	moveCmd.MarkFlagRequired("destination")
	moveCmd.MarkFlagRequired("output")

	metadataCmd.AddCommand(moveCmd)
	rootCmd.AddCommand(metadataCmd)
}

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
	return nil, fmt.Errorf("no match found for '%s'", sourceTrack.Name)
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
