package query

import (
	"strings"

	"github.com/llttlltt/dj-library-tools/internal/models"
)

// Global Schemas generated via reflection
var (
	TrackSchema, AllowedTrackFields = ReflectSchema(models.Track{})
	GroupSchema, AllowedGroupFields = ReflectSchema(models.ResourceGroup{})
)

type CustomMatcher interface {
	CustomMatch(track models.Track, field string, op Operator, value string) bool
}

type Evaluator struct {
	Query   Query
	Matcher CustomMatcher
}

func NewEvaluator(q Query) *Evaluator {
	return &Evaluator{Query: q}
}

func NewEvaluatorWithMatcher(q Query, m CustomMatcher) *Evaluator {
	return &Evaluator{Query: q, Matcher: m}
}

func (e *Evaluator) Matches(track models.Track) bool {
	return e.MatchesWithPlaylists(track, nil)
}

func (e *Evaluator) MatchesWithPlaylists(track models.Track, playlists []string) bool {
	if e.Query.Root == nil { return true }
	return e.eval(e.Query.Root, track, playlists)
}

func (e *Evaluator) MatchesGroup(group models.ResourceGroup) bool {
	if e.Query.Root == nil { return true }
	return e.evalGroup(e.Query.Root, group)
}

func (e *Evaluator) eval(expr Expression, track models.Track, playlists []string) bool {
	switch v := expr.(type) {
	case Comparison:
		return e.matchComparison(track, playlists, v)
	case Logical:
		if v.Op == "AND" {
			return e.eval(v.Left, track, playlists) && e.eval(v.Right, track, playlists)
		}
		return e.eval(v.Left, track, playlists) || e.eval(v.Right, track, playlists)
	case Not:
		return !e.eval(v.Expr, track, playlists)
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

func (e *Evaluator) matchComparison(track models.Track, playlists []string, c Comparison) bool {
	field := strings.ToLower(c.Field)
	targetValue := ResolveValue(c.Field, c.Value)

	// Custom delegation (e.g. Cues by color)
	if e.Matcher != nil && isCalculatedField(field) && !isNumericIntent(c) {
		return e.Matcher.CustomMatch(track, c.Field, c.Operator, c.Value)
	}

	// Membership domain
	if field == "playlists" {
		return Compare(field, strings.Join(playlists, ","), targetValue, c.Operator)
	}

	fieldValue := GetFieldValue(track, field)
	return Compare(field, fieldValue, targetValue, c.Operator)
}

func (e *Evaluator) matchGroupComparison(group models.ResourceGroup, c Comparison) bool {
	field := strings.ToLower(c.Field)
	fieldValue := GetFieldValue(group, field)
	return Compare(field, fieldValue, c.Value, c.Operator)
}

// Helpers

func isCalculatedField(field string) bool {
	return field == "hotcues" || field == "memorycues" || field == "beatgrids"
}

func isNumericIntent(c Comparison) bool {
	if c.Quoted { return false }
	return c.Operator == OpRange || c.Operator == OpGt || c.Operator == OpLt
}
