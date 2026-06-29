package rekordbox

import (
	"html"
	"strconv"
	"strings"

	"github.com/llttlltt/dj-library-tools/internal/models"
)

func ToNeutralTrack(t Track) models.Track {
	mt := models.Track{
		ID:           strconv.Itoa(t.TrackID),
		Title:        strings.TrimSpace(html.UnescapeString(t.Name)),
		Artist:       strings.TrimSpace(html.UnescapeString(t.Artist)),
		Album:        strings.TrimSpace(html.UnescapeString(t.Album)),
		Key:          t.Tonality,
		Genre:        strings.TrimSpace(html.UnescapeString(t.Genre)),
		Comment:      strings.TrimSpace(html.UnescapeString(t.Comments)),
		Label:        strings.TrimSpace(html.UnescapeString(t.Label)),
		Year:         int(t.Year),
		Location:     t.Location,
		Rating:       int(t.Rating), // Rekordbox already uses 0-255 in XML
		Plays:        int(t.PlayCount),
		DateAdded:    t.DateAdded,
		DateModified: t.DateModified,
		Bitrate:      int(t.BitRate),
		SampleRate:   int(t.SampleRate),
		Size:         t.Size,
		Remixer:      t.Remixer,
		Mix:          t.Mix,
		ImplementationState: t,
	}

	if t.AverageBpm != "" {
		mt.BPM, _ = strconv.ParseFloat(t.AverageBpm, 64)
	}

	for _, pm := range t.PositionMark {
		cueType := models.CueTypeMemory
		if pm.Num != -1 {
			cueType = models.CueTypeHot
		}
		mt.CuePoints = append(mt.CuePoints, models.CuePoint{
			Name:     pm.Name,
			Position: parsePosition(pm.Start),
			Color:    GetHotCueColorName(pm),
			Type:     cueType,
			Index:    int(pm.Num),
		})
	}

	for _, tm := range t.Tempo {
		bpm, _ := strconv.ParseFloat(tm.Bpm, 64)
		mt.TempoMarkers = append(mt.TempoMarkers, models.TempoMarker{
			Position: parsePosition(tm.Inizio),
			BPM:      bpm,
		})
	}
	
	if t.Colour != "" {
		mt.Color = GetTrackColorName(t.Colour)
	} else {
		mt.Color = t.Colour
	}

	return mt
}

func parsePosition(s string) float64 {
	v, _ := strconv.ParseFloat(s, 64)
	return v
}

func ToNeutralGroup(n Node, parentFolder string) models.ResourceGroup {
	name := strings.TrimSpace(html.UnescapeString(n.Name))
	// Construction of the ID: use full path to ensure uniqueness
	id := name
	if parentFolder != "" {
		id = parentFolder + "/" + name
	}

	groupKind := models.GroupKindPlaylist
	if n.Type == 0 {
		groupKind = models.GroupKindFolder
	}

	// Folders (Type=0) store their child-node count in Count;
	// playlists (Type=1) store their track count in Entries.
	items := DerefInt32(n.Entries)
	if n.Type == 0 {
		items = DerefInt32(n.Count)
	}
	return models.ResourceGroup{
		ID:                  id,
		Name:                name,
		Items:               int(items),
		ParentFolder:        parentFolder,
		Kind:                groupKind,
		ImplementationState: n,
	}
}
