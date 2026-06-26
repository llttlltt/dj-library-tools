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
	"github.com/llttlltt/dj-library-tools/pkg/rekordbox"
)

// AllowedTrackFields is a list of valid fields for track queries.
var AllowedTrackFields = []string{
	"playlists", "title", "artist", "album", "bpm", "key", "genre", "comment",
	"year", "label", "rating", "plays", "added", "modified", "color", "bitrate",
	"samplerate", "size", "beatgrids", "hotcues", "memorycues", "id", "location",
	"remixer", "mix",
}

// AllowedNodeFields is a list of valid fields for playlist and folder queries.
var AllowedNodeFields = []string{
	"name", "parent", "folder", "items", "type",
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

	// Membership domain: handles both "how many?" and "is in?"
	if field == "playlists" {
		// If range or comparison operator, force numeric evaluation
		if c.Operator == OpRange {
			return e.matchRange(strconv.Itoa(len(playlists)), targetValue)
		}
		if c.Operator != OpSubstring && c.Operator != OpExact {
			return e.matchNumericComparison(strconv.Itoa(len(playlists)), targetValue, c.Operator)
		}
		// If substring/exact match, check if it's a number (unless quoted)
		if !c.Quoted {
			if _, err := strconv.Atoi(targetValue); err == nil {
				// This matches our "Type-Inference" precedence: raw numbers are counts
				return e.matchNumericComparison(strconv.Itoa(len(playlists)), targetValue, c.Operator)
			}
		}
		// Otherwise, search playlist names
		return e.matchPlaylistStrings(playlists, c)
	}

	// Cue domains: handles both "how many?" and properties
	if field == "hotcues" || field == "memorycues" {
		// If comparison or range, or numeric exact match, it's a count
		if c.Operator == OpRange {
			return e.matchRange(e.getFieldValue(track, playlists, field), targetValue)
		}
		if c.Operator != OpSubstring && c.Operator != OpExact {
			return e.matchNumericComparison(e.getFieldValue(track, playlists, field), targetValue, c.Operator)
		}
		if !c.Quoted {
			if _, err := strconv.Atoi(targetValue); err == nil {
				return e.matchNumericComparison(e.getFieldValue(track, playlists, field), targetValue, c.Operator)
			}
		}
		// Otherwise, searching cue properties (regex/exact/substring on color/comment)
		return e.matchCueProperties(track, field, c)
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
	case "id":
		return track.ID
	case "location":
		return track.Location
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

func (e *Evaluator) matchCueProperties(track models.Track, field string, c Comparison) bool {
	// Reconstruct the full value including ID targeting if it was a double-colon style query
	// (though the parser currently flattens it)
	target := strings.ToLower(c.Value)

	// Rekordbox specific: extracting cue data from Raw if present
	if rt, ok := track.Raw.(rekordbox.Track); ok {
		if field == "hotcues" {
			for _, pm := range rt.PositionMark {
				if pm.Num == -1 {
					continue
				} // skip memory cues
				if e.matchCueMetadata(pm, target, c.Operator) {
					return true
				}
			}
		} else {
			for _, pm := range rt.PositionMark {
				if pm.Num != -1 {
					continue
				} // skip hot cues
				if e.matchCueMetadata(pm, target, c.Operator) {
					return true
				}
			}
		}
	}
	return false
}

func (e *Evaluator) matchCueMetadata(pm rekordbox.PositionMark, target string, op Operator) bool {
	// Match by name/comment
	if op == OpExact {
		if strings.EqualFold(pm.Name, target) {
			return true
		}
	} else if strings.Contains(strings.ToLower(pm.Name), target) {
		return true
	}

	// Match by color name (Rekordbox hex map)
	colorName := strings.ToLower(e.getHotCueColorName(pm))
	if op == OpExact {
		if colorName == target {
			return true
		}
	} else if strings.Contains(colorName, target) {
		return true
	}

	return false
}

func (e *Evaluator) getHotCueColorName(pm rekordbox.PositionMark) string {
	// Simple map for pad colors
	rgb := fmt.Sprintf("%02X%02X%02X", pm.Red, pm.Green, pm.Blue)
	switch rgb {
	case "E62828":
		return "red"
	case "DE44CF":
		return "hotpink"
	case "FFFF00", "B4BE04", "C3AF04":
		return "yellow"
	case "28E214", "10B176":
		return "green"
	case "00E0FF", "50B4FF":
		return "aqua"
	case "305AFF", "6473FF":
		return "blue"
	case "B432FF", "AA72FF":
		return "purple"
	case "E0641B", "FFA500":
		return "orange"
	}
	return ""
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

func (e *Evaluator) MatchesNode(node models.Node) bool {
	if e.Query.Root == nil {
		return true
	}
	return e.evalNode(e.Query.Root, node)
}

func (e *Evaluator) evalNode(expr Expression, node models.Node) bool {
	switch v := expr.(type) {
	case Comparison:
		return e.matchNodeComparison(node, v)
	case Logical:
		if v.Op == "AND" {
			return e.evalNode(v.Left, node) && e.evalNode(v.Right, node)
		}
		return e.evalNode(v.Left, node) || e.evalNode(v.Right, node)
	case Not:
		return !e.evalNode(v.Expr, node)
	}
	return false
}

func (e *Evaluator) matchNodeComparison(node models.Node, c Comparison) bool {
	val := ""
	switch strings.ToLower(c.Field) {
	case "name":
		val = node.Name
	case "parent", "folder":
		val = node.ParentFolder
	case "items":
		val = strconv.Itoa(node.Items)
	case "type":
		val = strconv.Itoa(node.Type)
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
