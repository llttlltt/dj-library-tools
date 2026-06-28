package query

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/llttlltt/dj-library-tools/internal/models"
)

// AllowedTrackFields is a list of valid fields for track queries.
var AllowedTrackFields = []string{
	"playlists", "title", "artist", "album", "bpm", "key", "genre", "comment",
	"year", "label", "rating", "plays", "added", "modified", "color", "bitrate",
	"samplerate", "size", "beatgrids", "hotcues", "memorycues", "id", "location",
	"remixer", "mix", "display",
}

// AllowedGroupFields is a list of valid fields for playlist and folder queries.
var AllowedGroupFields = []string{
	"name", "parent", "folder", "items", "kind",
}

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

	// Implementation-specific delegation (e.g. Cues by color)
	if e.Matcher != nil && isCalculatedField(field) && !isNumericIntent(c) {
		return e.Matcher.CustomMatch(track, c.Field, c.Operator, c.Value)
	}

	fieldValue := e.getTrackFieldValue(track, playlists, field)
	return Compare(field, fieldValue, targetValue, c.Operator)
}

func (e *Evaluator) matchGroupComparison(group models.ResourceGroup, c Comparison) bool {
	field := strings.ToLower(c.Field)
	fieldValue := e.getGroupFieldValue(group, field)
	return Compare(field, fieldValue, c.Value, c.Operator)
}

func (e *Evaluator) getTrackFieldValue(track models.Track, playlists []string, field string) string {
	switch field {
	case "id":       return track.ID
	case "location": return track.Location
	case "display":  return track.Display
	case "playlists": return strconv.Itoa(len(playlists))
	case "title":    return track.Title
	case "artist":   return track.Artist
	case "album":    return track.Album
	case "bpm":      return fmt.Sprintf("%.2f", track.BPM)
	case "key":      return track.Key
	case "genre":    return track.Genre
	case "comment":  return track.Comment
	case "year":     return strconv.Itoa(track.Year)
	case "label":    return track.Label
	case "rating":   return strconv.Itoa(track.Rating)
	case "plays":    return strconv.Itoa(track.Plays)
	case "added":    return track.DateAdded
	case "modified": return track.DateModified
	case "color":    return track.Color
	case "bitrate":  return strconv.Itoa(track.Bitrate)
	case "samplerate": return strconv.Itoa(track.SampleRate)
	case "size":     return strconv.FormatInt(track.Size, 10)
	case "remixer":  return track.Remixer
	case "mix":      return track.Mix
	case "hotcues":    return strconv.Itoa(countCues(track, models.CueTypeHot))
	case "memorycues": return strconv.Itoa(countCues(track, models.CueTypeMemory))
	case "beatgrids":  return strconv.Itoa(len(track.TempoMarkers))
	}
	return ""
}

func (e *Evaluator) getGroupFieldValue(group models.ResourceGroup, field string) string {
	switch field {
	case "name":   return group.Name
	case "parent", "folder": return group.ParentFolder
	case "items":  return strconv.Itoa(group.Items)
	case "kind":   return string(group.Kind)
	}
	return ""
}

// Helpers

func isCalculatedField(field string) bool {
	return field == "hotcues" || field == "memorycues" || field == "beatgrids"
}

func isNumericIntent(c Comparison) bool {
	if c.Quoted { return false }
	if c.Operator == OpRange || c.Operator == OpGt || c.Operator == OpLt { return true }
	_, err := strconv.Atoi(c.Value)
	return err == nil
}

func countCues(track models.Track, cueType models.CueType) int {
	count := 0
	for _, cp := range track.CuePoints {
		if cp.Type == cueType { count++ }
	}
	return count
}
