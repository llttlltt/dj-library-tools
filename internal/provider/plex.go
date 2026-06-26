package provider

import (
	"context"
	"fmt"
	"strings"

	"github.com/llttlltt/dj-library-tools/internal/plex"
	"github.com/llttlltt/dj-library-tools/pkg/rekordbox"
)

type PlexProvider struct {
	client *plex.Client
	host   string
	port   int
}

func NewPlexProvider(token string, host string, port int) *PlexProvider {
	return &PlexProvider{
		client: plex.NewClient(token),
		host:   host,
		port:   port,
	}
}

func (p *PlexProvider) Name() string {
	return "plex"
}

func (p *PlexProvider) GetTracks(query string) ([]rekordbox.Track, error) {
	ctx := context.Background()
	baseURL, err := p.resolveBaseURL(ctx)
	if err != nil {
		return nil, err
	}

	// For now, we assume query is a playlist ID (RatingKey) as seen in current list.go
	// In the future, this could be a more complex query.
	path := "/playlists/" + query + "/items"
	plexTracks, err := p.client.GetPlaylistTracks(ctx, baseURL, path)
	if err != nil {
		return nil, err
	}

	var tracks []rekordbox.Track
	for _, pt := range plexTracks {
		t := rekordbox.Track{
			Name:   pt.Title,
			Artist: pt.Artist,
			Album:  pt.Album,
		}
		if len(pt.Media) > 0 && len(pt.Media[0].Part) > 0 {
			t.Location = pt.Media[0].Part[0].File
		}
		tracks = append(tracks, t)
	}

	return tracks, nil
}

func (p *PlexProvider) GetPlaylists(query string) ([]NodeResult, error) {
	ctx := context.Background()
	baseURL, err := p.resolveBaseURL(ctx)
	if err != nil {
		return nil, err
	}

	plexPlaylists, err := p.client.GetPlaylists(ctx, baseURL)
	if err != nil {
		return nil, err
	}

	var results []NodeResult
	for _, pl := range plexPlaylists {
		if query != "" && !strings.Contains(strings.ToLower(pl.Title), strings.ToLower(query)) {
			continue
		}
		results = append(results, NodeResult{
			Name:    pl.Title,
			Entries: pl.LeafCount,
			Raw:     pl,
		})
	}

	return results, nil
}

func (p *PlexProvider) resolveBaseURL(ctx context.Context) (string, error) {
	if p.host != "" {
		port := p.port
		if port == 0 {
			port = 32400
		}
		return fmt.Sprintf("http://%s:%d", p.host, port), nil
	}

	// If no host provided, we need to find the best server connection.
	// This is a bit complex as Plex has multiple servers.
	// For simplicity, we'll look for the first available server.
	resources, err := p.client.GetResources(ctx)
	if err != nil {
		return "", err
	}

	for _, res := range resources {
		if res.Provides != "server" {
			continue
		}
		// Probe best connection for this server
		probe, err := p.client.ProbeBestConnection(res)
		if err == nil {
			return probe.BaseURL, nil
		}
	}

	return "", fmt.Errorf("could not find an active Plex server")
}
