package query

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/llttlltt/dj-library-tools/pkg/rekordbox"
)

// Evaluator checks if a track matches the given query criteria
type Evaluator struct {
	Query Query
}

func NewEvaluator(q Query) *Evaluator {
	return &Evaluator{Query: q}
}

// Matches returns true if the track meets all query criteria
func (e *Evaluator) Matches(track rekordbox.Track) bool {
	if len(e.Query.Criteria) == 0 {
		return true
	}

	allMatched := true
	for _, c := range e.Query.Criteria {
		if !e.matchCriterion(track, c) {
			allMatched = false
			break
		}
	}

	if e.Query.Negated {
		return !allMatched
	}
	return allMatched
}

func (e *Evaluator) matchCriterion(track rekordbox.Track, c Criterion) bool {
	fieldValue := e.getFieldValue(track, c.Field)

	if c.Operator == OpRange {
		return e.matchRange(fieldValue, c.Value)
	}

	switch c.Operator {
	case OpExact:
		return strings.EqualFold(fieldValue, c.Value)
	case OpSubstring:
		return strings.Contains(strings.ToLower(fieldValue), strings.ToLower(c.Value))
	case OpRegex:
		re, err := regexp.Compile(c.Value)
		if err != nil {
			return false
		}
		return re.MatchString(fieldValue)
	}

	return false
}

func (e *Evaluator) getFieldValue(track rekordbox.Track, field string) string {
	switch strings.ToLower(field) {
	case "name", "title":
		return track.Name
	case "artist":
		return track.Artist
	case "album":
		return track.Album
	case "bpm", "tempo":
		if len(track.Tempo) > 0 {
			return fmt.Sprintf("%.2f", track.Tempo[0].Bpm)
		}
		return "0.00"
	case "key":
		return track.Tonality
	case "genre":
		return track.Genre
	case "comment", "comments":
		return track.Comments
	case "year":
		return strconv.Itoa(int(track.Year))
	case "label":
		return track.Label
	case "grouping":
		return track.Grouping
	case "rating":
		return strconv.Itoa(int(track.Rating))
	case "playcount":
		return strconv.Itoa(int(track.PlayCount))
	case "added", "dateadded":
		return track.DateAdded
	case "kind":
		return track.Kind
	case "size":
		return strconv.FormatInt(track.Size, 10)
	}
	return ""
}

func (e *Evaluator) matchRange(fieldValue string, rangeValue string) bool {
	parts := strings.Split(rangeValue, "..")
	if len(parts) != 2 {
		return false
	}

	val, err := strconv.ParseFloat(fieldValue, 64)
	if err != nil {
		return false
	}

	min, errMin := strconv.ParseFloat(parts[0], 64)
	max, errMax := strconv.ParseFloat(parts[1], 64)

	if errMin != nil || errMax != nil {
		return false
	}

	return val >= min && val <= max
}
