package orchestrator

import (
	"context"
	"testing"

	provider "github.com/llttlltt/dj-library-tools/internal/providers"
	_ "github.com/llttlltt/dj-library-tools/internal/providers/mock"
	"github.com/stretchr/testify/assert"
)

type MockFeedback struct {
	Previews  []string
	Successes []string
	Warnings  []string
	Progress  []struct{ Done, Total int }
}

func (f *MockFeedback) OnPreview(msg string) { f.Previews = append(f.Previews, msg) }
func (f *MockFeedback) OnSuccess(msg string) { f.Successes = append(f.Successes, msg) }
func (f *MockFeedback) OnWarning(msg string) { f.Warnings = append(f.Warnings, msg) }
func (f *MockFeedback) OnStatus(msg string)  {}
func (f *MockFeedback) OnProgress(done, total int) {
	f.Progress = append(f.Progress, struct{ Done, Total int }{done, total})
}
func (f *MockFeedback) OnTable(headers []string, rows [][]string) {}

func TestOrchestrator_List(t *testing.T) {
	fb := &MockFeedback{}
	o := New(fb, Options{})

	ctx := context.Background()
	// Use the registered mock provider
	res, err := o.List(ctx, "mock/tracks", "", RunOptions{}, "")
	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Len(t, res.Tracks, 1)
	assert.Equal(t, "Mock Track", res.Tracks[0].Title)
}

func TestOrchestrator_Edit_Feedback(t *testing.T) {
	fb := &MockFeedback{}
	o := New(fb, Options{})

	ctx := context.Background()
	changes := map[string]string{"comment": "test"}

	// Test Dry-run (Apply: false)
	_, err := o.Edit(ctx, "mock/tracks", "", RunOptions{Apply: false}, changes)
	assert.NoError(t, err)

	// GatedProvider (which Edit uses via resolver) should emit a preview
	assert.Len(t, fb.Previews, 1)
	assert.Contains(t, fb.Previews[0], "update tracks matching")
}

func TestOrchestrator_Stats(t *testing.T) {
	fb := &MockFeedback{}
	o := New(fb, Options{})

	ctx := context.Background()
	res, err := o.Stats(ctx, "mock/tracks", "", RunOptions{})
	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, 1, res.Count)
	assert.Equal(t, 0.0, res.AvgBPM) // Mock track has 0 BPM
}

func TestOrchestrator_List_Sorted(t *testing.T) {
	fb := &MockFeedback{}
	o := New(fb, Options{})

	ctx := context.Background()
	res, err := o.List(ctx, "mock/tracks", "", RunOptions{}, "artist")
	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, []string{"Artist", "Title"}, res.DefaultColumns)

	// Test invalid sort
	res2, err := o.List(ctx, "mock/tracks", "", RunOptions{}, "invalid")
	assert.Error(t, err)
	assert.Nil(t, res2)
	assert.Contains(t, err.Error(), "invalid sort field")
}

func TestOrchestrator_Sync_Feedback(t *testing.T) {
	fb := &MockFeedback{}
	o := New(fb, Options{})

	ctx := context.Background()
	err := o.Sync(ctx, "mock/tracks", "mock/tracks", "", RunOptions{Apply: true}, SyncOptions{})
	assert.NoError(t, err)
}

func TestOrchestrator_Sync_MultiGroup(t *testing.T) {
	fb := &MockFeedback{}
	o := New(fb, Options{})

	ctx := context.Background()
	// Target query "Multi" matches two playlists in our updated mock provider
	err := o.Sync(ctx, "mock/tracks", "mock/playlists name:Multi", "", RunOptions{Apply: true}, SyncOptions{})
	assert.NoError(t, err)
}

func TestOrchestrator_Fix_Feedback(t *testing.T) {
	fb := &MockFeedback{}
	o := New(fb, Options{})

	ctx := context.Background()
	_, err := o.Fix(ctx, "mock/tracks", "", RunOptions{Apply: true}, FixOptions{
		Actions: map[provider.FixType][]string{
			provider.FixDuplicates: {"tracks"},
		},
	})
	assert.NoError(t, err)
}
