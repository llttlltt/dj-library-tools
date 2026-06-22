package media

import (
	"bytes"
	"path/filepath"
	"strings"
	"text/template"
)

type PathMetadata struct {
	Artist string
	Album  string
	Title  string
}

func (t *Transcoder) ApplyPathMap(source string) string {
	if t.Config.PathMaps == nil {
		return source
	}

	for remote, local := range t.Config.PathMaps {
		if strings.HasPrefix(source, remote) {
			return strings.Replace(source, remote, local, 1)
		}
	}
	return source
}

// sanitizePathComponent replaces characters that are invalid or problematic in
// file and directory names across macOS, Linux, and Windows.
func sanitizePathComponent(s string) string {
	replacer := strings.NewReplacer(
		"/", "-",
		"\\", "-",
		":", "-",
		"*", "-",
		"?", "",
		`"`, "",
		"<", "",
		">", "",
		"|", "-",
	)
	return strings.TrimSpace(replacer.Replace(s))
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

	// Sanitize each component before template execution so that characters
	// such as '/' in an artist name cannot create unintended subdirectories.
	sanitized := PathMetadata{
		Artist: sanitizePathComponent(metadata.Artist),
		Album:  sanitizePathComponent(metadata.Album),
		Title:  sanitizePathComponent(metadata.Title),
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, sanitized); err != nil {
		return "", err
	}

	return buf.String() + "." + t.Config.Format, nil
}

func (t *Transcoder) GetDestinationPath(metadata PathMetadata) (string, error) {
	relPath, err := t.FormatPath(metadata)
	if err != nil {
		return "", err
	}
	return filepath.Join(t.Config.Dest, relPath), nil
}
