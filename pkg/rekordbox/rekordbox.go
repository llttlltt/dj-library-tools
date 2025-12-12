package rekordbox

import (
	"encoding/xml"
	"fmt"
	"io"
	"os"
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

	return &rekordboxLibrary, nil
}

func WriteRekordboxLibrary(path string, library *RekordboxLibraryXML) error {
	output, err := xml.MarshalIndent(library, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshall library data: %w", err)
	}
	xmlHeader := []byte(`<?xml version="1.0" encoding="UTF-8"?>` + "\n")
	output = append(xmlHeader, output...)

	err = os.WriteFile(path, output, 0644)
	if err != nil {
		return fmt.Errorf("failed to write XML to file %s: %w", path, err)
	}

	return nil
}
