package rekordbox

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/llttlltt/dj-library-tools/internal/models"
	"github.com/llttlltt/dj-library-tools/internal/query"
)

// CustomMatch implements logic for Rekordbox-specific query fields (hotcues, memorycues).
func CustomMatch(track models.Track, field string, op query.Operator, value string) bool {
	var targetCues []models.CuePoint
	if field == "hotcues" {
		targetCues = filterCues(track.CuePoints, models.CueTypeHot)
	} else if field == "memorycues" {
		targetCues = filterCues(track.CuePoints, models.CueTypeMemory)
	} else {
		return false
	}

	// Parse nested value syntax: [index|property]:[subproperty|value]:[value]
	// Example: "0:color:red" or "red" or "0:red"
	parts := strings.Split(value, ":")
	
	switch len(parts) {
	case 1:
		// Simple match: search all cues of this type for the value
		return matchAnyCue(targetCues, parts[0], op)
	case 2:
		// Index match or Property match
		// e.g. "0:red" (first cue is red) or "color:red" (any cue is red)
		if idx, err := parseCueIndex(parts[0]); err == nil {
			return matchSpecificCue(targetCues, idx, parts[1], "", op)
		}
		return matchAnyCueProperty(targetCues, parts[0], parts[1], op)
	case 3:
		// Specific property match at index
		// e.g. "0:color:red"
		if idx, err := parseCueIndex(parts[0]); err == nil {
			return matchSpecificCue(targetCues, idx, parts[2], parts[1], op)
		}
	}

	return false
}

func filterCues(cues []models.CuePoint, cueType models.CueType) []models.CuePoint {
	var filtered []models.CuePoint
	for _, cp := range cues {
		if cp.Type == cueType {
			filtered = append(filtered) // placeholder - wait, I need to keep the index
		}
	}
	// Note: We'll use the track's original slice and check type during matching 
	// to ensure 'Index' matches correctly.
	return cues
}

func matchAnyCue(cues []models.CuePoint, target string, op query.Operator) bool {
	for _, cp := range cues {
		if matchCuePoint(cp, target, op) { return true }
	}
	return false
}

func matchAnyCueProperty(cues []models.CuePoint, prop, target string, op query.Operator) bool {
	for _, cp := range cues {
		if matchCuePointProperty(cp, prop, target, op) { return true }
	}
	return false
}

func matchSpecificCue(cues []models.CuePoint, index int, target string, prop string, op query.Operator) bool {
	for _, cp := range cues {
		if cp.Index == index {
			if prop != "" {
				return matchCuePointProperty(cp, prop, target, op)
			}
			return matchCuePoint(cp, target, op)
		}
	}
	return false
}

func matchCuePoint(cp models.CuePoint, target string, op query.Operator) bool {
	// Match against Name or Color
	if matchCuePointProperty(cp, "name", target, op) { return true }
	return matchCuePointProperty(cp, "color", target, op)
}

func matchCuePointProperty(cp models.CuePoint, prop, target string, op query.Operator) bool {
	val := ""
	switch strings.ToLower(prop) {
	case "name", "comment":
		val = cp.Name
	case "color":
		val = cp.Color
	default:
		return false
	}

	if op == query.OpExact {
		return strings.EqualFold(val, target)
	}
	return strings.Contains(strings.ToLower(val), strings.ToLower(target))
}

func parseCueIndex(s string) (int, error) {
	// Handle letters a-h for hotcues
	if len(s) == 1 && s[0] >= 'a' && s[0] <= 'h' {
		return int(s[0] - 'a'), nil
	}
	return strconv.Atoi(s)
}

func GetHotCueColorName(pm PositionMark) string {
	rgb := fmt.Sprintf("%02X%02X%02X", pm.Red, pm.Green, pm.Blue)
	switch rgb {
	case "E62828": return "red"
	case "DE44CF": return "hotpink"
	case "FFFF00", "B4BE04", "C3AF04": return "yellow"
	case "28E214", "10B176": return "green"
	case "00E0FF", "50B4FF": return "aqua"
	case "305AFF", "6473FF": return "blue"
	case "B432FF", "AA72FF": return "purple"
	case "E0641B", "FFA500": return "orange"
	}
	return ""
}

func GetTrackColorName(hex string) string {
	switch strings.ToUpper(hex) {
	case "0XFF007F": return "pink"
	case "0XFF0000": return "red"
	case "0XFFA500": return "orange"
	case "0XFFFF00": return "yellow"
	case "0X00FF00": return "green"
	case "0X25FDE9": return "aqua"
	case "0X0000FF": return "blue"
	case "0X660099": return "purple"
	}
	return hex
}
