package workflow

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"github.com/llttlltt/dj-library-tools/internal/config"
	provider "github.com/llttlltt/dj-library-tools/internal/providers"
	"github.com/llttlltt/dj-library-tools/internal/services/orchestrator"
)

// Engine executes Workflows against an Orchestrator.
type Engine struct {
	orch *orchestrator.Orchestrator
}

// New returns an Engine backed by orch.
func New(orch *orchestrator.Orchestrator) *Engine {
	return &Engine{orch: orch}
}

// Execute runs every Step in wf. When apply is false (Preview mode) all Steps
// run with ExecutionContext.Apply=false — reads only, safe for concurrency.
// When apply is true (Run mode) Steps targeting the same Source must execute
// sequentially to avoid concurrent writes to the same in-memory library.
//
// Concurrency model (RISK-002):
//   - Steps with no After field and whose Source does not overlap with a
//     running Step execute in their own goroutine.
//   - In Run mode (apply=true), a per-source mutex serialises Steps that share
//     a Source, preventing concurrent writes to the same in-memory provider.
//   - In Preview mode (apply=false) all Steps may run fully in parallel since
//     no mutations occur.
//
// Failure propagation (REQ-007):
//   - A Step that errors is marked "failed".
//   - Any Step whose After list contains a failed or blocked Step is marked
//     "blocked" without executing.
//   - Unrelated Steps always run to completion.
func (e *Engine) Execute(ctx context.Context, wf config.Workflow, apply bool) (WorkflowResult, error) {
	result := WorkflowResult{WorkflowID: wf.ID}

	if len(wf.Steps) == 0 {
		return result, nil
	}

	// Build and validate the dependency graph. Rejects cycles (RISK-003).
	_, err := buildGraph(wf.Steps)
	if err != nil {
		return result, fmt.Errorf("invalid workflow graph: %w", err)
	}

	// Initialise one StepResult per Step.
	stepResults := make(map[string]*StepResult, len(wf.Steps))
	for i := range wf.Steps {
		sr := &StepResult{
			StepID: wf.Steps[i].ID,
			Status: StatusPending,
		}
		stepResults[wf.Steps[i].ID] = sr
	}

	// Per-source mutex for Run mode (apply=true). In Preview mode this map is
	// populated but the mutexes are never contended (all Steps are read-only).
	var sourceMu sync.Map // map[sourceID → *sync.Mutex]
	getSourceMu := func(sourceID string) *sync.Mutex {
		v, _ := sourceMu.LoadOrStore(sourceID, &sync.Mutex{})
		return v.(*sync.Mutex)
	}

	// doneCh is closed when a Step's goroutine finishes (success or failure).
	type stepDone struct {
		id     string
		failed bool
	}
	doneCh := make(chan stepDone, len(wf.Steps))

	// Track which Steps have finished and whether they failed.
	finished := make(map[string]bool) // id → failed?
	var finishedMu sync.Mutex

	// stepReady reports whether all After dependencies have completed
	// successfully (not failed/blocked).
	stepReady := func(step config.Step) (ready bool, anyFailed bool) {
		finishedMu.Lock()
		defer finishedMu.Unlock()
		for _, dep := range step.After {
			failed, done := finished[dep]
			if !done {
				return false, false
			}
			if failed {
				return false, true
			}
		}
		return true, false
	}

	var wg sync.WaitGroup

	// runStep executes a single Step and sends its outcome to doneCh.
	runStep := func(step config.Step, sr *StepResult) {
		defer wg.Done()
		defer func() {
			failed := sr.Status == StatusFailed
			finishedMu.Lock()
			finished[step.ID] = failed
			finishedMu.Unlock()
			doneCh <- stepDone{id: step.ID, failed: failed}
		}()

		sr.Status = StatusRunning
		fb := NewGUIFeedback(sr)

		// In Run mode serialise Steps that share the same source to prevent
		// concurrent in-memory writes to the same provider instance.
		if apply {
			mu := getSourceMu(step.Source.SourceID)
			mu.Lock()
			defer mu.Unlock()
		}

		runOpts := orchestrator.RunOptions{Apply: apply}

		// Build the orchestrator with a per-step Feedback so output is captured.
		stepOrch := orchestrator.NewWithFeedback(fb, e.orch)

		if err := executeStep(ctx, stepOrch, step, runOpts, sr); err != nil {
			sr.Status = StatusFailed
			sr.Error = err.Error()
			return
		}
		sr.Status = StatusSuccess
	}

	// Scheduling loop — starts Steps as soon as their dependencies are met.
	pending := make([]config.Step, len(wf.Steps))
	copy(pending, wf.Steps)
	started := make(map[string]bool, len(wf.Steps))

	for len(started) < len(wf.Steps) {
		launched := 0
		for _, step := range pending {
			if started[step.ID] {
				continue
			}
			sr := stepResults[step.ID]
			ready, anyFailed := stepReady(step)

			if anyFailed {
				sr.Status = StatusBlocked
				finishedMu.Lock()
				finished[step.ID] = true // treat blocked as failed for dependents
				finishedMu.Unlock()
				doneCh <- stepDone{id: step.ID, failed: true}
				started[step.ID] = true
				launched++
				continue
			}

			if len(step.After) == 0 || ready {
				started[step.ID] = true
				wg.Add(1)
				go runStep(step, sr)
				launched++
			}
		}

		// If nothing launched in this pass, wait for a Step to finish before
		// retrying — avoids a busy-loop while waiting on dependencies.
		if launched == 0 {
			<-doneCh
		}
	}

	wg.Wait()

	// Collect ordered results.
	for _, step := range wf.Steps {
		result.Steps = append(result.Steps, *stepResults[step.ID])
	}
	return result, nil
}

