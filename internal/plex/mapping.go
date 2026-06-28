package plex

import (
	"github.com/llttlltt/dj-library-tools/internal/models"
	"github.com/llttlltt/dj-library-tools/internal/utils"
)

func ToNeutralTrack(t Track) models.Track {
	mt := models.Track{
		ID:     t.RatingKey,
		Title:  t.Title,
		Artist: t.Artist,
		Album:  t.Album,
		BPM:    t.BPM,
		Key:    t.KeyTag,
		Rating: utils.NormalizeRating(t.UserRating, 10.0), // Plex uses a 10-point internal scale
		ImplementationState: t,
	}
	if len(t.Media) > 0 && len(t.Media[0].Part) > 0 {
		mt.Location = t.Media[0].Part[0].File
	}
	return mt
}

func ToNeutralGroup(p Playlist) models.ResourceGroup {
	return models.ResourceGroup{
		ID:    p.RatingKey,
		Name:  p.Title,
		Items: p.LeafCount,
		Type:  models.GroupTypePlaylist,
		ImplementationState: p,
	}
}
