package workflow

import (
	"fmt"

	"github.com/llttlltt/dj-library-tools/internal/config"
)

// buildGraph validates the Step dependency declarations in wf and returns a
// reverse-adjacency map: for each Step ID, the list of Step IDs that depend on
// it (i.e. declare it in their After field).
//
// Returns an error if:
//   - any After entry references a Step ID that does not exist in the Workflow
//   - the After declarations form a cycle (which would deadlock Execute)
func buildGraph(steps []config.Step) (map[string][]string, error) {
	// Index steps by ID for fast lookup.
	byID := make(map[string]struct{}, len(steps))
	for _, s := range steps {
		byID[s.ID] = struct{}{}
	}

	// Validate all After references and build reverse adjacency.
	// reverse[A] = [B, C] means B and C declared A in their After list.
	reverse := make(map[string][]string, len(steps))
	for _, s := range steps {
		for _, dep := range s.After {
			if _, ok := byID[dep]; !ok {
				return nil, fmt.Errorf("step %q references unknown dependency %q", s.ID, dep)
			}
			reverse[dep] = append(reverse[dep], s.ID)
		}
	}

	// Cycle detection via DFS on the forward dependency graph.
	// colour: 0 = white (unvisited), 1 = grey (in stack), 2 = black (done).
	colour := make(map[string]int, len(steps))
	var visit func(id string) error
	visit = func(id string) error {
		switch colour[id] {
		case 1:
			return fmt.Errorf("cycle detected involving step %q", id)
		case 2:
			return nil
		}
		colour[id] = 1
		for _, dep := range afterOf(steps, id) {
			if err := visit(dep); err != nil {
				return err
			}
		}
		colour[id] = 2
		return nil
	}
	for _, s := range steps {
		if err := visit(s.ID); err != nil {
			return nil, err
		}
	}

	return reverse, nil
}

// afterOf returns the After slice for the step with the given ID.
func afterOf(steps []config.Step, id string) []string {
	for _, s := range steps {
		if s.ID == id {
			return s.After
		}
	}
	return nil
}
