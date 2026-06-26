package provider

import (
	"context"
	"fmt"
	"strings"

	"github.com/llttlltt/dj-library-tools/internal/plex"
	"github.com/llttlltt/dj-library-tools/internal/query"
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

func (p *PlexProvider) GetTracks(queryString string) ([]rekordbox.Track, error) {
	ctx := context.Background()
	baseURL, err := p.resolveBaseURL(ctx)
	if err != nil {
		return nil, err
	}

	q := query.NewParser().Parse(queryString)
	if err := q.Validate(); err != nil {
		return nil, err
	}

	playlistID := ""
	if queryString != "" {
		if q.Root == nil {
			return nil, fmt.Errorf("query must specify a field (e.g. id:%q or name:%q)", queryString, queryString)
		}

		// Helper to extract fields from query
		var playlistName string
		var walk func(expr query.Expression)
		walk = func(expr query.Expression) {
			switch v := expr.(type) {
			case query.Comparison:
				switch strings.ToLower(v.Field) {
				case "id", "ratingkey":
					playlistID = v.Value
				case "name", "title", "playlist":
					playlistName = v.Value
				}
			case query.Logical:
				walk(v.Left)
				walk(v.Right)
			}
		}
		walk(q.Root)

		if playlistID == "" && playlistName != "" {
			// Resolve name to ID
			plexPlaylists, err := p.client.GetPlaylists(ctx, baseURL)
			if err != nil {
				return nil, err
			}
			for _, pl := range plexPlaylists {
				if strings.EqualFold(pl.Title, playlistName) {
					playlistID = pl.RatingKey
					break
				}
			}
			if playlistID == "" {
				return nil, fmt.Errorf("plex playlist %q not found", playlistName)
			}
		}
	}

	if playlistID == "" {
		return nil, fmt.Errorf("query must resolve to a playlist ID (use id:123 or name:'Summer')")
	}

	path := "/playlists/" + playlistID + "/items"
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

func (p *PlexProvider) GetPlaylists(queryString string) ([]NodeResult, error) {
	ctx := context.Background()
	baseURL, err := p.resolveBaseURL(ctx)
	if err != nil {
		return nil, err
	}

	q := query.NewParser().Parse(queryString)
	if err := q.Validate(); err != nil {
		return nil, err
	}

	plexPlaylists, err := p.client.GetPlaylists(ctx, baseURL)
	if err != nil {
		return nil, err
	}

	var results []NodeResult
	for _, pl := range plexPlaylists {
		if queryString != "" {
			matched := false
			// Simple evaluator for playlists
			var walk func(expr query.Expression) bool
			walk = func(expr query.Expression) bool {
				switch v := expr.(type) {
				case query.Comparison:
					f := strings.ToLower(v.Field)
					if f == "name" || f == "title" {
						return strings.Contains(strings.ToLower(pl.Title), strings.ToLower(v.Value))
					}
					if f == "id" || f == "ratingkey" {
						return pl.RatingKey == v.Value
					}
					return false
				case query.Logical:
					if v.Op == "AND" {
						return walk(v.Left) && walk(v.Right)
					}
					return walk(v.Left) || walk(v.Right)
				case query.Not:
					return !walk(v.Expr)
				}
				return true
			}
			matched = walk(q.Root)
			if !matched {
				continue
			}
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
