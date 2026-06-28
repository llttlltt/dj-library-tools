package query

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

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
	"name", "parent", "folder", "items", "type",
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
	if e.Query.Root == nil {
		return true
	}
	return e.eval(e.Query.Root, track, playlists)
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

func isNumericField(field string) bool {
	switch strings.ToLower(field) {
	case "playlists", "hotcues", "memorycues", "beatgrids", "rating", "plays", "year",
		"bpm", "bitrate", "samplerate", "size", "items", "type":
		return true
	}
	return false
}

func (e *Evaluator) matchComparison(track models.Track, playlists []string, c Comparison) bool {
	field := strings.ToLower(c.Field)
	targetValue := c.Value
	if field == "added" || field == "modified" {
		targetValue = e.resolveDateShorthand(c.Value)
	}

	// Membership domain
	if field == "playlists" {
		if c.Operator == OpRange {
			return e.matchRange(strconv.Itoa(len(playlists)), targetValue)
		}
		if c.Operator != OpSubstring && c.Operator != OpExact {
			return e.matchNumericComparison(strconv.Itoa(len(playlists)), targetValue, c.Operator)
		}
		if !c.Quoted {
			if _, err := strconv.Atoi(targetValue); err == nil {
				return e.matchNumericComparison(strconv.Itoa(len(playlists)), targetValue, c.Operator)
			}
		}
		return e.matchPlaylistStrings(playlists, c)
	}

	// Delegated custom fields
	if e.Matcher != nil {
		switch field {
		case "hotcues", "memorycues", "beatgrids":
			if !c.Quoted {
				if _, err := strconv.Atoi(targetValue); err == nil || c.Operator == OpRange || c.Operator == OpGt || c.Operator == OpLt {
					return e.matchNumericCount(track, playlists, c)
				}
			}
			return e.Matcher.CustomMatch(track, c.Field, c.Operator, c.Value)
		}
	}

	fieldValue := e.getFieldValue(track, playlists, c.Field)
	if c.Operator == OpRange {
		return e.matchRange(fieldValue, targetValue)
	}

	switch c.Operator {
	case OpGt, OpGte, OpLt, OpLte:
		return e.matchNumericComparison(fieldValue, targetValue, c.Operator)
	case OpExact:
		return strings.EqualFold(fieldValue, targetValue)
	case OpSubstring:
		if isNumericField(c.Field) {
			fv, errF := strconv.ParseFloat(fieldValue, 64)
			tv, errT := strconv.ParseFloat(targetValue, 64)
			if errF == nil && errT == nil {
				return fv == tv
			}
		}
		return strings.Contains(strings.ToLower(fieldValue), strings.ToLower(targetValue))
	case OpRegex:
		re, err := regexp.Compile(targetValue)
		if err != nil {
			return false
		}
		return re.MatchString(fieldValue)
	}
	return false
}

func (e *Evaluator) matchNumericCount(track models.Track, playlists []string, c Comparison) bool {
	fieldValue := e.getFieldValue(track, playlists, c.Field)
	if c.Operator == OpRange {
		return e.matchRange(fieldValue, c.Value)
	}
	return e.matchNumericComparison(fieldValue, c.Value, c.Operator)
}

func (e *Evaluator) resolveDateShorthand(val string) string {
	val = strings.ToLower(val)
	now := time.Now()
	switch val {
	case "today":
		return now.Format("2006-01-02")
	case "yesterday":
		return now.AddDate(0, 0, -1).Format("2006-01-02")
	}
	if strings.HasPrefix(val, "-") {
		unit := val[len(val)-1:]
		amount, _ := strconv.Atoi(val[1 : len(val)-1])
		switch unit {
		case "d": return now.AddDate(0, 0, -amount).Format("2006-01-02")
		case "m": return now.AddDate(0, -amount, 0).Format("2006-01-02")
		case "y": return now.AddDate(-amount, 0, 0).Format("2006-01-02")
		}
	}
	return val
}

