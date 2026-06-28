package rb

import (
	"testing"
	"github.com/llttlltt/dj-library-tools/internal/rekordbox"
	"github.com/llttlltt/dj-library-tools/internal/query"
)

func TestQuery(t *testing.T) {
	rt := rekordbox.Track{Name: "Test"}
	tr := ToNeutralTrack(rt)
	parser := query.NewParser()
	q := parser.Parse("title:Test")
	eval := query.NewEvaluator(q)
	if !eval.Matches(tr) {
		t.Error("expected match")
	}
}
