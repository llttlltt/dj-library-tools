package rekordbox

import (
	"context"
	"testing"

	"github.com/llttlltt/dj-library-tools/internal/core/models"
	"github.com/llttlltt/dj-library-tools/internal/providers"
	"github.com/stretchr/testify/assert"
)

func TestFixDuplicateTracks(t *testing.T) {
	xml := &RekordboxLibraryXML{
		Collection: Collection{
			TRACK: []Track{
				{TrackID: 1, Name: "Title", Artist: "Artist", Size: 100},
				{TrackID: 2, Name: "Title", Artist: "Artist", Size: 100}, // Duplicate
				{TrackID: 3, Name: "Other", Artist: "Other", Size: 200},
			},
			Entries: 3,
		},
	}

	sel := provider.Selection{
		Tracks: []models.Track{
			ToNeutralTrack(xml.Collection.TRACK[0]),
			ToNeutralTrack(xml.Collection.TRACK[1]),
			ToNeutralTrack(xml.Collection.TRACK[2]),
		},
	}

	res, err := FixDuplicateTracks(context.Background(), xml, sel, true)
	assert.NoError(t, err)
	assert.Equal(t, 1, res.TotalApplied)
	assert.Equal(t, 2, int(xml.Collection.Entries))
	assert.Len(t, xml.Collection.TRACK, 2)
}

func TestFixDuplicateMembers(t *testing.T) {
	xml := &RekordboxLibraryXML{
		Collection: Collection{
			TRACK: []Track{
				{TrackID: 1, Name: "Title", Artist: "Artist"},
			},
		},
		Playlists: Playlists{
			Node: RootNode{
				Name: "ROOT",
				Type: 0,
				Node: []Node{
					{
						Name: "Test",
						Type: 1,
						TRACK: []struct {
							Key string `xml:"Key,attr"`
						}{
							{Key: "1"},
							{Key: "1"}, // Duplicate member
						},
					},
				},
			},
		},
	}

	sel := provider.Selection{
		Items: []models.Resource{
			models.ResourceGroup{Name: "Test", Kind: models.GroupKindPlaylist},
		},
	}

	res, err := FixDuplicateMembers(context.Background(), xml, sel, true)
	assert.NoError(t, err)
	assert.Equal(t, 1, res.TotalApplied)
	assert.Len(t, xml.Playlists.Node.Node[0].TRACK, 1)
}
