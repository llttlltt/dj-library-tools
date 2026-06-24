package rekordbox

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"os"
	"regexp"
)

func ReadRekordboxLibrary(path string) (*RekordboxLibraryXML, error) {
	xmlFile, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open XML file at %s: %w", path, err)
	}
	defer xmlFile.Close()

	xmlFileBytes, err := io.ReadAll(xmlFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read XML file %s: %w", path, err)
	}

	var rekordboxLibrary RekordboxLibraryXML
	err = xml.Unmarshal(xmlFileBytes, &rekordboxLibrary)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal XML data: %w", err)
	}

	rekordboxLibrary.Format = DetectFormat(xmlFileBytes)
	rekordboxLibrary.OriginalRaw = xmlFileBytes

	return &rekordboxLibrary, nil
}

func WriteRekordboxLibrary(path string, library *RekordboxLibraryXML) error {
	// If the collection hasn't changed and we have the original raw bytes,
	// we perform a surgical save of just the playlists section.
	if !library.CollectionChanged && len(library.OriginalRaw) > 0 {
		return writeSurgically(path, library)
	}

	format := library.Format
	if format == nil {
		format = DefaultFormat()
	}

	// Use our high-fidelity TokenStreamFormatter instead of standard xml.Marshal
	marshaled, err := xml.Marshal(library)
	if err != nil {
		return fmt.Errorf("failed to marshal library data: %w", err)
	}

	formatter := NewTokenStreamFormatter(format)
	var output bytes.Buffer
	xmlHeader := `<?xml version="1.0" encoding="UTF-8"?>` + format.LineEnding + format.LineEnding
	output.WriteString(xmlHeader)

	if err := formatter.Format(bytes.NewReader(marshaled), &output); err != nil {
		return fmt.Errorf("failed to format XML: %w", err)
	}

	finalOutput := output.Bytes()
	if format.LineEnding != "\n" {
		finalOutput = bytes.ReplaceAll(finalOutput, []byte("\n"), []byte(format.LineEnding))
	}

	return os.WriteFile(path, finalOutput, 0644)
}

var playlistBlockRegex = regexp.MustCompile(`(?s)<PLAYLISTS>.*</PLAYLISTS>`)

func writeSurgically(path string, library *RekordboxLibraryXML) error {
	// If nothing changed at all, just return
	if !library.CollectionChanged && !library.PlaylistsChanged {
		return nil
	}

	format := library.Format
	if format == nil {
		format = DefaultFormat()
	}

	// Marshal only the playlists section
	newPlaylistsRaw, err := xml.Marshal(library.Playlists)
	if err != nil {
		return fmt.Errorf("failed to marshal playlists: %w", err)
	}

	formatter := NewTokenStreamFormatter(format)
	var formattedPlaylists bytes.Buffer
	if err := formatter.Format(bytes.NewReader(newPlaylistsRaw), &formattedPlaylists); err != nil {
		return fmt.Errorf("failed to format playlists: %w", err)
	}
	newPlaylists := formattedPlaylists.Bytes()

	// Apply line endings
	if format.LineEnding != "\n" {
		newPlaylists = bytes.ReplaceAll(newPlaylists, []byte("\n"), []byte(format.LineEnding))
	}

	// Find the playlists block in the original file
	loc := playlistBlockRegex.FindIndex(library.OriginalRaw)
	if loc == nil {
		// Fallback to full write if we can't find the block
		library.CollectionChanged = true
		return WriteRekordboxLibrary(path, library)
	}

	// Stitch it together
	var output bytes.Buffer
	output.Write(library.OriginalRaw[:loc[0]])
	output.Write(newPlaylists)
	output.Write(library.OriginalRaw[loc[1]:])

	return os.WriteFile(path, output.Bytes(), 0644)
}
