package engine

import (
	"github.com/llttlltt/dj-library-tools/internal/query"
	"github.com/llttlltt/dj-library-tools/pkg/rekordbox"
)

// NodeResult is a matched playlist or folder node along with its direct parent folder name.
// ParentFolder is empty when the node lives at the root level.
type NodeResult struct {
	Node         rekordbox.Node
	ParentFolder string
}

// LsPlaylists returns all playlist nodes (Type=1) matching the given query string.
// Supported fields: name, folder (parent folder name), entries, type.
func (e *Engine) LsPlaylists(queryString string) ([]NodeResult, error) {
	return e.lsNodes(queryString, 1)
}

// LsFolders returns all folder nodes (Type=0) matching the given query string.
// Supported fields: name, folder (parent folder name), type.
func (e *Engine) LsFolders(queryString string) ([]NodeResult, error) {
	return e.lsNodes(queryString, 0)
}

func (e *Engine) lsNodes(queryString string, nodeType int32) ([]NodeResult, error) {
	parser := query.NewParser()
	q := parser.Parse(queryString)
	if err := q.ValidateWithFields(query.AllowedNodeFields); err != nil {
		return nil, err
	}
	eval := query.NewEvaluator(q)

	var matched []NodeResult
	e.collectNodes(e.Library.GetPlaylists(), eval, nodeType, "", &matched)
	return matched, nil
}

// collectNodes walks the node tree recursively, tracking the parent folder name at
// each level so that folder: queries work correctly even when playlist names are
// not unique across different folders.
func (e *Engine) collectNodes(nodes []rekordbox.Node, eval *query.Evaluator, nodeType int32, parentFolder string, out *[]NodeResult) {
	for _, node := range nodes {
		if node.Type == nodeType {
			if eval.MatchesNode(node.ToNeutral(parentFolder)) {
				*out = append(*out, NodeResult{Node: node, ParentFolder: parentFolder})
			}
		}
		if len(node.Node) > 0 {
			e.collectNodes(node.Node, eval, nodeType, node.Name, out)
		}
	}
}
