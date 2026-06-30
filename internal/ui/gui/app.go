// Package gui exposes the Wails application bindings. App is the only type
// the frontend imports; all methods on App are callable as Go bindings.
package gui

import (
	"context"
	"fmt"

	"github.com/llttlltt/dj-library-tools/internal/config"
	"github.com/llttlltt/dj-library-tools/internal/services/orchestrator"
	"github.com/llttlltt/dj-library-tools/internal/services/workflow"

	// models needed for TrackRow population
	"github.com/llttlltt/dj-library-tools/internal/core/models"

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

// ── Diff ──────────────────────────────────────────────────────────────────────

// TrackRow is a lightweight track summary used in StepDiff lists.
type TrackRow struct {
	ID     string `json:"id"`
	Title  string `json:"title"`
	Artist string `json:"artist"`
	BPM    string `json:"bpm"`
}

// StepDiff holds per-target diff data for one Step.
type StepDiff struct {
	StepID     string     `json:"step_id"`
	TargetName string     `json:"target_name"`
	Added      []TrackRow `json:"added"`
	Removed    []TrackRow `json:"removed"`
	Unchanged  []TrackRow `json:"unchanged"`
}

// GetWorkflowDiff returns per-target, track-level diff data for every sync
// Step in the Workflow without applying any changes.
func (a *App) GetWorkflowDiff(id string) ([]StepDiff, error) {
	wf, err := a.GetWorkflow(id)
	if err != nil {
		return nil, err
	}

	runOpts := orchestrator.RunOptions{Apply: false}
	var out []StepDiff

	for _, step := range wf.Steps {
		if step.Kind != "sync" {
			continue
		}

		src, err := sourceProviderLoc(step.Source)
		if err != nil {
			return nil, fmt.Errorf("step %s source: %w", step.ID, err)
		}

		for _, tgt := range step.Targets {
			tgtLoc, err := sourceProviderLoc(tgt)
			if err != nil {
				return nil, fmt.Errorf("step %s target: %w", step.ID, err)
			}
			if tgt.Query != "" {
				tgtLoc = tgtLoc + " " + tgt.Query
			}

			diff, err := a.orch.GetSyncDiff(a.ctx, src, tgtLoc, step.Source.Query, runOpts, false)
			if err != nil {
				return nil, fmt.Errorf("step %s diff: %w", step.ID, err)
			}

			// Build the removed set for fast lookup.
			removedSet := make(map[string]bool, len(diff.RemovedIDs))
			for _, rid := range diff.RemovedIDs {
				removedSet[rid] = true
			}

			sd := StepDiff{
				StepID:     step.ID,
				TargetName: diff.TargetName,
				Added:      toTrackRows(diff.AddedIDs, diff.TrackLookup),
				Removed:    toTrackRows(diff.RemovedIDs, diff.TrackLookup),
			}
			// Unchanged = current members that are NOT being removed.
			var unchangedIDs []string
			for _, cid := range diff.CurrentIDs {
				if !removedSet[cid] {
					unchangedIDs = append(unchangedIDs, cid)
				}
			}
			sd.Unchanged = toTrackRows(unchangedIDs, diff.TrackLookup)
			out = append(out, sd)
		}
	}
	return out, nil
}

// sourceProviderLoc resolves an Endpoint to "<provider>/<resource>".
func sourceProviderLoc(ep config.Endpoint) (string, error) {
	src, err := config.FindSourceByID(ep.SourceID)
	if err != nil {
		return "", err
	}
	return src.Provider + "/" + ep.Resource, nil
}

// toTrackRows converts a slice of track IDs to TrackRow summaries using the
// lookup map populated by GetSyncDiff.
func toTrackRows(ids []string, lookup map[string]models.Track) []TrackRow {
	rows := make([]TrackRow, 0, len(ids))
	for _, id := range ids {
		t := lookup[id]
		bpm := ""
		if t.BPM > 0 {
			bpm = fmt.Sprintf("%.0f", t.BPM)
		}
		rows = append(rows, TrackRow{
			ID:     id,
			Title:  t.Title,
			Artist: t.Artist,
			BPM:    bpm,
		})
	}
	return rows
}
