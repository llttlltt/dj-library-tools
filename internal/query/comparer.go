package query

import (
	"regexp"
	"strconv"
	"strings"
)

// FieldType defines the data type of a field for evaluation purposes.
type FieldType int

const (
	TypeString FieldType = iota
	TypeNumeric
)

// Schema maps field names to their types.
var Schema = map[string]FieldType{
	"playlists":  TypeNumeric,
	"hotcues":    TypeNumeric,
	"memorycues": TypeNumeric,
	"beatgrids":  TypeNumeric,
	"rating":     TypeNumeric,
	"plays":      TypeNumeric,
	"year":       TypeNumeric,
	"bpm":        TypeNumeric,
	"bitrate":    TypeNumeric,
	"samplerate": TypeNumeric,
	"size":       TypeNumeric,
	"items":      TypeNumeric,
	"duration":   TypeNumeric,
}

// Compare executes a comparison between two values based on the operator and field type.
func Compare(field string, fieldValue, targetValue string, op Operator) bool {
	if op == OpRange {
		return matchRange(fieldValue, targetValue)
	}

	fieldType := Schema[strings.ToLower(field)]
	if fieldType == TypeNumeric {
		return matchNumeric(fieldValue, targetValue, op)
	}

	return matchString(fieldValue, targetValue, op)
}

func matchRange(fieldValue, rangeValue string) bool {
	parts := strings.Split(rangeValue, "..")
	if len(parts) != 2 { return false }
	val, _ := strconv.ParseFloat(fieldValue, 64)
	min, _ := strconv.ParseFloat(parts[0], 64)
	max, _ := strconv.ParseFloat(parts[1], 64)
	return val >= min && val <= max
}

func matchNumeric(fieldValue, targetValue string, op Operator) bool {
	f, _ := strconv.ParseFloat(fieldValue, 64)
	t, _ := strconv.ParseFloat(targetValue, 64)
	switch op {
	case OpGt:  return f > t
	case OpGte: return f >= t
	case OpLt:  return f < t
	case OpLte: return f <= t
	case OpExact, OpSubstring: return f == t
	}
	return false
}

func matchString(fieldValue, targetValue string, op Operator) bool {
	switch op {
	case OpExact:
		return strings.EqualFold(fieldValue, targetValue)
	case OpSubstring:
		return strings.Contains(strings.ToLower(fieldValue), strings.ToLower(targetValue))
	case OpRegex:
		re, err := regexp.Compile(targetValue)
		if err != nil { return false }
		return re.MatchString(fieldValue)
	}
	return false
}
