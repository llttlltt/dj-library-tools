package engine

import (
	"github.com/llttlltt/dj-library-tools/internal/models"
	"github.com/llttlltt/dj-library-tools/internal/query"
)

// LsPlaylists returns all playlist nodes (Type=1) matching the given query string.
func (e *Engine) LsPlaylists(queryString string) ([]models.Node, error) {
	return e.lsNodes(queryString, 1)
}

// LsFolders returns all folder nodes (Type=0) matching the given query string.
func (e *Engine) LsFolders(queryString string) ([]models.Node, error) {
	return e.lsNodes(queryString, 0)
}

func (e *Engine) lsNodes(queryString string, nodeType int) ([]models.Node, error) {
	parser := query.NewParser()
	q := parser.Parse(queryString)
	if err := q.ValidateWithFields(query.AllowedNodeFields); err != nil {
		return nil, err
	}
	eval := query.NewEvaluator(q)

	var matched []models.Node
	e.collectNodes(e.Library.GetPlaylists(), eval, nodeType, &matched)
	return matched, nil
}

func (e *Engine) collectNodes(nodes []models.Node, eval *query.Evaluator, nodeType int, out *[]models.Node) {
	for _, node := range nodes {
		if node.Type == nodeType {
			if eval.MatchesNode(node) {
				*out = append(*out, node)
			}
		}
		// If node has children (rekordbox specific usually)
		// We'll need a better way to handle hierarchy in neutral nodes soon
	}
}
