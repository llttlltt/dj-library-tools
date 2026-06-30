package cli

import (
	"context"
	"testing"

	"github.com/llttlltt/dj-library-tools/internal/core/models"
	"github.com/llttlltt/dj-library-tools/internal/providers"
)

type MockProvider struct{}

func (m *MockProvider) Name() string                   { return "mock" }
func (m *MockProvider) Tracks() provider.TrackService  { return &mockTrackService{} }
func (m *MockProvider) Groups() provider.GroupService  { return &mockGroupService{} }
func (m *MockProvider) System() provider.SystemService { return &mockSystemService{} }

type mockTrackService struct{}

func (s *mockTrackService) List(_ context.Context, _ provider.ExecutionContext, _ string) ([]models.Track, error) {
	return []models.Track{{ID: "1", Title: "Test"}}, nil
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
	return provider.ProviderCapabilities{}
}
func (s *mockSystemService) Containment() provider.ContainmentPolicy {
	return provider.ContainmentPolicy{}
}
func (s *mockSystemService) MetadataCapabilities() []string { return nil }
func (s *mockSystemService) TableHeaders() []string         { return []string{"BPM", "Artist", "Title"} }
func (s *mockSystemService) SupportedResources() []string   { return nil }
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

func TestResolveSelection(t *testing.T) {
	// The current ResolveSelection uses factory.NewProvider, which we can't easily mock
	// without a real registry entry or more refactoring.
	// This test file might need a separate update if we want to test CLI resolution.
}
