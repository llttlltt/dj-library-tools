package query

import (
	"strings"

	"github.com/llttlltt/dj-library-tools/internal/core/models"
)

// Allowed Fields (Derived from models at startup)
var AllowedTrackFields []string
var AllowedGroupFields []string

func init() {
	for k := range models.TrackFields {
		AllowedTrackFields = append(AllowedTrackFields, k)
	}
	for k := range models.GroupFields {
		AllowedGroupFields = append(AllowedGroupFields, k)
	}
}

type CustomMatcher interface {
	CustomMatch(track models.Track, field string, op Operator, value string) bool
}

type Evaluator struct {
	Query     Query
	Matcher   CustomMatcher
	pathCache map[string]string // key: "trackID:path"
}

func NewEvaluator(q Query) *Evaluator {
	return &Evaluator{
		Query:     q,
		pathCache: make(map[string]string),
	}
}

func NewEvaluatorWithMatcher(q Query, m CustomMatcher) *Evaluator {
	return &Evaluator{
		Query:     q,
		Matcher:   m,
		pathCache: make(map[string]string),
	}
}

func (e *Evaluator) Matches(track models.Track) bool {
	if e.Query.Root == nil {
		return true
	}
	return e.eval(e.Query.Root, track)
}

func (e *Evaluator) MatchesGroup(group models.ResourceGroup) bool {
	if e.Query.Root == nil {
		return true
	}
	return e.evalGroup(e.Query.Root, group)
}

func (e *Evaluator) eval(expr Expression, track models.Track) bool {
	switch v := expr.(type) {
	case Comparison:
		return e.matchComparison(track, v)
	case Logical:
		if v.Op == "AND" {
			return e.eval(v.Left, track) && e.eval(v.Right, track)
		}
		return e.eval(v.Left, track) || e.eval(v.Right, track)
	case Not:
		return !e.eval(v.Expr, track)
	}
	return false
}

func (e *Evaluator) evalGroup(expr Expression, group models.ResourceGroup) bool {
	switch v := expr.(type) {
	case Comparison:
		return e.matchGroupComparison(group, v)
	case Logical:
		if v.Op == "AND" {
			return e.evalGroup(v.Left, group) && e.evalGroup(v.Right, group)
		}
		return e.evalGroup(v.Left, group) || e.evalGroup(v.Right, group)
	case Not:
		return !e.evalGroup(v.Expr, group)
	}
	return false
}

func (e *Evaluator) matchComparison(track models.Track, c Comparison) bool {
	field := strings.ToLower(c.Field)
	targetValue := ResolveValue(c.Field, c.Value)

	// Path-based resolution
	if isPath(c.Field) {
		cacheKey := track.ID + ":" + c.Field
		val, ok := e.pathCache[cacheKey]
		if !ok {
			var found bool
			val, found = ResolvePath(track, c.Field)
			if found {
				e.pathCache[cacheKey] = val
			} else {
				return false
			}
		}
		return Compare(c.Field, val, targetValue, c.Operator)
	}

	return Compare(field, track.Value(field), targetValue, c.Operator)
}

func (e *Evaluator) matchGroupComparison(group models.ResourceGroup, c Comparison) bool {
	field := strings.ToLower(c.Field)
	return Compare(field, group.Value(field), c.Value, c.Operator)
}

// Helpers

func isPath(field string) bool {
	return strings.ContainsAny(field, "./-")
}
