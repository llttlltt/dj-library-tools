package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/llttlltt/dj-library-tools/internal/utils"
	"github.com/llttlltt/dj-library-tools/pkg/rekordbox"
	"github.com/urfave/cli/v3"
)

type TrackMatch struct {
	SourceTrack      rekordbox.Track
	DestinationTrack rekordbox.Track
	FoundMatch       bool
}

func main() {
	var forceOverwrite bool
	var sourceLibraryPath string
	var destinationLibraryPath string
	var outputFileExtension string
	var outputLibraryPath string

	forceFlag := &cli.BoolFlag{
		Name:        "force",
		Aliases:     []string{"f"},
		Usage:       "Force overwrite of the output file if it already exists.",
		Destination: &forceOverwrite,
		Value:       false,
	}
	sourcePathFlag := &cli.StringFlag{
		Name:        "source",
		Usage:       "Path to the source Rekordbox XML library. Metadata will be extracted from tracks in this file.",
		Aliases:     []string{"s"},
		TakesFile:   true,
		Required:    true,
		Destination: &sourceLibraryPath,
	}
	destinationPathFlag := &cli.StringFlag{
		Name:        "destination",
		Usage:       "Path to the destination Rekordbox XML library. Tracks in this file will receive metadata.",
		Aliases:     []string{"d"},
		TakesFile:   true,
		Required:    true,
		Destination: &destinationLibraryPath,
	}
	outputFileExtensionFlag := &cli.StringFlag{
		Name:        "output-extension",
		Value:       "xml",
		Usage:       "Internal: Sets the file extension for output files. (Developers only)",
		Destination: &outputFileExtension,
		Hidden:      true,
	}
	outputPathFlag := &cli.StringFlag{
		Name:        "output",
		Usage:       "Path where the new merged Rekordbox XML library will be saved.",
		TakesFile:   true,
		Aliases:     []string{"o"},
		Required:    true,
		Destination: &outputLibraryPath,
	}
	cmd := &cli.Command{
		Name:                  "rb-cli",
		Usage:                 "A CLI tool for managing Rekordbox XML libraries.",
		EnableShellCompletion: true,
		Commands: []*cli.Command{
			{
				Name:  "metadata",
				Usage: "Manage track metadata operations between libraries.",
				Flags: []cli.Flag{
					forceFlag,
					sourcePathFlag,
					destinationPathFlag,
					outputFileExtensionFlag,
					outputPathFlag,
				},
				Commands: []*cli.Command{
					{
						Name:  "move",
						Usage: "Move metadata from a source library's tracks to matching tracks in a destination library.",
						Description: "This command reads two Rekordbox XML libraries: a 'source' library " +
							"from which specific metadata (currently only Tempo) will be taken, and a " +
							"'destination' library whose tracks will receive this metadata if a match " +
							"is found. A new merged library is created with the updated tracks. " +
							"Tracks are matched using strict metadata equality (Name, Artist, Album, etc.).",
						Action: func(ctx context.Context, cmd *cli.Command) error {
							finalOutputPath := utils.EnforceExtension(outputLibraryPath, outputFileExtension)
							if err := utils.CheckFileOverwrite(finalOutputPath, forceOverwrite); err != nil {
								return fmt.Errorf("output file validation failed: %w", err)
							}

							sourceLibrary, err := rekordbox.ReadRekordboxLibrary(sourceLibraryPath)
							if err != nil {
								return fmt.Errorf("error getting source library '%s': %w", sourceLibraryPath, err)
							}

							destinationLibrary, err := rekordbox.ReadRekordboxLibrary(destinationLibraryPath)
							if err != nil {
								return fmt.Errorf("error getting destination library '%s': %w", destinationLibraryPath, err)
							}

							mergedLibrary := mergeLibraryData(sourceLibrary, destinationLibrary)

							// 3. Write the merged library to the validated and finalized path
							if err := rekordbox.WriteRekordboxLibrary(finalOutputPath, mergedLibrary); err != nil {
								return fmt.Errorf("error writing merged library to '%s': %w", finalOutputPath, err)
							}
							fmt.Printf("Successfully wrote merged library to '%s'\n", finalOutputPath)
							return nil
						},
					},
				},
			},
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}

}

func matchTrackInLibrary(sourceTrack rekordbox.Track, destinationLibrary *rekordbox.RekordboxLibraryXML) (*rekordbox.Track, error) {
	for _, destinationTrack := range destinationLibrary.Collection.TRACK {
		isMetadataMatch, _ := tracksMatchMetadata(sourceTrack, destinationTrack)
		if isMetadataMatch {
			return &destinationTrack, nil
		}
	}

	return &rekordbox.Track{}, fmt.Errorf("no match found for '%s' (%s)", sourceTrack.Name, sourceTrack.Location)
}

func tracksMatchMetadata(t1, t2 rekordbox.Track) (bool, []string) {
	mismatches := []string{}

	if t1.Name != t2.Name {
		mismatches = append(mismatches, fmt.Sprintf("Name: '%s' != '%s'", t1.Name, t2.Name))
	}
	if t1.Artist != t2.Artist {
		mismatches = append(mismatches, fmt.Sprintf("Artist: '%s' != '%s'", t1.Artist, t2.Artist))
	}
	if t1.Composer != t2.Composer {
		mismatches = append(mismatches, fmt.Sprintf("Composer: '%s' != '%s'", t1.Composer, t2.Composer))
	}
	if t1.Album != t2.Album {
		mismatches = append(mismatches, fmt.Sprintf("Album: '%s' != '%s'", t1.Album, t2.Album))
	}
	if t1.Comments != t2.Comments {
		mismatches = append(mismatches, fmt.Sprintf("Comments: '%s' != '%s'", t1.Comments, t2.Comments))
	}
	if t1.DiscNumber != t2.DiscNumber {
		mismatches = append(mismatches, fmt.Sprintf("DiscNumber: %d != %d", t1.DiscNumber, t2.DiscNumber))
	}
	if t1.TrackNumber != t2.TrackNumber {
		mismatches = append(mismatches, fmt.Sprintf("TrackNumber: %d != %d", t1.TrackNumber, t2.TrackNumber))
	}
	if t1.Year != t2.Year {
		mismatches = append(mismatches, fmt.Sprintf("Year: %d != %d", t1.Year, t2.Year))
	}

	return len(mismatches) == 0, mismatches
}

type TrackMatches struct {
	tracks         []*TrackMatch
	matchedCount   int
	unmatchedCount int
}

func getLibraryMatches(sourceLibrary, destinationLibrary *rekordbox.RekordboxLibraryXML) *TrackMatches {
	trackMatches := TrackMatches{}

	for _, sourceTrack := range sourceLibrary.Collection.TRACK {
		matchedTrack, err := matchTrackInLibrary(sourceTrack, destinationLibrary)

		if err != nil {
			trackMatches.tracks = append(trackMatches.tracks, &TrackMatch{
				SourceTrack:      sourceTrack,
				DestinationTrack: rekordbox.Track{},
				FoundMatch:       false,
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

	return &trackMatches
}

func mergeLibraryData(sourceLibrary, destinationLibrary *rekordbox.RekordboxLibraryXML) *rekordbox.RekordboxLibraryXML {
	trackMatches := getLibraryMatches(sourceLibrary, destinationLibrary)

	fmt.Println("\n--- Summary ---")
	fmt.Printf("Source Library Collection Entries:      %d\n", sourceLibrary.Collection.Entries)
	fmt.Printf("Destination Library Collection Entries: %d\n", destinationLibrary.Collection.Entries)
	fmt.Printf("Tracks Matched:   %d\n", trackMatches.matchedCount)
	fmt.Printf("Tracks Unmatched: %d\n", trackMatches.unmatchedCount)

	if trackMatches.unmatchedCount > 0 {
		fmt.Println("\n--- Unmatched Source Tracks ---")
		for _, track := range trackMatches.tracks {
			if !track.FoundMatch {

				fmt.Printf("  - %s (%s)\n", track.SourceTrack.Name, track.SourceTrack.Location)
			}
		}
	}

	var mergedTracks []rekordbox.Track
	for _, trackMatch := range trackMatches.tracks {
		if trackMatch.FoundMatch {
			mergedTrack := trackMatch.DestinationTrack
			mergedTrack.Tempo = trackMatch.SourceTrack.Tempo
			mergedTracks = append(mergedTracks, mergedTrack)
		} else {
			// handle no match case
		}
	}

	mergedCollection := rekordbox.Collection{
		Entries: int32(len(mergedTracks)),
		TRACK:   mergedTracks,
	}

	mergedLibrary := rekordbox.RekordboxLibraryXML{
		Version:    destinationLibrary.Version,
		Product:    destinationLibrary.Product,
		Collection: mergedCollection,
		Playlists:  destinationLibrary.Playlists,
	}
	return &mergedLibrary
}
