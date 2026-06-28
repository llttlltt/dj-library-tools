package plex

import (
	"context"
	"fmt"
	"strings"

	"github.com/llttlltt/dj-library-tools/internal/models"
	"github.com/llttlltt/dj-library-tools/internal/plex"
	"github.com/llttlltt/dj-library-tools/internal/provider"
	"github.com/llttlltt/dj-library-tools/internal/provider/factory"
	"github.com/llttlltt/dj-library-tools/internal/query"
)

func init() {
	factory.Register("plex", func(opts factory.ProviderOptions) (provider.Provider, error) {
		token := ""
		if opts.Config != nil {
			token = opts.Config.PlexToken
		}
		if token == "" {
			token = "PLEX_TOKEN"
		}
		return NewPlexProvider(token, opts.Config.PlexHost, opts.Config.PlexPort), nil
	})
}

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

func (p *PlexProvider) Name() string { return "plex" }

func (p *PlexProvider) Tracks() provider.TrackService { return &plexTrackService{p} }
func (p *PlexProvider) Groups() provider.GroupService { return &plexGroupService{p} }
func (p *PlexProvider) System() provider.SystemService { return &plexSystemService{p} }

type plexTrackService struct{ *PlexProvider }

func (s *plexTrackService) List(ctx provider.ExecutionContext, queryString string) ([]models.Track, error) {
	raw, err := s.getRawTracksInternal(ctx, queryString)
	if err != nil { return nil, err }
	var tracks []models.Track
	for _, pt := range raw { tracks = append(tracks, pt.ToNeutral()) }
	return tracks, nil
}

func (s *plexTrackService) Update(ctx provider.ExecutionContext, query string, changes map[string]string) (int, error) {
	return 0, provider.ErrReadOnly
}

func (s *plexTrackService) UpdateBatch(ctx provider.ExecutionContext, matches []models.MetadataMatch, fields []string) error {
	return provider.ErrReadOnly
}

func (s *plexTrackService) Delete(ctx provider.ExecutionContext, query string) (int, error) {
	return 0, provider.ErrReadOnly
}

func (s *plexTrackService) Groups() provider.TrackGroupService { return s }

func (s *plexTrackService) Add(ctx provider.ExecutionContext, tracks []models.Track, target models.ResourceGroup) (int, error) {
	return 0, provider.ErrReadOnly
}

func (s *plexTrackService) Remove(ctx provider.ExecutionContext, tracks []models.Track, group models.ResourceGroup) (int, error) {
	return 0, provider.ErrReadOnly
}

func (s *plexTrackService) Move(ctx provider.ExecutionContext, tracks []models.Track, from models.ResourceGroup, to models.ResourceGroup) (int, error) {
	return 0, provider.ErrReadOnly
}

func (s *plexTrackService) Sort(ctx provider.ExecutionContext, tracks []models.Track, field string) {}

type plexGroupService struct{ *PlexProvider }

func (s *plexGroupService) List(ctx provider.ExecutionContext, queryString string) ([]models.ResourceGroup, error) {
	goCtx := context.Background()
	baseURL, err := s.resolveBaseURL(goCtx)
	if err != nil { return nil, err }
	q := query.NewParser().Parse(queryString)
	plexPlaylists, err := s.client.GetPlaylists(goCtx, baseURL)
	if err != nil { return nil, err }
	var results []models.ResourceGroup
	eval := query.NewEvaluator(q)
	for _, pl := range plexPlaylists {
		n := pl.ToNeutralGroup()
		if eval.MatchesGroup(n) { results = append(results, n) }
	}
	return results, nil
}

func (s *plexGroupService) Create(ctx provider.ExecutionContext, parent models.ResourceGroup, name string, gt models.GroupType, pos int) (models.ResourceGroup, error) {
	return models.ResourceGroup{}, provider.ErrReadOnly
}

func (s *plexGroupService) Update(ctx provider.ExecutionContext, group models.ResourceGroup, newName string, newParent *models.ResourceGroup) error {
	return provider.ErrReadOnly
}

