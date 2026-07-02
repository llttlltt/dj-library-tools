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

	// Register providers so orchestrator calls succeed.
	_ "github.com/llttlltt/dj-library-tools/internal/providers/m3u"
	_ "github.com/llttlltt/dj-library-tools/internal/providers/mock"
)

// setupConnections writes a mock Connection file to a temp config dir and returns its ID.
func setupConnections(t *testing.T) (connectionID string) {
	t.Helper()
	tmp := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", tmp)

	conn := config.Connection{
		ID:       config.NewConnectionID(),
		Name:     "Test Mock",
		Provider: "mock",
		Config:   map[string]string{},
	}
	require.NoError(t, config.SaveConnection(conn))
	return conn.ID
}

func newTestEngine() *Engine {
	orch := orchestrator.New(nil, orchestrator.Options{})
	return New(orch)
}

// ── single-step Workflow ───────────────────────────────────────────────────

func TestEngine_SingleStep_Success(t *testing.T) {
	connID := setupConnections(t)
	engine := newTestEngine()

	wf := config.Workflow{
		ID:   config.NewWorkflowID(),
		Name: "Single",
		Steps: []config.Step{
			{
				ID:   config.NewStepID(),
				Kind: "sync",
				Source: config.Endpoint{
					ConnectionID: connID, Resource: "tracks",
				},
				Targets: []config.Endpoint{
					{ConnectionID: connID, Resource: "tracks"},
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
	connID := setupConnections(t)
	engine := newTestEngine()

	mkStep := func() config.Step {
		return config.Step{
			ID:   config.NewStepID(),
			Kind: "sync",
			Source: config.Endpoint{
				ConnectionID: connID, Resource: "tracks",
			},
			Targets: []config.Endpoint{
				{ConnectionID: connID, Resource: "tracks"},
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
	connID := setupConnections(t)
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
					ConnectionID: connID, Resource: "tracks",
				},
				Targets: []config.Endpoint{
					{ConnectionID: connID, Resource: "tracks"},
				},
			},
			{
				ID:    secondID,
				Kind:  "sync",
				After: []string{firstID},
				Source: config.Endpoint{
					ConnectionID: connID, Resource: "tracks",
				},
				Targets: []config.Endpoint{
					{ConnectionID: connID, Resource: "tracks"},
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
	connID := setupConnections(t)
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
					ConnectionID: connID, Resource: "tracks",
				},
			},
			// Depends on failID — should be blocked.
			{
				ID:    dependentID,
				Kind:  "sync",
				After: []string{failID},
				Source: config.Endpoint{
					ConnectionID: connID, Resource: "tracks",
				},
				Targets: []config.Endpoint{
					{ConnectionID: connID, Resource: "tracks"},
				},
			},
			// No dependency — should succeed regardless.
			{
				ID:   unrelatedID,
				Kind: "sync",
				Source: config.Endpoint{
					ConnectionID: connID, Resource: "tracks",
				},
				Targets: []config.Endpoint{
					{ConnectionID: connID, Resource: "tracks"},
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
	connID := setupConnections(t)
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
					ConnectionID: connID, Resource: "tracks",
				},
			},
			{
				ID:    bID,
				Kind:  "sync",
				After: []string{aID}, // B waits for A — cycle!
				Source: config.Endpoint{
					ConnectionID: connID, Resource: "tracks",
				},
			},
		},
	}

	_, err := engine.Execute(context.Background(), wf, false)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "cycle")
}

func TestEngine_AddAndRemove_Success(t *testing.T) {
	connID := setupConnections(t)
	engine := newTestEngine()

	wf := config.Workflow{
		ID:   config.NewWorkflowID(),
		Name: "AddRemove",
		Steps: []config.Step{
			{
				ID:   "add-step",
				Kind: "add",
				Source: config.Endpoint{
					ConnectionID: connID, Resource: "tracks",
				},
				Targets: []config.Endpoint{
					{ConnectionID: connID, Resource: "playlists", Query: "New Playlist"},
				},
			},
			{
				ID:    "remove-step",
				Kind:  "remove",
				After: []string{"add-step"},
				Source: config.Endpoint{
					ConnectionID: connID, Resource: "tracks",
				},
				Targets: []config.Endpoint{
					{ConnectionID: connID, Resource: "playlists", Query: "New Playlist"},
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

func TestEngine_M3UExportAndReuse(t *testing.T) {
	connID := setupConnections(t)
	engine := newTestEngine()

	tmpFile := filepath.Join(t.TempDir(), "export.m3u")

	wf := config.Workflow{
		ID:   config.NewWorkflowID(),
		Name: "M3UExport",
		Steps: []config.Step{
			{
				ID:   "export",
				Kind: "m3u_export",
				Source: config.Endpoint{
					ConnectionID: connID, Resource: "tracks",
				},
				Options: map[string]interface{}{
					"path": tmpFile,
				},
			},
			{
				ID:    "use-exported",
				Kind:  "sync",
				After: []string{"export"},
				Source: config.Endpoint{
					ConnectionID: "m3u", Resource: tmpFile,
				},
				Targets: []config.Endpoint{
					{ConnectionID: connID, Resource: "playlists", Query: "From M3U"},
				},
			},
		},
	}

	// Must use Apply=true to actually write the file
	res, err := engine.Execute(context.Background(), wf, true)
	require.NoError(t, err)
	for _, sr := range res.Steps {
		if sr.Status == StatusFailed {
			t.Errorf("Step %s failed: %s", sr.StepID, sr.Error)
		}
	}
	require.Len(t, res.Steps, 2)
	assert.Equal(t, StatusSuccess, res.Steps[0].Status)
	assert.Equal(t, StatusSuccess, res.Steps[1].Status)

	// Verify file exists
	_, err = os.Stat(tmpFile)
	assert.NoError(t, err)
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
	_ = setupConnections(t)
	engine := newTestEngine()

	wf := config.Workflow{ID: config.NewWorkflowID(), Name: "Empty"}
	res, err := engine.Execute(context.Background(), wf, false)
	require.NoError(t, err)
	assert.Empty(t, res.Steps)
}

// Ensure the test helper actually creates the connections directory.
func TestSetupConnections_CreatesFile(t *testing.T) {
	connID := setupConnections(t)
	dir, err := config.GetConnectionsDir()
	require.NoError(t, err)
	_, err = os.Stat(filepath.Join(dir, connID+".json"))
	require.NoError(t, err)
}
