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

	return &rekordboxLibrary, nil
}

func WriteRekordboxLibrary(path string, library *RekordboxLibraryXML) error {
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
