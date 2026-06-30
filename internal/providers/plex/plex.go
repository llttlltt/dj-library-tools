package plex

import (
	"context"
	"fmt"

	"github.com/llttlltt/dj-library-tools/internal/services/library"
	"github.com/llttlltt/dj-library-tools/internal/core/models"
	"github.com/llttlltt/dj-library-tools/internal/providers"
	"github.com/llttlltt/dj-library-tools/internal/providers/factory"
)

func init() {
	factory.Register("plex", func(opts factory.ProviderOptions) (provider.Provider, error) {
		token := ""
		if opts.Config != nil {
			token = opts.Config.Plex.Token
		}
		if token == "" {
			token = "PLEX_TOKEN"
		}
		return NewPlexProvider(token, opts.Config.Plex.Host, opts.Config.Plex.Port), nil
	})
}

type PlexProvider struct {
	client *Client
	host   string
	port   int
}

func NewPlexProvider(token string, host string, port int) *PlexProvider {
	return &PlexProvider{
		client: NewClient(token),
		host:   host,
		port:   port,
	}
}

func (p *PlexProvider) Name() string { return "plex" }

func (p *PlexProvider) Tracks() provider.TrackService   { return &plexTrackService{p} }
func (p *PlexProvider) Groups() provider.GroupService   { return &plexGroupService{p} }
func (p *PlexProvider) System() provider.SystemService { return &plexSystemService{p} }

type plexTrackService struct{ *PlexProvider }

func (s *plexTrackService) List(ctx context.Context, ectx provider.ExecutionContext, queryString string) ([]models.Track, error) {
	baseURL, err := s.resolveBaseURL(context.Background())
	if err != nil { return nil, err }

	eng := library.NewEngine(NewLibrary(s.client, baseURL))
	return eng.Ls(queryString, nil)
}

func (s *plexTrackService) Update(ctx context.Context, ectx provider.ExecutionContext, query string, changes map[string]string) (int, error) {
	return 0, provider.ErrReadOnly
}

func (s *plexTrackService) UpdateBatch(ctx context.Context, ectx provider.ExecutionContext, matches []models.MetadataMatch, fields []string) error {
	return provider.ErrReadOnly
}

func (s *plexTrackService) Delete(ctx context.Context, ectx provider.ExecutionContext, query string) (int, error) {
	return 0, provider.ErrReadOnly
}

func (s *plexTrackService) Groups() provider.TrackGroupService { return s }
func (s *plexTrackService) Add(ctx context.Context, ectx provider.ExecutionContext, tracks []models.Track, target models.ResourceGroup) (int, error) { return 0, provider.ErrReadOnly }
func (s *plexTrackService) Remove(ctx context.Context, ectx provider.ExecutionContext, tracks []models.Track, group models.ResourceGroup) (int, error) { return 0, provider.ErrReadOnly }
func (s *plexTrackService) Move(ctx context.Context, ectx provider.ExecutionContext, tracks []models.Track, from models.ResourceGroup, to models.ResourceGroup) (int, error) { return 0, provider.ErrReadOnly }
func (s *plexTrackService) Sort(ctx context.Context, ectx provider.ExecutionContext, tracks []models.Track, field string) {}

type plexGroupService struct{ *PlexProvider }

func (s *plexGroupService) List(ctx context.Context, ectx provider.ExecutionContext, queryString string) ([]models.ResourceGroup, error) {
	baseURL, err := s.resolveBaseURL(context.Background())
	if err != nil { return nil, err }

	eng := library.NewEngine(NewLibrary(s.client, baseURL))
	return eng.LsGroups(queryString)
}

func (s *plexGroupService) Create(ctx context.Context, ectx provider.ExecutionContext, parent models.ResourceGroup, name string, gt models.GroupKind, pos int) (models.ResourceGroup, error) {
	return models.ResourceGroup{}, provider.ErrReadOnly
}

func (s *plexGroupService) Update(ctx context.Context, ectx provider.ExecutionContext, group models.ResourceGroup, newName string, newParent *models.ResourceGroup) error {
	return provider.ErrReadOnly
}

func (s *plexGroupService) Delete(ctx context.Context, ectx provider.ExecutionContext, group models.ResourceGroup) error {
	return provider.ErrReadOnly
}

func (s *plexGroupService) Sort(ctx context.Context, ectx provider.ExecutionContext, groups []models.ResourceGroup, field string) {}

type plexSystemService struct{ *PlexProvider }

func (s *plexSystemService) Capabilities() provider.ProviderCapabilities { return provider.ProviderCapabilities{} }
func (s *plexSystemService) Containment() provider.ContainmentPolicy { return provider.ContainmentPolicy{} }
func (s *plexSystemService) MetadataCapabilities() []string {
	return provider.ResolveAvailableFields(s.Capabilities())
}

func (s *plexSystemService) SupportedResources() []string { return []string{"tracks", "playlists"} }
func (s *plexSystemService) TableHeaders() []string {
	return []string{"bpm", "key", "artist", "title"}
}
func (s *plexSystemService) Save(ctx context.Context, ectx provider.ExecutionContext, path string) error { return nil }
func (s *plexSystemService) Fix(ctx context.Context, ectx provider.ExecutionContext, selection provider.Selection, options provider.FixOptions) (int, error) {
	return 0, provider.ErrReadOnly
}
func (s *plexSystemService) Sync(ctx context.Context, ectx provider.ExecutionContext, tracks []models.Track, targetQuery string, opts provider.SyncOptions) error { return provider.ErrReadOnly }
func (s *plexSystemService) Identify(name string, gt models.GroupKind) string { return "" }

func (p *PlexProvider) resolveBaseURL(ctx context.Context) (string, error) {
	if p.host != "" {
		port := p.port
		if port == 0 { port = 32400 }
		return fmt.Sprintf("http://%s:%d", p.host, port), nil
	}
	resources, err := p.client.GetResources(ctx)
	if err != nil { return "", err }
	for _, res := range resources {
		if res.Provides == "server" {
			probe, err := p.client.ProbeBestConnection(res)
			if err == nil { return probe.BaseURL, nil }
		}
	}
	return "", fmt.Errorf("could find no active Plex server")
}
