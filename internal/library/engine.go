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

// Modify applies changes to matched tracks
func (e *Engine) Modify(queryString string, matcher query.CustomMatcher, action func(track models.Track, changes map[string]string) error, changes map[string]string) (int, error) {
	tracks, err := e.Ls(queryString, matcher)
	if err != nil {
		return 0, err
	}

	for _, t := range tracks {
		if err := action(t, changes); err != nil {
			return 0, err
		}
	}
	return len(tracks), nil
}
