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
	LineLength  int
}

// DefaultFormat returns the default Rekordbox XML formatting.
func DefaultFormat() *XMLFormat {
	return &XMLFormat{
		Indent:      "  ",
		SelfClosing: true,
		LineEnding:  "\n",
		LineLength:  0, // No wrapping by default
	}
}

// DetectFormat attempts to guess the formatting from the XML data.
func DetectFormat(data []byte) *XMLFormat {
	format := DefaultFormat()

	// Detect line endings
	if bytes.Contains(data, []byte("\r\n")) {
		format.LineEnding = "\r\n"
	}

	// Detect indentation and average line length
	lines := bytes.Split(data, []byte(format.LineEnding))
	totalLength := 0
	lineCount := 0
	for _, line := range lines {
		if len(line) == 0 {
			continue
		}
		// Only look for indentation on tags that aren't the first line
		if bytes.Contains(line, []byte("<?xml")) {
			continue
		}

		// Detect indentation
		if format.Indent == "  " { // Still default
			matches := regexp.MustCompile(`^([ \t]+)<`).FindSubmatch(line)
			if len(matches) > 1 {
				format.Indent = string(matches[1])
			}
		}

		// Track line lengths to detect wrapping
		if len(line) > 10 { // ignore very short lines
			totalLength += len(line)
			lineCount++
		}
	}

	if lineCount > 0 {
		avgLen := totalLength / lineCount
		// If average line length is significantly less than standard 120-150, 
		// it's likely the file has attribute wrapping.
		// The user noticed a ~6000 line difference in a ~300k line file.
		// 315757 - 309037 = 6720 lines.
		// Standard Rekordbox seems to wrap at 80 characters for attributes.
		if avgLen < 150 {
			format.LineLength = 80
		} else {
			format.LineLength = 0 // No wrapping
		}
	}

	// Always prefer self-closing for Rekordbox
	format.SelfClosing = true

	return format
}