func (s *plexGroupService) Delete(ctx provider.ExecutionContext, group models.ResourceGroup) error {
	return provider.ErrReadOnly
}

func (s *plexGroupService) Sort(ctx provider.ExecutionContext, groups []models.ResourceGroup, field string) {}

type plexSystemService struct{ *PlexProvider }

func (s *plexSystemService) Capabilities() provider.ProviderCapabilities { return provider.ProviderCapabilities{} }
func (s *plexSystemService) Containment() provider.ContainmentPolicy { return provider.ContainmentPolicy{} }
func (s *plexSystemService) MetadataCapabilities() []string { return []string{"rating", "plays"} }
func (s *plexSystemService) SupportedResources() []string { return []string{"tracks", "playlists"} }
func (s *plexSystemService) Save(ctx provider.ExecutionContext, path string) error { return nil }
func (s *plexSystemService) Fix(ctx provider.ExecutionContext, resource, query string) error { return provider.ErrReadOnly }
func (s *plexSystemService) Sync(ctx provider.ExecutionContext, tracks []models.Track, srcQ, tgtQ string, opts provider.SyncOptions) error { return provider.ErrReadOnly }
func (s *plexSystemService) Identify(name string, gt models.GroupType) string { return "" }

func (p *PlexProvider) resolveBaseURL(ctx context.Context) (string, error) {
	if p.host != "" {
		port := p.port
		if port == 0 { port = 32400 }
		return fmt.Sprintf("http://%s:%d", p.host, port), nil
	}
	resources, err := p.client.GetResources(ctx)
	if err != nil { return "", err }
	for _, res := range resources {
		if res.Provides != "server" { continue }
		probe, err := p.client.ProbeBestConnection(res)
		if err == nil { return probe.BaseURL, nil }
	}
	return "", fmt.Errorf("could not find an active Plex server")
}

func (p *PlexProvider) getRawTracksInternal(_ provider.ExecutionContext, queryString string) ([]plex.Track, error) {
	goCtx := context.Background()
	baseURL, err := p.resolveBaseURL(goCtx)
	if err != nil { return nil, err }
	q := query.NewParser().Parse(queryString)
	playlistIDs := []string{}
	if queryString != "" {
		var playlistName string
		var playlistOp query.Operator
		var walkResolve func(expr query.Expression)
		walkResolve = func(expr query.Expression) {
			switch v := expr.(type) {
			case query.Comparison:
				f := strings.ToLower(v.Field)
				if f == "id" || f == "ratingkey" { playlistIDs = append(playlistIDs, v.Value)
				} else if f == "playlist" { playlistName = v.Value; playlistOp = v.Operator }
			case query.Logical: walkResolve(v.Left); walkResolve(v.Right)
			}
		}
		walkResolve(q.Root)
		if len(playlistIDs) == 0 && playlistName != "" {
			plexPlaylists, err := p.client.GetPlaylists(goCtx, baseURL)
			if err != nil { return nil, err }
			for _, pl := range plexPlaylists {
				match := false
				if playlistOp == query.OpExact { match = pl.Title == playlistName
				} else { match = strings.Contains(strings.ToLower(pl.Title), strings.ToLower(playlistName)) }
				if match { playlistIDs = append(playlistIDs, pl.RatingKey) }
			}
		}
	}
	var plexTracks []plex.Track
	if len(playlistIDs) > 0 {
		seen := make(map[string]bool)
		for _, id := range playlistIDs {
			pt, err := p.client.GetPlaylistTracks(goCtx, baseURL, "/playlists/"+id+"/items")
			if err != nil { continue }
			for _, t := range pt { if !seen[t.RatingKey] { plexTracks = append(plexTracks, t); seen[t.RatingKey] = true } }
		}
	} else {
		plexTracks, err = p.client.GetAllTracks(goCtx, baseURL)
	}
	if err != nil { return nil, err }
	var tracks []plex.Track
	eval := query.NewEvaluator(q)
	for _, pt := range plexTracks { if eval.Matches(pt.ToNeutral()) { tracks = append(tracks, pt) } }
	return tracks, nil
}
