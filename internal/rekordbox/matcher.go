package rekordbox

import (
	"fmt"
	"strings"

	"github.com/llttlltt/dj-library-tools/internal/models"
	"github.com/llttlltt/dj-library-tools/internal/query"
)

// CustomMatch implements logic for Rekordbox-specific query fields (hotcues, memorycues).
func CustomMatch(track models.Track, field string, op query.Operator, value string) bool {
	target := strings.ToLower(value)

	if field == "hotcues" {
		for _, cp := range track.CuePoints {
			if cp.Type == models.CueTypeMemory { continue }
			if matchCuePoint(cp, target, op) { return true }
		}
	} else if field == "memorycues" {
		for _, cp := range track.CuePoints {
			if cp.Type != models.CueTypeMemory { continue }
			if matchCuePoint(cp, target, op) { return true }
		}
	}
	return false
}

func matchCuePoint(cp models.CuePoint, target string, op query.Operator) bool {
	if op == query.OpExact {
		if strings.EqualFold(cp.Name, target) { return true }
	} else if strings.Contains(strings.ToLower(cp.Name), target) {
		return true
	}

	colorName := strings.ToLower(cp.Color)
	if op == query.OpExact {
		if colorName == target { return true }
	} else if strings.Contains(colorName, target) {
		return true
	}

	return false
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
