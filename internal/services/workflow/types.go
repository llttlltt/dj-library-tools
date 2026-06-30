package workflow

import (
	provider "github.com/llttlltt/dj-library-tools/internal/providers"
)

// Step status values.
const (
	StatusPending = "pending"
	StatusRunning = "running"
	StatusSuccess = "success"
	StatusFailed  = "failed"
	StatusBlocked = "blocked"
)

// StepResult records the outcome of one Step execution.
type StepResult struct {
	StepID    string   `json:"step_id"`
	Status    string   `json:"status"`
	Previews  []string `json:"previews,omitempty"`
	Successes []string `json:"successes,omitempty"`
	Warnings  []string `json:"warnings,omitempty"`
	Error     string   `json:"error,omitempty"`
}

// WorkflowResult is the aggregate outcome of a full Workflow execution.
type WorkflowResult struct {
	WorkflowID string       `json:"workflow_id"`
	Steps      []StepResult `json:"steps"`
}

// GUIFeedback implements provider.Feedback by appending messages to a
// *StepResult so the engine can collect per-Step output.
type GUIFeedback struct {
	result *StepResult
}

// NewGUIFeedback returns a GUIFeedback targeting r.
func NewGUIFeedback(r *StepResult) *GUIFeedback { return &GUIFeedback{result: r} }

func (g *GUIFeedback) OnPreview(msg string)             { g.result.Previews = append(g.result.Previews, msg) }
func (g *GUIFeedback) OnSuccess(msg string)             { g.result.Successes = append(g.result.Successes, msg) }
func (g *GUIFeedback) OnWarning(msg string)             { g.result.Warnings = append(g.result.Warnings, msg) }
func (g *GUIFeedback) OnStatus(_ string)                {}
func (g *GUIFeedback) OnProgress(_, _ int)              {}
func (g *GUIFeedback) OnTable(_ []string, _ [][]string) {}

// Ensure GUIFeedback satisfies the provider.Feedback interface.
var _ provider.Feedback = (*GUIFeedback)(nil)
