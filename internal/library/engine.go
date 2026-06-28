package library

import (
	"github.com/llttlltt/dj-library-tools/internal/models"
	"github.com/llttlltt/dj-library-tools/internal/query"
)

// Engine performs operations on a library using queries
type Engine struct {
	Library Library
}

func NewEngine(lib Library) *Engine {
	return &Engine{
		Library: lib,
	}
}

// Ls returns all tracks that match the given query string
func (e *Engine) Ls(queryString string, matcher query.CustomMatcher) ([]models.Track, error) {
	parser := query.NewParser()
	q := parser.Parse(queryString)
	if err := q.ValidateWithFields(query.AllowedTrackFields); err != nil {
		return nil, err
	}
	eval := query.NewEvaluatorWithMatcher(q, matcher)

	membership := e.Library.GetMembershipMap()

	var matched []models.Track
	for _, track := range e.Library.GetTracks() {
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
	// We allow AllowedNodeFields here. The query package will handle the 'type' field.
	if err := q.ValidateWithFields(query.AllowedNodeFields); err != nil {
		return nil, err
	}
	eval := query.NewEvaluator(q)

	var matched []models.ResourceGroup
	for _, group := range e.Library.GetPlaylists() {
		if eval.MatchesGroup(group) {
			matched = append(matched, group)
		}
	}
	return matched, nil
}
