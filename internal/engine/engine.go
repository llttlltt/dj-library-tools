package engine

import (
	"github.com/llttlltt/dj-library-tools/internal/query"
	"github.com/llttlltt/dj-library-tools/pkg/rekordbox"
)

// Engine performs operations on a Rekordbox library using queries
type Engine struct {
	Library *rekordbox.RekordboxLibraryXML
}

func NewEngine(lib *rekordbox.RekordboxLibraryXML) *Engine {
	return &Engine{Library: lib}
}

// Ls returns all tracks that match the given query string
func (e *Engine) Ls(queryString string) ([]rekordbox.Track, error) {
	parser := query.NewParser()
	q := parser.Parse(queryString)
	eval := query.NewEvaluator(q)

	var matched []rekordbox.Track
	for _, track := range e.Library.Collection.TRACK {
		if eval.Matches(track) {
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
		if t.Tonality != "" {
			res.Keys[t.Tonality]++
		}
		if t.Artist != "" {
			res.Artists[t.Artist]++
		}
		if len(t.Tempo) > 0 {
			res.TotalTempo += t.Tempo[0].Bpm
		}
	}
	res.AvgBPM = res.TotalTempo / float64(len(tracks))

	return res, nil
}

// Modify applies changes to matched tracks
// Example: modifyQuery: "artist:Four", changes: map[string]string{"comment": "Verified"}
func (e *Engine) Modify(queryString string, changes map[string]string) (int, error) {
	parser := query.NewParser()
	q := parser.Parse(queryString)
	eval := query.NewEvaluator(q)

	modifyCount := 0
	for i := range e.Library.Collection.TRACK {
		if eval.Matches(e.Library.Collection.TRACK[i]) {
			e.applyChanges(&e.Library.Collection.TRACK[i], changes)
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
