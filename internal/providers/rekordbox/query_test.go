package rekordbox

import (
	"testing"
	"github.com/llttlltt/dj-library-tools/internal/core/query"
)

func TestQuery(t *testing.T) {
	rt := Track{Name: "Test"}
	tr := ToNeutralTrack(rt)
	parser := query.NewParser()
	q := parser.Parse("title:Test")
	eval := query.NewEvaluator(q)
	if !eval.Matches(tr) {
		t.Error("expected match")
	}
}
