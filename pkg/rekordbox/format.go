package rekordbox

import (
	"bytes"
	"regexp"
)

// XMLFormat stores formatting preferences for Rekordbox XML.
type XMLFormat struct {
	Indent      string
	SelfClosing bool
	LineEnding  string
}

// DefaultFormat returns the default Rekordbox XML formatting.
func DefaultFormat() *XMLFormat {
	return &XMLFormat{
		Indent:      "  ",
		SelfClosing: true,
		LineEnding:  "\n",
	}
}

// DetectFormat attempts to guess the formatting from the XML data.
func DetectFormat(data []byte) *XMLFormat {
	format := DefaultFormat()

	// Detect line endings
	if bytes.Contains(data, []byte("\r\n")) {
		format.LineEnding = "\r\n"
	}

	// Detect indentation
	lines := bytes.Split(data, []byte(format.LineEnding))
	for _, line := range lines {
		// Only look for indentation on tags that aren't the first line
		if bytes.Contains(line, []byte("<?xml")) {
			continue
		}
		matches := regexp.MustCompile(`^([ \t]+)<`).FindSubmatch(line)
		if len(matches) > 1 {
			format.Indent = string(matches[1])
			break
		}
	}

	// Detect self-closing tags
	if bytes.Contains(data, []byte("/>")) {
		format.SelfClosing = true
	} else {
		format.SelfClosing = false
	}

	return format
}