// executeStep dispatches a single Step to the correct orchestrator method.
// sr receives preview messages directly for step kinds (sync) that compute a
// diff before applying; this mirrors the pattern used by internal/ui/cli/sync.go.
func executeStep(ctx context.Context, orch *orchestrator.Orchestrator, step config.Step, runOpts orchestrator.RunOptions, sr *StepResult) error {
	// Resolve source location string from the Endpoint's Source.
	src, err := sourceLocation(step.Source)
	if err != nil {
		return err
	}

	switch strings.ToLower(step.Kind) {
	case "sync":
		for _, tgt := range step.Targets {
			tgtLoc, err := sourceLocation(tgt)
			if err != nil {
				return err
			}
			syncOpts := orchestrator.SyncOptions{}
			if m, ok := step.Options["metadata"]; ok {
				if fields, ok := toStringSlice(m); ok {
					syncOpts.MetadataFields = fields
				}
			}
			if m, ok := step.Options["match"]; ok {
				if fields, ok := toStringSlice(m); ok {
					syncOpts.MatchFields = fields
				}
			}
			// Append target query into the location string for group resolution.
			if tgt.Query != "" {
				tgtLoc = tgtLoc + " " + tgt.Query
			}

			// Mirror the CLI sync pattern: always compute the diff first so that
			// Preview mode produces meaningful output. Only call Sync when Apply=true.
			diffs, err := orch.GetSyncDiff(ctx, src, tgtLoc, step.Source.Query, runOpts, syncOpts.AppendOnly)
			if err != nil {
				return err
			}
			for _, diff := range diffs {
				sr.Previews = append(sr.Previews, fmt.Sprintf(
					"%s — add: %d, remove: %d, final: %d",
					diff.TargetName,
					len(diff.AddedIDs),
					len(diff.RemovedIDs),
					len(diff.SourceIDs),
				))
			}

			if runOpts.Apply {
				if err := orch.Sync(ctx, src, tgtLoc, step.Source.Query, runOpts, syncOpts); err != nil {
					return err
				}
			}
		}
		return nil

	case "fix":
		actions := make(map[provider.FixType][]string)
		if opts, ok := step.Options["actions"]; ok {
			if m, ok := opts.(map[string]interface{}); ok {
				for k, v := range m {
					if fields, ok := toStringSlice(v); ok {
						actions[provider.FixType(k)] = fields
					}
				}
			}
		}
		_, err := orch.Fix(ctx, src, step.Source.Query, runOpts, orchestrator.FixOptions{Actions: actions})
		return err

	case "edit":
		changes := make(map[string]string)
		if opts, ok := step.Options["set"]; ok {
			if m, ok := opts.(map[string]interface{}); ok {
				for k, v := range m {
					if s, ok := v.(string); ok {
						changes[k] = s
					}
				}
			}
		}
		_, err := orch.Edit(ctx, src, step.Source.Query, runOpts, changes)
		return err

	default:
		return fmt.Errorf("unknown step kind %q", step.Kind)
	}
}

// toStringSlice coerces v (from JSON unmarshalled map[string]interface{}) into
// a []string. Accepts []interface{} where every element is a string.
func toStringSlice(v interface{}) ([]string, bool) {
	switch t := v.(type) {
	case []string:
		return t, true
	case []interface{}:
		ss := make([]string, 0, len(t))
		for _, vi := range t {
			s, ok := vi.(string)
			if !ok {
				return nil, false
			}
			ss = append(ss, s)
		}
		return ss, true
	}
	return nil, false
}

// sourceLocation resolves an Endpoint to a provider location string of the
// form "<source-id>/<resource>". The resolver handles UUID lookup and config hydration.
func sourceLocation(ep config.Endpoint) (string, error) {
	if ep.SourceID == "" {
		return "", fmt.Errorf("source id missing in endpoint")
	}
	return ep.SourceID + "/" + ep.Resource, nil
}
