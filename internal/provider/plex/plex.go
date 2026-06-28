package plex

import (
	"context"
	"fmt"
	"strings"

	"github.com/llttlltt/dj-library-tools/internal/models"
	"github.com/llttlltt/dj-library-tools/internal/plex"
	"github.com/llttlltt/dj-library-tools/internal/provider"
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

func (p *PlexProvider) Capabilities() provider.ProviderCapabilities {
	return provider.ProviderCapabilities{
		CanWrite:          false,
		CanManageGroups:   false,
		CanUpdateMetadata: false,
		SupportsCues:      false,
		SupportsBeatgrids: false,
		IsFileBased:       false,
	}
}

func (p *PlexProvider) GetContainmentPolicy() provider.ContainmentPolicy {
	return provider.ContainmentPolicy{
		AllowTracksInFolders:    false,
		AllowFoldersInPlaylists: false,
		AllowNestedFolders:      false,
	}
}

func (p *PlexProvider) CustomMatch(track models.Track, field string, op query.Operator, value string) bool {
	return false
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

func (p *PlexProvider) getRawTracksInternal(_ provider.ExecutionContext, queryString string) ([]plex.Track, error) {
	goCtx := context.Background()
	baseURL, err := p.resolveBaseURL(goCtx)
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
				switch f {
				case "id", "ratingkey":
					playlistIDs = append(playlistIDs, v.Value)
				case "playlist":
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
			plexPlaylists, err := p.client.GetPlaylists(goCtx, baseURL)
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
			pt, err := p.client.GetPlaylistTracks(goCtx, baseURL, path)
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
		plexTracks, err = p.client.GetAllTracks(goCtx, baseURL)
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

func (p *PlexProvider) GetTracks(ctx provider.ExecutionContext, queryString string) ([]models.Track, error) {
	raw, err := p.getRawTracksInternal(ctx, queryString)
	if err != nil {
		return nil, err
	}

	var tracks []models.Track
	for _, pt := range raw {
		tracks = append(tracks, pt.ToNeutral())
	}
	return tracks, nil
}

func (p *PlexProvider) GetFolders(_ provider.ExecutionContext, _ string) ([]models.ResourceGroup, error) {
	return nil, nil // Plex has no folder concept
}

func (p *PlexProvider) GetPlaylists(ctx provider.ExecutionContext, queryString string) ([]models.ResourceGroup, error) {
	goCtx := context.Background()
	baseURL, err := p.resolveBaseURL(goCtx)
	if err != nil {
		return nil, err
	}

	q := query.NewParser().Parse(queryString)
	if err := q.ValidateWithFields(query.AllowedNodeFields); err != nil {
		return nil, err
	}

	plexPlaylists, err := p.client.GetPlaylists(goCtx, baseURL)
	if err != nil {
		return nil, err
	}

	var results []models.ResourceGroup
	eval := query.NewEvaluator(q)

	for _, pl := range plexPlaylists {
		n := pl.ToNeutralNode()
		if eval.MatchesGroup(n) {
			results = append(results, n)
		}
	}

	return results, nil
}

func (p *PlexProvider) CreateGroup(_ provider.ExecutionContext, _ models.ResourceGroup, _ string, _ int, _ int) (models.ResourceGroup, error) {
	return models.ResourceGroup{}, fmt.Errorf("plex does not support group creation via API")
}

func (p *PlexProvider) DeleteGroup(_ provider.ExecutionContext, _ models.ResourceGroup) error {
	return fmt.Errorf("plex does not support group deletion via API")
}

func (p *PlexProvider) RenameGroup(_ provider.ExecutionContext, _ models.ResourceGroup, _ string) error {
	return fmt.Errorf("plex does not support group renaming via API")
}

func (p *PlexProvider) MoveGroup(_ provider.ExecutionContext, _ models.ResourceGroup, _ models.ResourceGroup) error {
	return fmt.Errorf("plex does not support group moving via API")
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

func (p *PlexProvider) Sync(_ provider.ExecutionContext, _ []models.Track, _ string, _ string, _ provider.SyncOptions) error {
	return fmt.Errorf("sync not supported for plex")
}

func (p *PlexProvider) ModifyTracks(_ provider.ExecutionContext, _ string, _ map[string]string) (int, error) {
	return 0, fmt.Errorf("plex provider does not support metadata modification")
}

func (p *PlexProvider) ValidateAddTracks(_ models.ResourceGroup) error {
	return fmt.Errorf("plex provider is read-only")
}

func (p *PlexProvider) ValidateMoveGroup(_ models.ResourceGroup, _ models.ResourceGroup) error {
	return fmt.Errorf("plex provider is read-only")
}

func (p *PlexProvider) ValidateCreateGroup(_ models.ResourceGroup, _ models.GroupType) error {
	return fmt.Errorf("plex provider is read-only")
}

func (p *PlexProvider) Save(_ provider.ExecutionContext, _ string) error {
	return nil
}

func (p *PlexProvider) SortTracks(_ provider.ExecutionContext, tracks []models.Track, field string) {}
func (p *PlexProvider) SortGroups(_ provider.ExecutionContext, groups []models.ResourceGroup, field string) {}

func (p *PlexProvider) GetResources(ctx provider.ExecutionContext, resource string, query string) ([]models.Resource, error) {
	var items []models.Resource
	switch resource {
	case "tracks":
		tracks, err := p.GetTracks(ctx, query)
		if err != nil { return nil, err }
		for _, t := range tracks { items = append(items, t) }
	case "playlists":
		groups, err := p.GetPlaylists(ctx, query)
		if err != nil { return nil, err }
		for _, g := range groups { items = append(items, g) }
	default:
		return nil, fmt.Errorf("unknown resource: %s", resource)
	}
	return items, nil
}
