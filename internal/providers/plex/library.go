package plex

import (
	"context"
	"github.com/llttlltt/dj-library-tools/internal/core/models"
)

type Library struct {
	client  *Client
	baseURL string
}

func NewLibrary(client *Client, baseURL string) *Library {
	return &Library{client: client, baseURL: baseURL}
}

func (l *Library) GetResources(kind string) []models.Resource {
	ctx := context.Background()
	var results []models.Resource
	
	switch kind {
	case "track":
		tracks, err := l.client.GetAllTracks(ctx, l.baseURL)
		if err != nil {
			return nil
		}
		for _, t := range tracks {
			results = append(results, ToNeutralTrack(t))
		}
	case "group":
		playlists, err := l.client.GetPlaylists(ctx, l.baseURL)
		if err != nil {
			return nil
		}
		for _, p := range playlists {
			results = append(results, ToNeutralGroup(p))
		}
	}
	return results
}

func (l *Library) GetMembershipMap() map[string][]models.PlaylistMembership {
	// Note: For now, fetching memberships in Plex is expensive (N API calls)
	// We'll return an empty map and eventually optimize this if Plex becomes 
	// a primary source provider.
	return make(map[string][]models.PlaylistMembership)
}
