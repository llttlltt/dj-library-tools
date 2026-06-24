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
	case "playlistcount",
		"hotcues",
		"memorycues",
		"beatgrids",
		"rating", "playcount", "year",
		"bpm",
		"bitrate",
		"samplerate",
		"time",
		"size", "id",
		"disc",
		"track",
		"entries", "count",
		"type", "inizio", "battito", "start", "end",
		"num":
		return true
	}
	return false
}

func (e *Evaluator) matchComparison(track rekordbox.Track, playlists []string, c Comparison) bool {
	field := strings.ToLower(c.Field)
	switch field {
	case "playlist":
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
	case "tempo":
		return e.matchSpecificTempo(track, c)
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
	case "playlistcount":
		return strconv.Itoa(len(playlists))
	case "title":
		return track.Name
	case "artist":
		return track.Artist
	case "album":
		return track.Album
	case "bpm":
		return track.AverageBpm
	case "key":
		return track.Tonality
	case "genre":
		return track.Genre
	case "comment":
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
	case "added":
		return track.DateAdded
	case "modified":
		return track.DateModified
	case "played":
		return track.LastPlayed
	case "color":
		// Check track-level colour attribute (Standard Spec)
		if name := e.getTrackColorName(track.Colour); name != track.Colour {
			return name
		}
		return track.Colour
	case "kind":
		return track.Kind
	case "size":
		return strconv.FormatInt(track.Size, 10)
	case "beatgrids":
		return strconv.Itoa(len(track.Tempo))
	case "hotcues":
		count := 0
		for _, pm := range track.PositionMark {
			if pm.Num != -1 {
				count++
			}
		}
		return strconv.Itoa(count)
	case "memorycues":
		count := 0
		for _, pm := range track.PositionMark {
			if pm.Num == -1 {
				count++
			}
		}
		return strconv.Itoa(count)
	case "id":
		return strconv.Itoa(track.TrackID)
	case "composer":
		return track.Composer
	case "time":
		return strconv.Itoa(int(track.TotalTime))
	case "disc":
		return strconv.Itoa(int(track.DiscNumber))
	case "track":
		return strconv.Itoa(int(track.TrackNumber))
	case "bitrate":
		return strconv.Itoa(int(track.BitRate))
	case "samplerate":
		return strconv.Itoa(int(track.SampleRate))
	case "location":
		return track.Location
	case "remixer":
		return track.Remixer
	case "mix":
		return track.Mix
	}
	return ""
}

const (
	// NullValue is a global constant representing no value or empty state.
	NullValue = "none"
)

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
			// Check if it's a numeric slot (0-7 for A-H)
			if slot, err := strconv.Atoi(slotStr); err == nil && slot >= 0 {
				for _, pm := range track.PositionMark {
					if int(pm.Num) == slot {
						targetPMs = append(targetPMs, pm)
					}
				}
			} else {
				return false
			}
		} else {
			targetSlot := int(strings.ToLower(slotStr)[0] - 'a')
			for _, pm := range track.PositionMark {
				if int(pm.Num) == targetSlot {
					targetPMs = append(targetPMs, pm)
				}
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

		switch prop {
		case "comment":
			targetLabel := ""
			if i+1 < len(props) {
				targetLabel = props[i+1]
				i++
			}
			if targetLabel == `""` || targetLabel == "" || targetLabel == NullValue || targetLabel == "empty" {
				matched = pm.Name == ""
			} else {
				matched = strings.Contains(strings.ToLower(pm.Name), strings.ToLower(targetLabel))
			}
		case "loop", "active-loop", "activeloop":
			matched = pm.Type == 4
		case "type":
			if i+1 < len(props) {
				targetType := strings.ToLower(props[i+1])
				i++
				switch targetType {
				case "cue":
					matched = pm.Type == 0
				case "fade-in", "fadein":
					matched = pm.Type == 1
				case "fade-out", "fadeout":
					matched = pm.Type == 2
				case "load":
					matched = pm.Type == 3
				case "loop":
					matched = pm.Type == 4
				default:
					iv, _ := strconv.Atoi(targetType)
					matched = pm.Type == int32(iv)
				}
			}
		case "start", "time", "pos":
			if i+1 < len(props) {
				val := props[i+1]
				i++
				matched = e.matchNumericProperty(pm.Start, val)
			}
		case "end":
			if i+1 < len(props) {
				val := props[i+1]
				i++
				matched = e.matchNumericProperty(pm.End, val)
			}
		case "num":
			if i+1 < len(props) {
				val := props[i+1]
				i++
				matched = e.matchNumericProperty(strconv.Itoa(int(pm.Num)), val)
			}
		case `""`:
			matched = pm.Name == ""
		default:
			matched = e.matchColor(pm, prop)
		}

		if !matched {
			return false
		}
		i++
	}

	return true
}

func (e *Evaluator) matchSpecificTempo(track rekordbox.Track, c Comparison) bool {
	val := strings.ReplaceAll(c.Value, " ", "")
	parts := strings.Split(val, ":")
	if len(parts) == 0 {
		return false
	}

	idx, err := strconv.Atoi(parts[0])
	if err != nil || idx < 1 || idx > len(track.Tempo) {
		return false
	}

	tempo := track.Tempo[idx-1]

	if len(parts) == 1 {
		return true
	}

	return e.matchTempoProperties(tempo, parts[1:])
}

