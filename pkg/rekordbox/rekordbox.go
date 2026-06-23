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

	output, err := xml.MarshalIndent(library, "", format.Indent)
	if err != nil {
		return fmt.Errorf("failed to marshall library data: %w", err)
	}

	if format.SelfClosing {
		output = postProcessSelfClosing(output)
	}

	if format.LineEnding != "\n" {
		output = bytes.ReplaceAll(output, []byte("\n"), []byte(format.LineEnding))
	}

	if format.LineLength > 0 {
		output = postProcessWrapping(output, format.LineLength, format.Indent)
	}

	xmlHeader := []byte(`<?xml version="1.0" encoding="UTF-8"?>` + format.LineEnding + format.LineEnding)
	output = append(xmlHeader, output...)

	err = os.WriteFile(path, output, 0644)
	if err != nil {
		return fmt.Errorf("failed to write XML to file %s: %w", path, err)
	}

	return nil
}

var tagRegex = regexp.MustCompile(`<([a-zA-Z0-9_]+)([^>]*?)></[a-zA-Z0-9_]+>`)

func postProcessSelfClosing(data []byte) []byte {
	// We want to make sure the opening and closing tags match if we use backreferences, 
	// but Go's regexp doesn't support them. For our XML, we can do a simpler replace
	// or use a more robust approach if needed.
	// Since we know the structure, a simple regex should work for most empty tags.
	return tagRegex.ReplaceAll(data, []byte(`<${1}${2}/>`))
}

func postProcessWrapping(data []byte, lineLength int, indent string) []byte {
	if lineLength <= 0 {
		return data
	}

	lines := bytes.Split(data, []byte("\n"))
	var result [][]byte

	for _, line := range lines {
		if len(line) <= lineLength {
			result = append(result, line)
			continue
		}

		// Look for attributes to wrap
		// Regex for: space + attribute name + ="
		attrRegex := regexp.MustCompile(` [a-zA-Z0-9]+="`)
		matches := attrRegex.FindAllIndex(line, -1)
		if len(matches) == 0 {
			result = append(result, line)
			continue
		}

		currentLine := line
		var wrappedLines [][]byte

		for i := len(matches) - 1; i >= 0; i-- {
			pos := matches[i][0]
			if pos > lineLength {
				// We wrap from the end to avoid shifting indices
				head := currentLine[:pos]
				tail := currentLine[pos+1:] // Skip the space
				
				// Keep the head and start a new line for the tail
				// Note: recursion/looping would be better for multiple wraps
				// This simple version handles the most common case
				wrappedLines = append([][]byte{append([]byte(indent+indent+indent), tail...)}, wrappedLines...)
				currentLine = head
			}
		}
		result = append(result, currentLine)
		result = append(result, wrappedLines...)
	}

	return bytes.Join(result, []byte("\n"))
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
	newPlaylists, err := xml.MarshalIndent(library.Playlists, "", format.Indent)
	if err != nil {
		return fmt.Errorf("failed to marshal playlists: %w", err)
	}

	if format.SelfClosing {
		newPlaylists = postProcessSelfClosing(newPlaylists)
	}
	
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
