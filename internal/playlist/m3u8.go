package playlist

import (
	"fmt"
	"io"
)

// WriteM3U8Entry writes an #EXTINF line followed by the file path to the writer.
func WriteM3U8Entry(w io.Writer, metadata AudioMetadata, path string, duration float64) error {
	// #EXTINF:<duration>,<artist> - <title>
	// Use -1 for duration if unknown to signal player to auto-scan
	d := duration
	if d == 0 {
		d = -1
	}
	line := fmt.Sprintf("#EXTINF:%.0f,%s - %s\n", d, metadata.Artist, metadata.Title)
	if _, err := io.WriteString(w, line); err != nil {
		return err
	}

	if _, err := io.WriteString(w, path+"\n"); err != nil {
		return err
	}

	return nil
}

// WriteM3U8Header writes the #EXTM3U header.
func WriteM3U8Header(w io.Writer) error {
	_, err := io.WriteString(w, "#EXTM3U\n")
	return err
}
