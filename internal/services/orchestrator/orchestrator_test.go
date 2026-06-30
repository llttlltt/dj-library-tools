package orchestrator

import (
	"context"
	"testing"

	"github.com/llttlltt/dj-library-tools/internal/core/models"
	"github.com/llttlltt/dj-library-tools/internal/providers"
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
	assert.NotNil(t, o)
	
	ctx := context.Background()
	res, _ := o.List(ctx, "", "", RunOptions{})
	// ListResult should be returned but empty if selection is empty/nil
	assert.NotNil(t, res)
	assert.Empty(t, res.Tracks)
}

func TestOrchestrator_Models(t *testing.T) {
	// Proving models import is used
	_ = models.Track{}
	_ = provider.Selection{}
}
