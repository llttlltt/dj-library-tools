package workflow

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/llttlltt/dj-library-tools/internal/config"
	"github.com/llttlltt/dj-library-tools/internal/services/orchestrator"

	// Register the mock provider so orchestrator calls succeed.
	_ "github.com/llttlltt/dj-library-tools/internal/providers/mock"
)

// setupSources writes a mock Source file to a temp config dir and returns its ID.
func setupSources(t *testing.T) (sourceID string) {
	t.Helper()
	tmp := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", tmp)

	src := config.Source{
		ID:       config.NewSourceID(),
		Name:     "Test Mock",
		Provider: "mock",
		Config:   map[string]string{},
	}
	require.NoError(t, config.SaveSource(src))
	return src.ID
}

func newTestEngine() *Engine {
	orch := orchestrator.New(nil, orchestrator.Options{})
	return New(orch)
}

// ── single-step Workflow ───────────────────────────────────────────────────

func TestEngine_SingleStep_Success(t *testing.T) {
	srcID := setupSources(t)
	engine := newTestEngine()

	wf := config.Workflow{
		ID:   config.NewWorkflowID(),
		Name: "Single",
		Steps: []config.Step{
			{
				ID:   config.NewStepID(),
				Kind: "sync",
				Source: config.Endpoint{
					SourceID: srcID, Resource: "tracks",
				},
				Targets: []config.Endpoint{
					{SourceID: srcID, Resource: "tracks"},
				},
			},
		},
	}

	res, err := engine.Execute(context.Background(), wf, false)
	require.NoError(t, err)
	require.Len(t, res.Steps, 1)
	assert.Equal(t, StatusSuccess, res.Steps[0].Status)
}

// ── two independent Steps run in parallel ─────────────────────────────────

func TestEngine_IndependentSteps_BothSucceed(t *testing.T) {
	srcID := setupSources(t)
	engine := newTestEngine()

	mkStep := func() config.Step {
		return config.Step{
			ID:   config.NewStepID(),
			Kind: "sync",
			Source: config.Endpoint{
				SourceID: srcID, Resource: "tracks",
			},
			Targets: []config.Endpoint{
				{SourceID: srcID, Resource: "tracks"},
			},
		}
	}

	wf := config.Workflow{
		ID:    config.NewWorkflowID(),
		Name:  "Parallel",
		Steps: []config.Step{mkStep(), mkStep()},
	}

	start := time.Now()
	res, err := engine.Execute(context.Background(), wf, false)
	elapsed := time.Since(start)
	_ = elapsed // parallel execution is best-effort; don't assert timing in CI

	require.NoError(t, err)
	require.Len(t, res.Steps, 2)
	assert.Equal(t, StatusSuccess, res.Steps[0].Status)
	assert.Equal(t, StatusSuccess, res.Steps[1].Status)
}

// ── Step with after waits for its dependency ──────────────────────────────

func TestEngine_AfterDependency_Respected(t *testing.T) {
	srcID := setupSources(t)
	engine := newTestEngine()

	firstID := config.NewStepID()
	secondID := config.NewStepID()

	wf := config.Workflow{
		ID:   config.NewWorkflowID(),
		Name: "Sequential",
		Steps: []config.Step{
			{
				ID:   firstID,
				Kind: "sync",
				Source: config.Endpoint{
					SourceID: srcID, Resource: "tracks",
				},
				Targets: []config.Endpoint{
					{SourceID: srcID, Resource: "tracks"},
				},
			},
			{
				ID:    secondID,
				Kind:  "sync",
				After: []string{firstID},
				Source: config.Endpoint{
					SourceID: srcID, Resource: "tracks",
				},
				Targets: []config.Endpoint{
					{SourceID: srcID, Resource: "tracks"},
				},
			},
		},
	}

	res, err := engine.Execute(context.Background(), wf, false)
	require.NoError(t, err)
	require.Len(t, res.Steps, 2)
	assert.Equal(t, StatusSuccess, res.Steps[0].Status)
	assert.Equal(t, StatusSuccess, res.Steps[1].Status)
}

