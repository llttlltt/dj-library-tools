package m3u

import (
	"fmt"
	"io"
)

// WriteM3U8EntryRaw writes a raw #EXTINF line followed by the file path.
func WriteM3U8EntryRaw(w io.Writer, display string, path string, duration float64) error {
	d := duration
	if d == 0 {
		d = -1
	}
	line := fmt.Sprintf("#EXTINF:%.0f,%s\n", d, display)
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
