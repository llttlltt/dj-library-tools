package mock

import (
	"context"
	"github.com/llttlltt/dj-library-tools/internal/core/models"
	"github.com/llttlltt/dj-library-tools/internal/providers"
	"github.com/llttlltt/dj-library-tools/internal/providers/factory"
)

func init() {
	factory.Register("mock", func(opts factory.ProviderOptions) (provider.Provider, error) {
		return &MockProvider{}, nil
	})
}

type MockProvider struct{}

func (f *MockProvider) Name() string                   { return "mock" }
func (f *MockProvider) Tracks() provider.TrackService  { return &mockTrackService{} }
func (f *MockProvider) Groups() provider.GroupService  { return &mockGroupService{} }
func (f *MockProvider) System() provider.SystemService { return &mockSystemService{} }

type mockTrackService struct{}

func (s *mockTrackService) List(_ context.Context, _ provider.ExecutionContext, _ string) ([]models.Track, error) {
	return []models.Track{{ID: "1", Title: "Mock Track"}}, nil
}
func (s *mockTrackService) Update(_ context.Context, _ provider.ExecutionContext, _ string, _ map[string]string) (int, error) {
	return 0, nil
}
func (s *mockTrackService) UpdateBatch(_ context.Context, _ provider.ExecutionContext, _ []models.MetadataMatch, _ []string) error {
	return nil
}
func (s *mockTrackService) Delete(_ context.Context, _ provider.ExecutionContext, _ string) (int, error) {
	return 0, nil
}
func (s *mockTrackService) Groups() provider.TrackGroupService { return s }
func (s *mockTrackService) Add(_ context.Context, _ provider.ExecutionContext, _ []models.Track, _ models.ResourceGroup) (int, error) {
	return 1, nil
}
func (s *mockTrackService) Remove(_ context.Context, _ provider.ExecutionContext, _ []models.Track, _ models.ResourceGroup) (int, error) {
	return 1, nil
}
func (s *mockTrackService) Move(_ context.Context, _ provider.ExecutionContext, _ []models.Track, _ models.ResourceGroup, _ models.ResourceGroup) (int, error) {
	return 1, nil
}
func (s *mockTrackService) Sort(_ context.Context, _ provider.ExecutionContext, _ []models.Track, _ string) {
}

type mockGroupService struct{}

func (s *mockGroupService) List(_ context.Context, _ provider.ExecutionContext, _ string) ([]models.ResourceGroup, error) {
	return nil, nil
}
func (s *mockGroupService) Create(_ context.Context, _ provider.ExecutionContext, _ models.ResourceGroup, name string, _ models.GroupKind, _ int) (models.ResourceGroup, error) {
	return models.ResourceGroup{Name: name}, nil
}
func (s *mockGroupService) Update(_ context.Context, _ provider.ExecutionContext, _ models.ResourceGroup, _ string, _ *models.ResourceGroup) error {
	return nil
}
func (s *mockGroupService) Delete(_ context.Context, _ provider.ExecutionContext, _ models.ResourceGroup) error {
	return nil
}
func (s *mockGroupService) Sort(_ context.Context, _ provider.ExecutionContext, _ []models.ResourceGroup, _ string) {
}

type mockSystemService struct{}

func (s *mockSystemService) Capabilities() provider.ProviderCapabilities {
	return provider.ProviderCapabilities{CanUpdateMetadata: true}
}
func (s *mockSystemService) Containment() provider.ContainmentPolicy {
	return provider.ContainmentPolicy{}
}
func (s *mockSystemService) MetadataCapabilities() []string { return nil }
func (s *mockSystemService) TableHeaders() []string         { return []string{"Artist", "Title"} }
func (s *mockSystemService) SupportedResources() []string   { return []string{"tracks"} }
func (s *mockSystemService) Save(_ context.Context, _ provider.ExecutionContext, _ string) error {
	return nil
}
func (s *mockSystemService) Fix(_ context.Context, _ provider.ExecutionContext, _ provider.Selection, _ provider.FixOptions) (int, error) {
	return 0, nil
}
func (s *mockSystemService) Sync(_ context.Context, _ provider.ExecutionContext, _ []models.Track, _ string, _ provider.SyncOptions) error {
	return nil
}
func (s *mockSystemService) Identify(_ string, _ models.GroupKind) string { return "" }