func (e *Evaluator) getFieldValue(track models.Track, playlists []string, field string) string {
	switch strings.ToLower(field) {
	case "id":
		return track.ID
	case "location":
		return track.Location
	case "display":
		return track.Display
	case "playlists":
		return strconv.Itoa(len(playlists))
	case "title":
		return track.Title
	case "artist":
		return track.Artist
	case "album":
		return track.Album
	case "bpm":
		return fmt.Sprintf("%.2f", track.BPM)
	case "key":
		return track.Key
	case "genre":
		return track.Genre
	case "comment":
		return track.Comment
	case "year":
		return strconv.Itoa(track.Year)
	case "label":
		return track.Label
	case "rating":
		return strconv.Itoa(track.Rating)
	case "plays":
		return strconv.Itoa(track.Plays)
	case "added":
		return track.DateAdded
	case "modified":
		return track.DateModified
	case "color":
		if e.Matcher != nil {
			if cm, ok := e.Matcher.(interface{ GetTrackColorName(string) string }); ok {
				return cm.GetTrackColorName(track.Color)
			}
		}
		return track.Color
	case "bitrate":
		return strconv.Itoa(track.Bitrate)
	case "samplerate":
		return strconv.Itoa(track.SampleRate)
	case "hotcues":
		return strconv.Itoa(track.HotCues)
	case "memorycues":
		return strconv.Itoa(track.MemoryCues)
	case "beatgrids":
		return strconv.Itoa(track.BeatgridCount)
	case "remixer":
		return track.Remixer
	case "mix":
		return track.Mix
	case "size":
		return strconv.FormatInt(track.Size, 10)
	}
	return ""
}

func (e *Evaluator) matchPlaylistStrings(playlists []string, c Comparison) bool {
	for _, p := range playlists {
		switch c.Operator {
		case OpExact:
			if strings.EqualFold(p, c.Value) {
				return true
			}
		case OpSubstring:
			if strings.Contains(strings.ToLower(p), strings.ToLower(c.Value)) {
				return true
			}
		case OpRegex:
			re, _ := regexp.Compile(c.Value)
			if re != nil && re.MatchString(p) {
				return true
			}
		}
	}
	return false
}

func (e *Evaluator) matchRange(fieldValue string, rangeValue string) bool {
	parts := strings.Split(rangeValue, "..")
	if len(parts) != 2 {
		return false
	}
	val, _ := strconv.ParseFloat(fieldValue, 64)
	min, _ := strconv.ParseFloat(parts[0], 64)
	max, _ := strconv.ParseFloat(parts[1], 64)
	return val >= min && val <= max
}

func (e *Evaluator) matchNumericComparison(fieldValue string, targetValue string, op Operator) bool {
	f, _ := strconv.ParseFloat(fieldValue, 64)
	t, _ := strconv.ParseFloat(targetValue, 64)
	switch op {
	case OpGt:
		return f > t
	case OpGte:
		return f >= t
	case OpLt:
		return f < t
	case OpLte:
		return f <= t
	case OpExact, OpSubstring:
		return f == t
	}
	return false
}

func (e *Evaluator) MatchesGroup(node models.ResourceGroup) bool {
	if e.Query.Root == nil {
		return true
	}
	return e.evalGroup(e.Query.Root, node)
}

func (e *Evaluator) evalGroup(expr Expression, node models.ResourceGroup) bool {
	switch v := expr.(type) {
	case Comparison:
		return e.matchGroupComparison(node, v)
	case Logical:
		if v.Op == "AND" {
			return e.evalGroup(v.Left, node) && e.evalGroup(v.Right, node)
		}
		return e.evalGroup(v.Left, node) || e.evalGroup(v.Right, node)
	case Not:
		return !e.evalGroup(v.Expr, node)
	}
	return false
}

func (e *Evaluator) matchGroupComparison(node models.ResourceGroup, c Comparison) bool {
	val := ""
	switch strings.ToLower(c.Field) {
	case "name":
		val = node.Name
	case "parent", "folder":
		val = node.ParentFolder
	case "items":
		val = strconv.Itoa(node.Items)
	case "type":
		val = strconv.Itoa(int(node.Type))
	}
	if c.Operator == OpRange {
		return e.matchRange(val, c.Value)
	}
	switch c.Operator {
	case OpExact:
		return strings.EqualFold(val, c.Value)
	case OpSubstring:
		if isNumericField(c.Field) {
			fv, _ := strconv.ParseFloat(val, 64)
			tv, _ := strconv.ParseFloat(c.Value, 64)
			return fv == tv
		}
		return strings.Contains(strings.ToLower(val), strings.ToLower(c.Value))
	case OpGt, OpGte, OpLt, OpLte:
		return e.matchNumericComparison(val, c.Value, c.Operator)
	}
	return false
}
