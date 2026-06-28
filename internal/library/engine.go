package library

import (
	"github.com/llttlltt/dj-library-tools/internal/models"
	"github.com/llttlltt/dj-library-tools/internal/query"
)

// Engine performs operations on a library using queries
type Engine struct {
	Library ReadableLibrary
}

func NewEngine(lib ReadableLibrary) *Engine {
	return &Engine{
		Library: lib,
	}
}

// Ls returns all tracks that match the given query string
func (e *Engine) Ls(queryString string, matcher query.CustomMatcher) ([]models.Track, error) {
	parser := query.NewParser()
	q := parser.Parse(queryString)
	eval := query.NewEvaluatorWithMatcher(q, matcher)

	membership := e.Library.GetMembershipMap()

	var matched []models.Track
	resources := e.Library.GetResources("track")
	for _, res := range resources {
		track := res.(models.Track)
		if eval.MatchesWithPlaylists(track, membership[track.ID]) {
			matched = append(matched, track)
		}
	}
	return matched, nil
}

// LsGroups returns all resource groups matching the given query string.
func (e *Engine) LsGroups(queryString string) ([]models.ResourceGroup, error) {
	parser := query.NewParser()
	q := parser.Parse(queryString)
	if err := q.ValidateWithFields(query.AllowedGroupFields); err != nil {
		return nil, err
	}
	eval := query.NewEvaluator(q)

	var matched []models.ResourceGroup
	resources := e.Library.GetResources("group")
	for _, res := range resources {
		group := res.(models.ResourceGroup)
		if eval.MatchesGroup(group) {
			matched = append(matched, group)
		}
	}
	return matched, nil
}
