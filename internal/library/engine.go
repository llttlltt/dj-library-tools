package library

import (
	"fmt"
	"strconv"

	"github.com/llttlltt/dj-library-tools/internal/models"
	"github.com/llttlltt/dj-library-tools/internal/query"
	"github.com/llttlltt/dj-library-tools/internal/rekordbox"
)

// Engine performs operations on a library using queries
type Engine struct {
	Library  Library
	trackMap map[string][]string // Map Track ID string to list of playlist names
}

func NewEngine(lib Library) *Engine {
	e := &Engine{
		Library:  lib,
		trackMap: make(map[string][]string),
	}
	e.indexPlaylists()
	return e
}

func (e *Engine) indexPlaylists() {
	// For Rekordbox specifically, we need to index membership.
	if r, ok := e.Library.(*RekordboxLibrary); ok {
		e.walkRekordboxPlaylists(r.XML.Playlists.Node.Node)
	}
}

func (e *Engine) walkRekordboxPlaylists(nodes []rekordbox.Node) {
	for _, node := range nodes {
		if node.Type == 1 {
			for _, t := range node.TRACK {
				e.trackMap[t.Key] = append(e.trackMap[t.Key], node.Name)
			}
		}
		if len(node.Node) > 0 {
			e.walkRekordboxPlaylists(node.Node)
		}
	}
}

// Ls returns all tracks that match the given query string
func (e *Engine) Ls(queryString string) ([]models.Track, error) {
	parser := query.NewParser()
	q := parser.Parse(queryString)
	if err := q.ValidateWithFields(query.AllowedTrackFields); err != nil {
		return nil, err
	}
	eval := query.NewEvaluator(q)

	var matched []models.Track
	for _, track := range e.Library.GetTracks() {
		if eval.MatchesWithPlaylists(track, e.trackMap[track.ID]) {
			matched = append(matched, track)
		}
	}
	return matched, nil
}

// StatResult holds statistical analysis of a selection
type StatResult struct {
	Count      int
	AvgBPM     float64
	Genres     map[string]int
	Labels     map[string]int
	Keys       map[string]int
	Artists    map[string]int
	TotalTempo float64
}

// Stat performs statistical analysis on matched tracks
func (e *Engine) Stat(queryString string) (*StatResult, error) {
	tracks, err := e.Ls(queryString)
	if err != nil {
		return nil, err
	}

	res := &StatResult{
		Count:   len(tracks),
		Genres:  make(map[string]int),
		Labels:  make(map[string]int),
		Keys:    make(map[string]int),
		Artists: make(map[string]int),
	}

	if len(tracks) == 0 {
		return res, nil
	}

	for _, t := range tracks {
		if t.Genre != "" {
			res.Genres[t.Genre]++
		}
		if t.Label != "" {
			res.Labels[t.Label]++
		}
		if t.Key != "" {
			res.Keys[t.Key]++
		}
		if t.Artist != "" {
			res.Artists[t.Artist]++
		}
		res.TotalTempo += t.BPM
	}
	res.AvgBPM = res.TotalTempo / float64(len(tracks))

	return res, nil
}

// Modify applies changes to matched tracks
func (e *Engine) Modify(queryString string, changes map[string]string) (int, error) {
	// This still requires the underlying writable library to save back.
	// For now we'll only support this on RekordboxLibrary.
	lib, ok := e.Library.(*RekordboxLibrary)
	if !ok {
		return 0, fmt.Errorf("modify only supported on rekordbox-backed libraries")
	}

	parser := query.NewParser()
	q := parser.Parse(queryString)
	eval := query.NewEvaluator(q)

	modifyCount := 0
	tracks := lib.XML.Collection.TRACK
	for i := range tracks {
		rt := &tracks[i]
		if eval.MatchesWithPlaylists(rt.ToNeutral(), e.trackMap[strconv.Itoa(rt.TrackID)]) {
			e.applyChanges(rt, changes)
			modifyCount++
		}
	}
	return modifyCount, nil
}

func (e *Engine) applyChanges(track *rekordbox.Track, changes map[string]string) {
	for field, value := range changes {
		switch field {
		case "comment", "comments":
			track.Comments = value
		case "genre":
			track.Genre = value
		case "label":
			track.Label = value
		case "artist":
			track.Artist = value
		case "album":
			track.Album = value
		}
	}
}
