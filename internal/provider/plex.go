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

func (p *PlexProvider) Client() *plex.Client {
	return p.client
}

func (p *PlexProvider) GetTracks(queryString string) ([]rekordbox.Track, error) {
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

		// 1. Resolve Playlist Contexts
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
			
			// If we specifically asked for a playlist but found none, return empty
			if len(playlistIDs) == 0 {
				return []rekordbox.Track{}, nil
			}
		}
	}

	var plexTracks []plex.Track
	if len(playlistIDs) > 0 {
		// Aggregate tracks from all matching playlists
		seen := make(map[string]bool)
		for _, id := range playlistIDs {
			path := "/playlists/" + id + "/items"
			pt, err := p.client.GetPlaylistTracks(ctx, baseURL, path)
			if err != nil {
				continue // Skip failing playlists
			}
			for _, t := range pt {
				if !seen[t.RatingKey] {
					plexTracks = append(plexTracks, t)
					seen[t.RatingKey] = true
				}
			}
		}
	} else {
		// Global search
		plexTracks, err = p.client.GetAllTracks(ctx, baseURL)
	}

	if err != nil {
		return nil, err
	}

	var tracks []rekordbox.Track
	eval := query.NewEvaluator(q)

	for _, pt := range plexTracks {
		// Map Plex Track to Rekordbox Track for the Evaluator
		t := rekordbox.Track{
			TrackID:    0, // Plex doesn't have an integer TrackID in the same way
			Name:       pt.Title,
			Artist:     pt.Artist,
			Album:      pt.Album,
			Tonality:   pt.KeyTag,
			AverageBpm: fmt.Sprintf("%.2f", pt.BPM),
		}
		if pt.BPM > 0 {
			t.Tempo = []rekordbox.Tempo{{Bpm: fmt.Sprintf("%.2f", pt.BPM)}}
		}
		if len(pt.Media) > 0 && len(pt.Media[0].Part) > 0 {
			t.Location = pt.Media[0].Part[0].File
		}

		// Use the standard Evaluator!
		if eval.Matches(t) {
			tracks = append(tracks, t)
		}
	}

	return tracks, nil
}

func (p *PlexProvider) GetRawTracks(queryString string) (interface{}, error) {
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

		// 1. Resolve Playlist Contexts
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
		// Evaluator needs a rekordbox.Track, so we map just for evaluation
		t := rekordbox.Track{
			Name:       pt.Title,
			Artist:     pt.Artist,
			Album:      pt.Album,
			Tonality:   pt.KeyTag,
			AverageBpm: fmt.Sprintf("%.2f", pt.BPM),
		}
		if eval.Matches(t) {
			tracks = append(tracks, pt)
		}
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
	if err := q.ValidateWithFields(query.AllowedNodeFields); err != nil {
		return nil, err
	}

	plexPlaylists, err := p.client.GetPlaylists(ctx, baseURL)
	if err != nil {
		return nil, err
	}

	var results []NodeResult
	eval := query.NewEvaluator(q)

	for _, pl := range plexPlaylists {
		// Mock a rekordbox.Node for the Evaluator
		n := rekordbox.Node{
			Name:    pl.Title,
			Type:    1,
			Entries: rekordbox.PtrInt32(int32(pl.LeafCount)),
		}

		if eval.MatchesNode(n, "") {
			results = append(results, NodeResult{
				Name:    pl.Title,
				Entries: pl.LeafCount,
				Raw:     pl,
			})
		}
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
