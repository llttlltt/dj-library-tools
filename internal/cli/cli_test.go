package cli

import (
	"bytes"
	"os"
	"strings"
	"testing"

	"github.com/llttlltt/dj-library-tools/internal/models"
	"github.com/llttlltt/dj-library-tools/internal/provider"
	"github.com/llttlltt/dj-library-tools/internal/provider/factory"
	"github.com/llttlltt/dj-library-tools/internal/query"
)

// MockProvider for CLI testing
type MockProvider struct{}

var _ provider.WritableProvider = (*MockProvider)(nil)

func (m *MockProvider) Name() string { return "mock" }
func (m *MockProvider) Capabilities() provider.ProviderCapabilities {
	return provider.ProviderCapabilities{CanWrite: true, CanManageGroups: true}
}
func (m *MockProvider) GetContainmentPolicy() provider.ContainmentPolicy { return provider.ContainmentPolicy{} }
func (m *MockProvider) CustomMatch(_ models.Track, _ string, _ query.Operator, _ string) bool { return false }
func (m *MockProvider) CanTranscode() bool { return true }
func (m *MockProvider) SupportedResources() []string { return []string{"tracks", "playlists"} }
func (m *MockProvider) GetResources(_ provider.ExecutionContext, resource, query string) ([]models.Resource, error) {
	if resource == "tracks" {
		return []models.Resource{models.Track{ID: "1", Title: "Test Track", Artist: "Test Artist", Location: "file://localhost/test.mp3"}}, nil
	}
	if resource == "playlists" {
		return []models.Resource{models.ResourceGroup{Name: "Inbox", Type: models.GroupTypePlaylist, Items: 1}}, nil
	}
	return nil, nil
}
func (m *MockProvider) SortTracks(_ provider.ExecutionContext, _ []models.Track, _ string) {}
func (m *MockProvider) SortGroups(_ provider.ExecutionContext, _ []models.ResourceGroup, _ string) {}

// Writable implementation
func (m *MockProvider) AddTracks(_ provider.ExecutionContext, _ models.ResourceGroup, _ []models.Track) (int, error) { return 1, nil }
func (m *MockProvider) RemoveTracks(_ provider.ExecutionContext, _ models.ResourceGroup, _ []models.Track) (int, error) { return 1, nil }
func (m *MockProvider) CreateGroup(_ provider.ExecutionContext, _ models.ResourceGroup, name string, gt models.GroupType, _ int) (models.ResourceGroup, error) {
	return models.ResourceGroup{Name: name, Type: gt}, nil
}
func (m *MockProvider) DeleteGroup(_ provider.ExecutionContext, _ models.ResourceGroup) error { return nil }
func (m *MockProvider) RenameGroup(_ provider.ExecutionContext, _ models.ResourceGroup, _ string, _ models.GroupType) error {
	return nil
}
func (m *MockProvider) MoveGroup(_ provider.ExecutionContext, _ models.ResourceGroup, _ models.ResourceGroup) error {
	return nil
}
func (m *MockProvider) MoveTracks(_ provider.ExecutionContext, _ models.ResourceGroup, _ models.ResourceGroup, _ []models.Track) (int, error) {
	return 1, nil
}
func (m *MockProvider) Sync(_ provider.ExecutionContext, _ []models.Track, _, _ string, _ provider.SyncOptions) error {
	return nil
}
func (m *MockProvider) UpdateTracks(_ provider.ExecutionContext, _ string, _ map[string]string) (int, error) { return 0, nil }
func (m *MockProvider) ValidateAddTracks(_ models.ResourceGroup) error { return nil }
func (m *MockProvider) ValidateMoveGroup(_, _ models.ResourceGroup) error { return nil }
func (m *MockProvider) ValidateCreateGroup(_ models.ResourceGroup, _ models.GroupType) error { return nil }
func (m *MockProvider) IdentifyGroup(n string, _ models.GroupType) string { return n }
func (m *MockProvider) Save(_ provider.ExecutionContext, _ string) error { return nil }

func init() {
	factory.Register("mock", func(opts factory.ProviderOptions) (provider.Provider, error) {
		return &MockProvider{}, nil
	})
}

func resetTestState() {
	dryRun = false
	verbose = false
	jsonOutput = false
	filePath = "mock.xml"
}

func executeCommand(args ...string) (string, error) {
	resetTestState()
	root := NewRootCmd()
	buf := new(bytes.Buffer)
	root.SetOut(buf)
	root.SetErr(buf)
	root.SetArgs(args[:])

	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w
	err := root.Execute()
	w.Close()
	os.Stdout = old
	var outBuf bytes.Buffer
	outBuf.ReadFrom(r)
	return buf.String() + outBuf.String(), err
}

func TestCommandConsistencyAgostic(t *testing.T) {
	tests := []struct {
		name     string
		args     []string
		wantIn   string
		wantErr  bool
	}{
		{
			name: "list mock positional query",
			args: []string{"ls", "mock/tracks", "title:'Test Track'"},
			wantIn: "Test Track",
		},
		{
			name: "add tracks merged into sync --append",
			args: []string{"sync", "mock/tracks", "title:Test", "--to", "mock/playlists name:Inbox", "--append", "--dry-run"},
			wantIn: "Would append to",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			out, err := executeCommand(tt.args...)
			if (err != nil) != tt.wantErr {
				t.Errorf("error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantIn != "" && !strings.Contains(out, tt.wantIn) {
				t.Errorf("out = %q, want %q", out, tt.wantIn)
			}
		})
	}
}

func (m *MockProvider) MetadataCapabilities() []string { return nil }
func (m *MockProvider) UpdateMetadata(_ provider.ExecutionContext, _ []models.MetadataMatch, _ []string) error { return nil }
func (m *MockProvider) Fix(_ provider.ExecutionContext, _, _ string) error { return nil }
