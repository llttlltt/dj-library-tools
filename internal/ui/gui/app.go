// Package gui exposes the Wails application bindings. App is the only type
// the frontend imports; all methods on App are callable as Go bindings.
package gui

import (
	"context"
	"fmt"

	"github.com/llttlltt/dj-library-tools/internal/config"
	"github.com/llttlltt/dj-library-tools/internal/services/orchestrator"
	"github.com/llttlltt/dj-library-tools/internal/services/workflow"

	// Provider registrations.
	_ "github.com/llttlltt/dj-library-tools/internal/providers/m3u"
	_ "github.com/llttlltt/dj-library-tools/internal/providers/plex"
	_ "github.com/llttlltt/dj-library-tools/internal/providers/rekordbox"
)

// App is the Wails application object. Its exported methods become TypeScript
// bindings in the frontend. All methods return JSON-serialisable types.
type App struct {
	ctx    context.Context
	orch   *orchestrator.Orchestrator
	engine *workflow.Engine
}

// NewApp constructs the App, resolving the primary rekordbox Source at startup.
// If no Source is found the orchestrator is still created; commands that
// require a file will fail with a clear error at call time.
func NewApp() *App {
	var primaryPath string
	if src, err := config.FindFirstSource("rb"); err == nil {
		primaryPath = config.ResolveProviderOptions(*src).FilePath
	}

	orch := orchestrator.New(nil, orchestrator.Options{
		RekordboxPrimaryPath: primaryPath,
	})

	return &App{
		orch:   orch,
		engine: workflow.New(orch),
	}
}

// Startup is called by Wails when the application starts.
func (a *App) Startup(ctx context.Context) { a.ctx = ctx }

// ── Sources ───────────────────────────────────────────────────────────────────

// ListSources returns all configured Sources.
func (a *App) ListSources() ([]config.Source, error) {
	return config.LoadSources()
}

// CreateSource generates a new UUID, saves the Source file, and returns it.
func (a *App) CreateSource(name, prov string, cfg map[string]string) (config.Source, error) {
	if name == "" {
		return config.Source{}, fmt.Errorf("source name is required")
	}
	if prov == "" {
		return config.Source{}, fmt.Errorf("source provider is required")
	}
	s := config.Source{
		ID:       config.NewSourceID(),
		Name:     name,
		Provider: prov,
		Config:   cfg,
	}
	if err := config.SaveSource(s); err != nil {
		return config.Source{}, err
	}
	return s, nil
}

// DeleteSource removes the Source with the given ID.
func (a *App) DeleteSource(id string) error {
	return config.DeleteSource(id)
}

// ── Workflows ─────────────────────────────────────────────────────────────────

// ListWorkflows returns all configured Workflows.
func (a *App) ListWorkflows() ([]config.Workflow, error) {
	return config.LoadWorkflows()
}

// GetWorkflow returns the Workflow with the given ID.
func (a *App) GetWorkflow(id string) (config.Workflow, error) {
	wfs, err := config.LoadWorkflows()
	if err != nil {
		return config.Workflow{}, err
	}
	for _, w := range wfs {
		if w.ID == id {
			return w, nil
		}
	}
	return config.Workflow{}, fmt.Errorf("workflow %q not found", id)
}

// SaveWorkflow assigns a UUID if the Workflow has no ID, saves it, and returns
// the saved Workflow (with ID populated).
func (a *App) SaveWorkflow(wf config.Workflow) (config.Workflow, error) {
	if wf.ID == "" {
		wf.ID = config.NewWorkflowID()
	}
	for i := range wf.Steps {
		if wf.Steps[i].ID == "" {
			wf.Steps[i].ID = config.NewStepID()
		}
	}
	if err := config.SaveWorkflow(wf); err != nil {
		return config.Workflow{}, err
	}
	return wf, nil
}

// DeleteWorkflow removes the Workflow with the given ID.
func (a *App) DeleteWorkflow(id string) error {
	return config.DeleteWorkflow(id)
}

// ── Execution ─────────────────────────────────────────────────────────────────

// PreviewWorkflow executes the Workflow with Apply=false and returns per-Step
// results showing what would change.
func (a *App) PreviewWorkflow(id string) (workflow.WorkflowResult, error) {
	wf, err := a.GetWorkflow(id)
	if err != nil {
		return workflow.WorkflowResult{}, err
	}
	return a.engine.Execute(a.ctx, wf, false)
}

// RunWorkflow executes the Workflow with Apply=true, committing all changes.
func (a *App) RunWorkflow(id string) (workflow.WorkflowResult, error) {
	wf, err := a.GetWorkflow(id)
	if err != nil {
		return workflow.WorkflowResult{}, err
	}
	return a.engine.Execute(a.ctx, wf, true)
}