func (e *Evaluator) matchTempoProperties(t rekordbox.Tempo, props []string) bool {
	i := 0
	for i < len(props) {
		prop := strings.ToLower(props[i])
		matched := false

		switch prop {
		case "bpm", "tempo":
			if i+1 < len(props) {
				val := props[i+1]
				i++
				matched = e.matchNumericProperty(t.Bpm, val)
			}
		case "inizio", "start", "time":
			if i+1 < len(props) {
				val := props[i+1]
				i++
				matched = e.matchNumericProperty(t.Inizio, val)
			}
		case "metro", "meter", "signature":
			if i+1 < len(props) {
				val := props[i+1]
				i++
				matched = strings.EqualFold(t.Metro, val) || strings.Contains(t.Metro, val)
			}
		case "battito", "beat":
			if i+1 < len(props) {
				val := props[i+1]
				i++
				matched = e.matchNumericProperty(strconv.Itoa(int(t.Battito)), val)
			}
		}

		if !matched {
			return false
		}
		i++
	}
	return true
}

func (e *Evaluator) matchNumericProperty(fieldValue, targetValue string) bool {
	// Properties don't have operators in the string "start:10.5".
	// We handle ranges (10..20) or exact equality.
	if strings.Contains(targetValue, "..") {
		return e.matchRange(fieldValue, targetValue)
	}
	fv, errF := strconv.ParseFloat(fieldValue, 64)
	tv, errT := strconv.ParseFloat(targetValue, 64)
	if errF == nil && errT == nil {
		return fv == tv
	}
	return strings.EqualFold(fieldValue, targetValue)
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
	if pm.Num == -1 {
		return false
	}
	return e.matchHotCueColor(pm, color)
}

func (e *Evaluator) matchHotCueColor(pm rekordbox.PositionMark, color string) bool {
	c := strings.ToLower(color)
	switch {
	case c == "hotpink":
		return pm.Red == 222 && pm.Green == 68 && pm.Blue == 207
	case c == "purple":
		return pm.Red == 180 && pm.Green == 50 && pm.Blue == 255
	case c == "violet":
		return pm.Red == 170 && pm.Green == 114 && pm.Blue == 255
	case c == "indigo":
		return pm.Red == 100 && pm.Green == 115 && pm.Blue == 255
	case c == "blue":
		return pm.Red == 48 && pm.Green == 90 && pm.Blue == 255
	case c == "skyblue":
		return pm.Red == 80 && pm.Green == 180 && pm.Blue == 255
	case c == "aqua":
		return pm.Red == 0 && pm.Green == 224 && pm.Blue == 255
	case c == "darkgreen":
		return pm.Red == 31 && pm.Green == 163 && pm.Blue == 146
	case c == "brightgreen":
		return pm.Red == 16 && pm.Green == 177 && pm.Blue == 118
	case c == "green":
		return pm.Red == 40 && pm.Green == 226 && pm.Blue == 20
	case c == "yellowgreen":
		return pm.Red == 165 && pm.Green == 225 && pm.Blue == 22
	case c == "yellow":
		return pm.Red == 180 && pm.Green == 190 && pm.Blue == 4
	case c == "orange":
		return pm.Red == 195 && pm.Green == 175 && pm.Blue == 4
	case c == "darkorange":
		return pm.Red == 224 && pm.Green == 100 && pm.Blue == 27
	case c == "red":
		return pm.Red == 230 && pm.Green == 40 && pm.Blue == 40
	case c == "pink":
		return pm.Red == 255 && pm.Green == 18 && pm.Blue == 123
	case c == NullValue:
		return pm.Red == 0 && pm.Green == 0 && pm.Blue == 0
	}

	return e.matchRawRGB(pm, color)
}

func (e *Evaluator) matchRawRGB(pm rekordbox.PositionMark, color string) bool {
	// Extensions: allow querying by exact RGB values if color name doesn't match
	// Format: red:255, green:0, blue:0
	if strings.Contains(color, ",") {
		parts := strings.Split(color, ",")
		var r, g, b int
		for _, p := range parts {
			kv := strings.Split(strings.TrimSpace(p), ":")
			if len(kv) != 2 {
				continue
			}
			val, _ := strconv.Atoi(kv[1])
			switch strings.TrimSpace(kv[0]) {
			case "red", "r":
				r = val
			case "green", "g":
				g = val
			case "blue", "b":
				b = val
			}
		}
		return int(pm.Red) == r && int(pm.Green) == g && int(pm.Blue) == b
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
	case "title", "name":
		return node.Name
	case "parent", "folder":
		return parentFolder
	case "entries":
		return strconv.Itoa(int(rekordbox.DerefInt32(node.Entries)))
	case "count":
		return strconv.Itoa(int(rekordbox.DerefInt32(node.Count)))
	case "type":
		return strconv.Itoa(int(node.Type))
	}
	return ""
}
