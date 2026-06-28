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

// LsPlaylists returns all playlist nodes matching the given query string.
func (e *Engine) LsPlaylists(queryString string) ([]models.ResourceGroup, error) {
	return e.lsNodes(queryString, models.GroupTypePlaylist)
}

// LsFolders returns all folder nodes matching the given query string.
func (e *Engine) LsFolders(queryString string) ([]models.ResourceGroup, error) {
	return e.lsNodes(queryString, models.GroupTypeFolder)
}

func (e *Engine) lsNodes(queryString string, nodeType models.GroupType) ([]models.ResourceGroup, error) {
	parser := query.NewParser()
	q := parser.Parse(queryString)
	if err := q.ValidateWithFields(query.AllowedNodeFields); err != nil {
		return nil, err
	}
	eval := query.NewEvaluator(q)

	var matched []models.ResourceGroup
	for _, node := range e.Library.GetPlaylists() {
		if node.Type == nodeType {
			if eval.MatchesGroup(node) {
				matched = append(matched, node)
			}
		}
	}
	return matched, nil
}
