package query

import (
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/llttlltt/dj-library-tools/pkg/rekordbox"
)

type Evaluator struct {
	Query Query
}

func NewEvaluator(q Query) *Evaluator {
	return &Evaluator{Query: q}
}

func (e *Evaluator) Matches(track rekordbox.Track) bool {
	return e.MatchesWithPlaylists(track, nil)
}

func (e *Evaluator) MatchesWithPlaylists(track rekordbox.Track, playlists []string) bool {
	if e.Query.Root == nil {
		return true
	}
	return e.eval(e.Query.Root, track, playlists)
}

func (e *Evaluator) eval(expr Expression, track rekordbox.Track, playlists []string) bool {
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

// isNumericField returns true for fields whose values are always numeric.
// For these fields, OpSubstring (`:`) performs exact numeric equality rather than
// a string-contains check, avoiding false positives such as playlistcount:0
// matching a track in 10 playlists because "10" contains "0".
func isNumericField(field string) bool {
	switch strings.ToLower(field) {
	case "playlistcount", "playlists",
		"hotcuecount", "hotcues",
		"memorycuecount", "memorycues",
		"beatgrids", "beatgridcount", "tempos", "tempocount",
		"rating", "playcount", "year",
		"bpm", "tempo",
		"bitrate", "kbps",
		"samplerate", "khz",
		"time", "length", "duration",
		"size", "id", "trackid",
		"disc", "discnumber",
		"track", "tracknumber",
		"entries",
		"type":
		return true
	}
	return false
}

func (e *Evaluator) matchComparison(track rekordbox.Track, playlists []string, c Comparison) bool {
	field := strings.ToLower(c.Field)
	switch field {
	case "playlist", "playlists":
		// If it's a numeric comparison or range, treat it as a count check
		if _, err := strconv.Atoi(c.Value); err == nil || strings.Contains(c.Value, "..") {
			break
		}
		return e.matchPlaylist(playlists, c)
	case "playlistcount":
		break
	case "hotcue":
		return e.matchSpecificCue(track, c, true)
	case "memorycue":
		return e.matchSpecificCue(track, c, false)
	case "hotcues":
		if _, err := strconv.Atoi(c.Value); err == nil || strings.Contains(c.Value, "..") {
			break
		}
		return e.matchAnyCue(track, true, c.Value)
	case "memorycues":
		if _, err := strconv.Atoi(c.Value); err == nil || strings.Contains(c.Value, "..") {
			break
		}
		return e.matchAnyCue(track, false, c.Value)
	}

	fieldValue := e.getFieldValue(track, playlists, c.Field)
	if c.Operator == OpRange {
		return e.matchRange(fieldValue, c.Value)
	}

	switch c.Operator {
	case OpGt, OpGte, OpLt, OpLte:
		return e.matchNumericComparison(fieldValue, c.Value, c.Operator)
	case OpExact:
		return strings.EqualFold(fieldValue, c.Value)
	case OpSubstring:
		// For numeric fields, use exact float equality to prevent substring false-positives
		// (e.g. playlistcount:0 must not match a track in 10 playlists).
		if isNumericField(c.Field) {
			fv, errF := strconv.ParseFloat(fieldValue, 64)
			tv, errT := strconv.ParseFloat(c.Value, 64)
			if errF == nil && errT == nil {
				return fv == tv
			}
		}
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

func (e *Evaluator) getFieldValue(track rekordbox.Track, playlists []string, field string) string {
	switch strings.ToLower(field) {
	case "playlistcount", "playlists":
		return strconv.Itoa(len(playlists))
	case "name", "title":
		return track.Name
	case "artist":
		return track.Artist
	case "album":
		return track.Album
	case "bpm", "tempo":
		return track.AverageBpm
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
		// Normalize 0-255 scale to 0-5 stars
		return strconv.Itoa(int(track.Rating / 51))
	case "playcount":
		return strconv.Itoa(int(track.PlayCount))
	case "added", "dateadded":
		return track.DateAdded
	case "kind":
		return track.Kind
	case "size":
		return strconv.FormatInt(track.Size, 10)
	case "beatgrids", "beatgridcount", "tempos", "tempocount":
		return strconv.Itoa(len(track.Tempo))
	case "hotcuecount", "hotcues":
		count := 0
		for _, pm := range track.PositionMark {
			if pm.Num != -1 {
				count++
			}
		}
		return strconv.Itoa(count)
	case "memorycuecount", "memorycues":
		count := 0
		for _, pm := range track.PositionMark {
			if pm.Num == -1 {
				count++
			}
		}
		return strconv.Itoa(count)
	case "id", "trackid":
		return strconv.Itoa(track.TrackID)
	case "composer":
		return track.Composer
	case "time", "length", "duration":
		return strconv.Itoa(int(track.TotalTime))
	case "disc", "discnumber":
		return strconv.Itoa(int(track.DiscNumber))
	case "track", "tracknumber":
		return strconv.Itoa(int(track.TrackNumber))
	case "bitrate", "kbps":
		return strconv.Itoa(int(track.BitRate))
	case "samplerate", "khz":
		return strconv.Itoa(int(track.SampleRate))
	case "path", "file", "location":
		return track.Location
	case "remixer":
		return track.Remixer
	case "mix", "version":
		return track.Mix
	}
	return ""
}

func (e *Evaluator) matchPlaylist(playlists []string, c Comparison) bool {
	for _, p := range playlists {
		matched := false
		switch c.Operator {
		case OpExact:
			matched = strings.EqualFold(p, c.Value)
		case OpSubstring, OpRange: // query parser might use Range operator for generic colon syntax
			matched = strings.Contains(strings.ToLower(p), strings.ToLower(c.Value))
		case OpRegex:
			re, err := regexp.Compile(c.Value)
			if err != nil {
				return false
			}
			matched = re.MatchString(p)
		}
		if matched {
			return true
		}
	}
	return false
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

func (e *Evaluator) matchNumericComparison(fieldValue string, targetValue string, op Operator) bool {
	fieldNum, errF := strconv.ParseFloat(fieldValue, 64)
	targetNum, errT := strconv.ParseFloat(targetValue, 64)
	if errF != nil || errT != nil {
		return false
	}

	switch op {
	case OpGt:
		return fieldNum > targetNum
	case OpGte:
		return fieldNum >= targetNum
	case OpLt:
		return fieldNum < targetNum
	case OpLte:
		return fieldNum <= targetNum
	}
	return false
}

func (e *Evaluator) matchSpecificCue(track rekordbox.Track, c Comparison, isHotCue bool) bool {
	// Trim any spaces that might have been introduced during parsing for property segments
	// but KEEP colons as they are our separators.
	val := strings.ReplaceAll(c.Value, " ", "")
	parts := strings.Split(val, ":")
	if len(parts) == 0 {
		return false
	}

	var targetPMs []rekordbox.PositionMark
	if isHotCue {
		slotStr := parts[0]
		if len(slotStr) != 1 {
			return false
		}
		targetSlot := int(strings.ToLower(slotStr)[0] - 'a')
		for _, pm := range track.PositionMark {
			if int(pm.Num) == targetSlot {
				targetPMs = append(targetPMs, pm)
			}
		}
	} else {
		idx, err := strconv.Atoi(parts[0])
		if err != nil || idx < 1 {
			return false
		}
		cues := e.getSortedMemoryCues(track)
		if idx <= len(cues) {
			targetPMs = append(targetPMs, cues[idx-1])
		}
	}

	if len(targetPMs) == 0 {
		return false
	}

	if len(parts) == 1 {
		return true
	}

	for _, pm := range targetPMs {
		if e.matchCueProperties(pm, parts[1:]) {
			return true
		}
	}
	return false
}

func (e *Evaluator) matchCueProperties(pm rekordbox.PositionMark, props []string) bool {
	if len(props) == 0 {
		return true
	}

	i := 0
	for i < len(props) {
		prop := strings.ToLower(props[i])
		matched := false

		if prop == "label" {
			targetLabel := ""
			if i+1 < len(props) {
				targetLabel = props[i+1]
				i++
			}
			if targetLabel == `""` || targetLabel == "" || targetLabel == "none" || targetLabel == "empty" {
				matched = pm.Name == ""
			} else {
				matched = strings.Contains(strings.ToLower(pm.Name), strings.ToLower(targetLabel))
			}
		} else if prop == "loop" || prop == "active loop" {
			matched = pm.Type == 4
		} else if prop == `""` { // Added explicit check for quoted empty string
			matched = pm.Name == ""
		} else {
			matched = e.matchColor(pm, prop)
		}

		if !matched {
			return false
		}
		i++
	}

	return true
}

func (e *Evaluator) matchAnyCue(track rekordbox.Track, isHotCue bool, value string) bool {
	val := strings.ReplaceAll(value, " ", "")
	parts := strings.Split(val, ":")
	for _, pm := range track.PositionMark {
		if isHotCue && pm.Num == -1 {
			continue
		}
		if !isHotCue && pm.Num != -1 {
			continue
		}
		if e.matchCueProperties(pm, parts) {
			return true
		}
	}
	return false
}

func (e *Evaluator) getSortedMemoryCues(track rekordbox.Track) []rekordbox.PositionMark {
	var cues []rekordbox.PositionMark
	for _, pm := range track.PositionMark {
		if pm.Num == -1 {
			cues = append(cues, pm)
		}
	}
	sort.Slice(cues, func(i, j int) bool {
		return cues[i].Start > cues[j].Start
	})
	return cues
}

func (e *Evaluator) matchColor(pm rekordbox.PositionMark, color string) bool {
	c := strings.ToLower(color)
	switch c {
	case "green":
		return pm.Red == 40 && pm.Green == 226 && pm.Blue == 20
	case "aqua":
		return pm.Red == 0 && pm.Green == 224 && pm.Blue == 255
	case "orange":
		return pm.Red == 224 && pm.Green == 100 && pm.Blue == 27
	case "red":
		return pm.Red == 230 && pm.Green == 40 && pm.Blue == 40
	case "blue":
		return pm.Red == 48 && pm.Green == 90 && pm.Blue == 255
	case "purple":
		return pm.Red == 180 && pm.Green == 50 && pm.Blue == 255
	case "yellow":
		return pm.Red == 255 && pm.Green == 255 && pm.Blue == 0
	case "pink":
		return pm.Red == 255 && pm.Green == 50 && pm.Blue == 180
	case "no color", "none":
		return (pm.Red == 0 && pm.Green == 0 && pm.Blue == 0) || (pm.Red == 40 && pm.Green == 226 && pm.Blue == 20 && pm.Num == -1)
	}
	return false
}

// MatchesNode evaluates whether a rekordbox Node matches the query.
// parentFolder is the name of the node's direct parent ("" if at root level).
func (e *Evaluator) MatchesNode(node rekordbox.Node, parentFolder string) bool {
	if e.Query.Root == nil {
		return true
	}
	return e.evalNode(e.Query.Root, node, parentFolder)
}

func (e *Evaluator) evalNode(expr Expression, node rekordbox.Node, parentFolder string) bool {
	switch v := expr.(type) {
	case Comparison:
		return e.matchNodeComparison(node, parentFolder, v)
	case Logical:
		if v.Op == "AND" {
			return e.evalNode(v.Left, node, parentFolder) && e.evalNode(v.Right, node, parentFolder)
		}
		return e.evalNode(v.Left, node, parentFolder) || e.evalNode(v.Right, node, parentFolder)
	case Not:
		return !e.evalNode(v.Expr, node, parentFolder)
	}
	return false
}

func (e *Evaluator) matchNodeComparison(node rekordbox.Node, parentFolder string, c Comparison) bool {
	fieldValue := e.getNodeFieldValue(node, parentFolder, c.Field)
	if c.Operator == OpRange {
		return e.matchRange(fieldValue, c.Value)
	}
	switch c.Operator {
	case OpGt, OpGte, OpLt, OpLte:
		return e.matchNumericComparison(fieldValue, c.Value, c.Operator)
	case OpExact:
		return strings.EqualFold(fieldValue, c.Value)
	case OpSubstring:
		if isNumericField(c.Field) {
			fv, errF := strconv.ParseFloat(fieldValue, 64)
			tv, errT := strconv.ParseFloat(c.Value, 64)
			if errF == nil && errT == nil {
				return fv == tv
			}
		}
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

func (e *Evaluator) getNodeFieldValue(node rekordbox.Node, parentFolder string, field string) string {
	switch strings.ToLower(field) {
	case "name":
		return node.Name
	case "folder", "parent":
		return parentFolder
	case "entries":
		return strconv.Itoa(int(node.Entries))
	case "type":
		return strconv.Itoa(int(node.Type))
	}
	return ""
}