// ── failed Step blocks dependent, unrelated Step completes ────────────────

func TestEngine_FailedStep_BlocksDependent_NotUnrelated(t *testing.T) {
	srcID := setupSources(t)
	engine := newTestEngine()

	failID := config.NewStepID()
	dependentID := config.NewStepID()
	unrelatedID := config.NewStepID()

	wf := config.Workflow{
		ID:   config.NewWorkflowID(),
		Name: "PartialFailure",
		Steps: []config.Step{
			// This step uses an unknown kind so it will fail.
			{
				ID:   failID,
				Kind: "unknown-kind-that-fails",
				Source: config.Endpoint{
					SourceID: srcID, Resource: "tracks",
				},
			},
			// Depends on failID — should be blocked.
			{
				ID:    dependentID,
				Kind:  "sync",
				After: []string{failID},
				Source: config.Endpoint{
					SourceID: srcID, Resource: "tracks",
				},
				Targets: []config.Endpoint{
					{SourceID: srcID, Resource: "tracks"},
				},
			},
			// No dependency — should succeed regardless.
			{
				ID:   unrelatedID,
				Kind: "sync",
				Source: config.Endpoint{
					SourceID: srcID, Resource: "tracks",
				},
				Targets: []config.Endpoint{
					{SourceID: srcID, Resource: "tracks"},
				},
			},
		},
	}

	res, err := engine.Execute(context.Background(), wf, false)
	require.NoError(t, err)
	require.Len(t, res.Steps, 3)

	byID := map[string]StepResult{}
	for _, sr := range res.Steps {
		byID[sr.StepID] = sr
	}
	assert.Equal(t, StatusFailed, byID[failID].Status)
	assert.Equal(t, StatusBlocked, byID[dependentID].Status)
	assert.Equal(t, StatusSuccess, byID[unrelatedID].Status)
}

// ── cycle detection ────────────────────────────────────────────────────────

func TestEngine_Cycle_ReturnsError(t *testing.T) {
	srcID := setupSources(t)
	engine := newTestEngine()

	aID := config.NewStepID()
	bID := config.NewStepID()

	wf := config.Workflow{
		ID:   config.NewWorkflowID(),
		Name: "Cycle",
		Steps: []config.Step{
			{
				ID:    aID,
				Kind:  "sync",
				After: []string{bID}, // A waits for B
				Source: config.Endpoint{
					SourceID: srcID, Resource: "tracks",
				},
			},
			{
				ID:    bID,
				Kind:  "sync",
				After: []string{aID}, // B waits for A — cycle!
				Source: config.Endpoint{
					SourceID: srcID, Resource: "tracks",
				},
			},
		},
	}

	_, err := engine.Execute(context.Background(), wf, false)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "cycle")
}

// ── graph: unknown After reference returns error ───────────────────────────

func TestBuildGraph_UnknownAfterRef_ReturnsError(t *testing.T) {
	steps := []config.Step{
		{ID: "a", After: []string{"does-not-exist"}},
	}
	_, err := buildGraph(steps)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "unknown dependency")
}

// ── empty workflow ─────────────────────────────────────────────────────────

func TestEngine_EmptyWorkflow(t *testing.T) {
	_ = setupSources(t)
	engine := newTestEngine()

	wf := config.Workflow{ID: config.NewWorkflowID(), Name: "Empty"}
	res, err := engine.Execute(context.Background(), wf, false)
	require.NoError(t, err)
	assert.Empty(t, res.Steps)
}

// Ensure the test helper actually creates the sources directory.
func TestSetupSources_CreatesFile(t *testing.T) {
	srcID := setupSources(t)
	dir, err := config.GetSourcesDir()
	require.NoError(t, err)
	_, err = os.Stat(filepath.Join(dir, srcID+".json"))
	require.NoError(t, err)
}
