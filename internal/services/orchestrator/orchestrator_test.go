package orchestrator

import (
	"context"
	"testing"

	"github.com/llttlltt/dj-library-tools/internal/core/models"
	"github.com/llttlltt/dj-library-tools/internal/providers"
	_ "github.com/llttlltt/dj-library-tools/internal/providers/mock"
	"github.com/stretchr/testify/assert"
)

type MockFeedback struct {
	Previews []string
	Successes []string
	Warnings []string
	Progress []struct{ Done, Total int }
}

func (f *MockFeedback) OnPreview(msg string)           { f.Previews = append(f.Previews, msg) }
func (f *MockFeedback) OnSuccess(msg string)           { f.Successes = append(f.Successes, msg) }
func (f *MockFeedback) OnWarning(msg string)           { f.Warnings = append(f.Warnings, msg) }
func (f *MockFeedback) OnProgress(done, total int)     { f.Progress = append(f.Progress, struct{ Done, Total int }{done, total}) }
func (f *MockFeedback) OnTable(headers []string, rows [][]string) {}

func TestOrchestrator_List(t *testing.T) {
	fb := &MockFeedback{}
	o := New(fb, Options{})
	
	ctx := context.Background()
	// Use the registered mock provider
	res, err := o.List(ctx, "mock/tracks", "", RunOptions{})
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

func TestOrchestrator_Models(t *testing.T) {
	_ = models.Track{}
	_ = provider.Selection{}
}
