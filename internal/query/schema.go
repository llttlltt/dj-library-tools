package query

import (
	"fmt"
	"strconv"

	"github.com/llttlltt/dj-library-tools/internal/models"
)

// TrackAccessor is a function that extracts a string value from a neutral Track.
type TrackAccessor func(models.Track) string

// GroupAccessor is a function that extracts a string value from a ResourceGroup.
type GroupAccessor func(models.ResourceGroup) string

// TrackAccessors maps query field names to their extraction logic.
var TrackAccessors = map[string]TrackAccessor{
	"id":         func(t models.Track) string { return t.ID },
	"title":      func(t models.Track) string { return t.Title },
	"artist":     func(t models.Track) string { return t.Artist },
	"album":      func(t models.Track) string { return t.Album },
	"genre":      func(t models.Track) string { return t.Genre },
	"comment":    func(t models.Track) string { return t.Comment },
	"label":      func(t models.Track) string { return t.Label },
	"year":       func(t models.Track) string { return strconv.Itoa(t.Year) },
	"color":      func(t models.Track) string { return t.Color },
	"bpm":        func(t models.Track) string { return fmt.Sprintf("%.2f", t.BPM) },
	"key":        func(t models.Track) string { return t.Key },
	"location":   func(t models.Track) string { return t.Location },
	"display":    func(t models.Track) string { return t.Display },
	"rating":     func(t models.Track) string { return strconv.Itoa(t.Rating) },
	"plays":      func(t models.Track) string { return strconv.Itoa(t.Plays) },
	"added":      func(t models.Track) string { return t.DateAdded },
	"modified":   func(t models.Track) string { return t.DateModified },
	"bitrate":    func(t models.Track) string { return strconv.Itoa(t.Bitrate) },
	"samplerate": func(t models.Track) string { return strconv.Itoa(t.SampleRate) },
	"size":       func(t models.Track) string { return strconv.FormatInt(t.Size, 10) },
	"remixer":    func(t models.Track) string { return t.Remixer },
	"mix":        func(t models.Track) string { return t.Mix },
	"hotcues":    func(t models.Track) string { return strconv.Itoa(countCues(t, models.CueTypeHot)) },
	"memorycues": func(t models.Track) string { return strconv.Itoa(countCues(t, models.CueTypeMemory)) },
	"beatgrids":  func(t models.Track) string { return strconv.Itoa(len(t.TempoMarkers)) },
}

// GroupAccessors maps query field names to their extraction logic.
var GroupAccessors = map[string]GroupAccessor{
	"name":   func(g models.ResourceGroup) string { return g.Name },
	"parent": func(g models.ResourceGroup) string { return g.ParentFolder },
	"folder": func(g models.ResourceGroup) string { return g.ParentFolder },
	"items":  func(g models.ResourceGroup) string { return strconv.Itoa(g.Items) },
	"kind":   func(g models.ResourceGroup) string { return string(g.Kind) },
}

func countCues(t models.Track, cueType models.CueType) int {
	count := 0
	for _, cp := range t.CuePoints {
		if cp.Type == cueType {
			count++
		}
	}
	return count
}
