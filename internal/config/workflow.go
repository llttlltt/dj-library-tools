package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/google/uuid"
)

// Endpoint is the {source_id, resource, query} triple that identifies one side
// of a Step. Source is the source Endpoint; Targets are the target Endpoints.
type Endpoint struct {
	SourceID string `json:"source_id"`
	Resource string `json:"resource"`
	Query    string `json:"query,omitempty"`
}

// Step is an atomic operation within a Workflow — one orchestrator call. Steps
// with an empty After slice execute concurrently with other such Steps. Steps
// that declare After wait for all listed Step IDs to complete.
type Step struct {
	ID      string                 `json:"id"`
	Kind    string                 `json:"kind"`
	Source  Endpoint               `json:"source"`
	Targets []Endpoint             `json:"targets"`
	After   []string               `json:"after,omitempty"`
	Options map[string]interface{} `json:"options,omitempty"`
}

// Workflow is a named, user-defined collection of Steps.
type Workflow struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Steps []Step `json:"steps"`
}

// NewWorkflowID returns a new UUID v4 string for a Workflow.
func NewWorkflowID() string { return uuid.New().String() }

// NewStepID returns a new UUID v4 string for a Step.
func NewStepID() string { return uuid.New().String() }

// LoadWorkflows reads all *.json files from ~/.config/djlt/workflows/.
func LoadWorkflows() ([]Workflow, error) {
	dir, err := GetWorkflowsDir()
	if err != nil {
		return nil, err
	}
	entries, err := filepath.Glob(filepath.Join(dir, "*.json"))
	if err != nil {
		return nil, err
	}
	var out []Workflow
	for _, p := range entries {
		data, err := os.ReadFile(p)
		if err != nil {
			return nil, fmt.Errorf("reading workflow %s: %w", p, err)
		}
		var w Workflow
		if err := json.Unmarshal(data, &w); err != nil {
			return nil, fmt.Errorf("parsing workflow %s: %w", p, err)
		}
		out = append(out, w)
	}
	return out, nil
}

// SaveWorkflow writes w to ~/.config/djlt/workflows/<id>.json.
func SaveWorkflow(w Workflow) error {
	dir, err := GetWorkflowsDir()
	if err != nil {
		return err
	}
	data, err := json.MarshalIndent(w, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(dir, w.ID+".json"), data, 0644)
}

// DeleteWorkflow removes ~/.config/djlt/workflows/<id>.json.
func DeleteWorkflow(id string) error {
	dir, err := GetWorkflowsDir()
	if err != nil {
		return err
	}
	p := filepath.Join(dir, id+".json")
	if err := os.Remove(p); err != nil && !os.IsNotExist(err) {
		return err
	}
	return nil
}
