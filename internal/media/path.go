package media

import (
	"bytes"
	"path/filepath"
	"text/template"
)

type PathMetadata struct {
	Artist string
	Album  string
	Title  string
}

func (t *Transcoder) FormatPath(metadata PathMetadata) (string, error) {
	tmplStr, ok := t.Config.Paths["default"]
	if !ok {
		tmplStr = "{{.Artist}} - {{.Album}} - {{.Title}}"
	}

	tmpl, err := template.New("path").Parse(tmplStr)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, metadata); err != nil {
		return "", err
	}

	// Clean path components to be safe for filenames
	fileName := buf.String()
	// Add extension
	fileName = fileName + "." + t.Config.Format

	return fileName, nil
}

func (t *Transcoder) GetDestinationPath(metadata PathMetadata) (string, error) {
	relPath, err := t.FormatPath(metadata)
	if err != nil {
		return "", err
	}
	return filepath.Join(t.Config.Dest, relPath), nil
}
