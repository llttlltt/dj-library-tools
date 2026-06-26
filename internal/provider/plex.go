package provider

import (
	"context"
	"fmt"
	"strings"

	"github.com/llttlltt/dj-library-tools/internal/models"
	"github.com/llttlltt/dj-library-tools/internal/plex"
	"github.com/llttlltt/dj-library-tools/internal/query"
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

func (p *PlexProvider) Client() *plex.Client {
	return p.client
}

func (p *PlexProvider) CanTranscode() bool {
	return true
}

func (p *PlexProvider) GetRawTracks(query string) (interface{}, error) {
	return p.getRawTracksInternal(query)
}

func (p *PlexProvider) getRawTracksInternal(queryString string) ([]plex.Track, error) {
	ctx := context.Background()
	baseURL, err := p.resolveBaseURL(ctx)
	if err != nil {
		return nil, err
	}

	q := query.NewParser().Parse(queryString)
	if err := q.ValidateWithFields(query.AllowedTrackFields); err != nil {
		return nil, err
	}

	playlistIDs := []string{}
	if queryString != "" {
		if q.Root == nil {
			return nil, fmt.Errorf("query must specify a field (e.g. playlist:%q or id:%q)", queryString, queryString)
		}

		var playlistName string
		var playlistOp query.Operator
		var walkResolve func(expr query.Expression)
		walkResolve = func(expr query.Expression) {
			switch v := expr.(type) {
			case query.Comparison:
				f := strings.ToLower(v.Field)
				if f == "id" || f == "ratingkey" {
					playlistIDs = append(playlistIDs, v.Value)
				} else if f == "playlist" {
					playlistName = v.Value
					playlistOp = v.Operator
				}
			case query.Logical:
				walkResolve(v.Left)
				walkResolve(v.Right)
			}
		}
		walkResolve(q.Root)

		if len(playlistIDs) == 0 && playlistName != "" {
			plexPlaylists, err := p.client.GetPlaylists(ctx, baseURL)
			if err != nil {
				return nil, err
			}
			for _, pl := range plexPlaylists {
				match := false
				if playlistOp == query.OpExact {
					match = pl.Title == playlistName
				} else {
					match = strings.Contains(strings.ToLower(pl.Title), strings.ToLower(playlistName))
				}

				if match {
					playlistIDs = append(playlistIDs, pl.RatingKey)
				}
			}

			if len(playlistIDs) == 0 {
				return []plex.Track{}, nil
			}
		}
	}

	var plexTracks []plex.Track
	if len(playlistIDs) > 0 {
		seen := make(map[string]bool)
		for _, id := range playlistIDs {
			path := "/playlists/" + id + "/items"
			pt, err := p.client.GetPlaylistTracks(ctx, baseURL, path)
			if err != nil {
				continue
			}
			for _, t := range pt {
				if !seen[t.RatingKey] {
					plexTracks = append(plexTracks, t)
					seen[t.RatingKey] = true
				}
			}
		}
	} else {
		plexTracks, err = p.client.GetAllTracks(ctx, baseURL)
	}

	if err != nil {
		return nil, err
	}

	var tracks []plex.Track
	eval := query.NewEvaluator(q)
	for _, pt := range plexTracks {
		if eval.Matches(pt.ToNeutral()) {
			tracks = append(tracks, pt)
		}
	}
	return tracks, nil
}

func (p *PlexProvider) GetTracks(queryString string) ([]models.Track, error) {
	raw, err := p.getRawTracksInternal(queryString)
	if err != nil {
		return nil, err
	}

	var tracks []models.Track
	for _, pt := range raw {
		tracks = append(tracks, pt.ToNeutral())
	}
	return tracks, nil
}

func (p *PlexProvider) GetPlaylists(queryString string) ([]models.Node, error) {
	ctx := context.Background()
	baseURL, err := p.resolveBaseURL(ctx)
	if err != nil {
		return nil, err
	}

	q := query.NewParser().Parse(queryString)
	if err := q.ValidateWithFields(query.AllowedNodeFields); err != nil {
		return nil, err
	}

	plexPlaylists, err := p.client.GetPlaylists(ctx, baseURL)
	if err != nil {
		return nil, err
	}

	var results []models.Node
	eval := query.NewEvaluator(q)

	for _, pl := range plexPlaylists {
		n := pl.ToNeutralNode()
		if eval.MatchesNode(n) {
			results = append(results, n)
		}
	}

	return results, nil
}

func (p *PlexProvider) CreateNode(parent models.Node, name string, nodeType int) (models.Node, error) {
	return models.Node{}, fmt.Errorf("plex does not support node creation via API")
}

func (p *PlexProvider) DeleteNode(node models.Node) error {
	return fmt.Errorf("plex does not support node deletion via API")
}

func (p *PlexProvider) RenameNode(node models.Node, newName string) error {
	return fmt.Errorf("plex does not support node renaming via API")
}

func (p *PlexProvider) MoveNode(node models.Node, targetParent models.Node) error {
	return fmt.Errorf("plex does not support node moving via API")
}

func (p *PlexProvider) resolveBaseURL(ctx context.Context) (string, error) {
	if p.host != "" {
		port := p.port
		if port == 0 {
			port = 32400
		}
		return fmt.Sprintf("http://%s:%d", p.host, port), nil
	}

	resources, err := p.client.GetResources(ctx)
	if err != nil {
		return "", err
	}

	for _, res := range resources {
		if res.Provides != "server" {
			continue
		}
		probe, err := p.client.ProbeBestConnection(res)
		if err == nil {
			return probe.BaseURL, nil
		}
	}

	return "", fmt.Errorf("could not find an active Plex server")
}
