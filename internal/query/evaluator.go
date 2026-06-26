package query

import (
	"fmt"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/llttlltt/dj-library-tools/internal/models"
)

// AllowedTrackFields is a list of valid fields for track queries.
var AllowedTrackFields = []string{
	"playlistcount", "title", "artist", "album", "bpm", "key", "genre", "comment",
	"year", "label", "rating", "playcount", "added", "modified", "color", "bitrate",
	"samplerate", "size", "beatgrids", "hotcues", "memorycues", "id", "location",
	"remixer", "mix", "playlist",
}

// AllowedNodeFields is a list of valid fields for playlist and folder queries.
var AllowedNodeFields = []string{
	"name", "parent", "folder", "entries", "count", "type",
}

type Evaluator struct {
	Query Query
}

func NewEvaluator(q Query) *Evaluator {
	return &Evaluator{Query: q}
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
	case "playlistcount", "hotcues", "memorycues", "beatgrids", "rating", "playcount", "year",
		"bpm", "bitrate", "samplerate", "size", "entries", "count", "type":
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

	if field == "playlist" {
		return e.matchPlaylist(playlists, c)
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
	case "playlistcount":
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
	case "playcount":
		return strconv.Itoa(track.PlayCount)
	case "added":
		return track.DateAdded
	case "modified":
		return track.DateModified
	case "color":
		if name := e.getTrackColorName(track.Color); name != track.Color {
			return name
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
	case "id":
		return track.ID
	case "location":
		return track.Location
	case "remixer":
		return track.Remixer
	case "mix":
		return track.Mix
	case "size":
		return strconv.FormatInt(track.Size, 10)
	}
	return ""
}

func (e *Evaluator) getTrackColorName(hex string) string {
	switch strings.ToUpper(hex) {
	case "0XFF007F":
		return "pink"
	case "0XFF0000":
		return "red"
	case "0XFFA500":
		return "orange"
	case "0XFFFF00":
		return "yellow"
	case "0X00FF00":
		return "green"
	case "0X25FDE9":
		return "aqua"
	case "0X0000FF":
		return "blue"
	case "0X660099":
		return "purple"
	}
	return hex
}

func (e *Evaluator) matchPlaylist(playlists []string, c Comparison) bool {
	for _, p := range playlists {
		switch c.Operator {
		case OpExact: if strings.EqualFold(p, c.Value) { return true }
		case OpSubstring: if strings.Contains(strings.ToLower(p), strings.ToLower(c.Value)) { return true }
		case OpRegex:
			re, _ := regexp.Compile(c.Value)
			if re != nil && re.MatchString(p) { return true }
		}
	}
	return false
}

func (e *Evaluator) matchRange(fieldValue string, rangeValue string) bool {
	parts := strings.Split(rangeValue, "..")
	if len(parts) != 2 { return false }
	val, _ := strconv.ParseFloat(fieldValue, 64)
	min, _ := strconv.ParseFloat(parts[0], 64)
	max, _ := strconv.ParseFloat(parts[1], 64)
	return val >= min && val <= max
}

func (e *Evaluator) matchNumericComparison(fieldValue string, targetValue string, op Operator) bool {
	f, _ := strconv.ParseFloat(fieldValue, 64)
	t, _ := strconv.ParseFloat(targetValue, 64)
	switch op {
	case OpGt: return f > t
	case OpGte: return f >= t
	case OpLt: return f < t
	case OpLte: return f <= t
	}
	return false
}

func (e *Evaluator) MatchesNode(node models.Node) bool {
	if e.Query.Root == nil { return true }
	return e.evalNode(e.Query.Root, node)
}

func (e *Evaluator) evalNode(expr Expression, node models.Node) bool {
	switch v := expr.(type) {
	case Comparison: return e.matchNodeComparison(node, v)
	case Logical:
		if v.Op == "AND" { return e.evalNode(v.Left, node) && e.evalNode(v.Right, node) }
		return e.evalNode(v.Left, node) || e.evalNode(v.Right, node)
	case Not: return !e.evalNode(v.Expr, node)
	}
	return false
}

func (e *Evaluator) matchNodeComparison(node models.Node, c Comparison) bool {
	val := ""
	switch strings.ToLower(c.Field) {
	case "name": val = node.Name
	case "parent", "folder": val = node.ParentFolder
	case "entries": val = strconv.Itoa(node.Entries)
	case "type": val = strconv.Itoa(node.Type)
	}
	if c.Operator == OpRange { return e.matchRange(val, c.Value) }
	switch c.Operator {
	case OpExact: return strings.EqualFold(val, c.Value)
	case OpSubstring:
		if isNumericField(c.Field) {
			fv, _ := strconv.ParseFloat(val, 64)
			tv, _ := strconv.ParseFloat(c.Value, 64)
			return fv == tv
		}
		return strings.Contains(strings.ToLower(val), strings.ToLower(c.Value))
	}
	return false
}

func printTop(m map[string]int, title string, limit int) {
	if len(m) == 0 {
		return
	}

	type kv struct {
		Key   string
		Value int
	}

	var ss []kv
	for k, v := range m {
		ss = append(ss, kv{k, v})
	}

	sort.Slice(ss, func(i, j int) bool {
		return ss[i].Value > ss[j].Value
	})

	headerFmt := color.New(color.FgCyan, color.Bold, color.Underline).SprintFunc()
	labelFmt := color.New(color.FgHiWhite).SprintFunc()
	valFmt := color.New(color.FgHiMagenta).SprintFunc()

	fmt.Printf("\n%s\n", headerFmt(title))
	for i, kv := range ss {
		if i >= limit {
			break
		}
		fmt.Printf("%-20s %s\n", labelFmt(kv.Key), valFmt(fmt.Sprintf("%d", kv.Value)))
	}
}

type StatResult struct {
	Count      int
	AvgBPM     float64
	Genres     map[string]int
	Labels     map[string]int
	Keys       map[string]int
	Artists    map[string]int
	TotalTempo float64
}
