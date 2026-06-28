package rb

import (
	"strconv"

	"github.com/llttlltt/dj-library-tools/internal/models"
	"github.com/llttlltt/dj-library-tools/internal/rekordbox"
)

func ToNeutralTrack(t rekordbox.Track) models.Track {
	mt := models.Track{
		ID:           strconv.Itoa(t.TrackID),
		Title:        t.Name,
		Artist:       t.Artist,
		Album:        t.Album,
		Key:          t.Tonality,
		Genre:        t.Genre,
		Comment:      t.Comments,
		Label:        t.Label,
		Year:         int(t.Year),
		Location:     t.Location,
		Rating:       int(t.Rating / 51),
		Plays:        int(t.PlayCount),
		DateAdded:    t.DateAdded,
		DateModified: t.DateModified,
		Bitrate:      int(t.BitRate),
		SampleRate:   int(t.SampleRate),
		Size:         t.Size,
		Remixer:      t.Remixer,
		Mix:          t.Mix,
		Raw:          t,
	}

	if t.AverageBpm != "" {
		mt.BPM, _ = strconv.ParseFloat(t.AverageBpm, 64)
	}

	mt.HotCues = 0
	mt.MemoryCues = 0
	for _, pm := range t.PositionMark {
		if pm.Num != -1 {
			mt.HotCues++
		} else {
			mt.MemoryCues++
		}
	}
	mt.BeatgridCount = len(t.Tempo)
	mt.Color = t.Colour

	return mt
}

func ToNeutralGroup(n rekordbox.Node, parentFolder string) models.ResourceGroup {
	// Folders (Type=0) store their child-node count in Count;
	// playlists (Type=1) store their track count in Entries.
	items := rekordbox.DerefInt32(n.Entries)
	if n.Type == 0 {
		items = rekordbox.DerefInt32(n.Count)
	}
	return models.ResourceGroup{
		Name:         n.Name,
		Items:        int(items),
		ParentFolder: parentFolder,
		Type:         models.GroupType(n.Type),
		Raw:          n,
	}
}
