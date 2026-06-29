package query

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/llttlltt/dj-library-tools/internal/models"
)

// Schema maps field names to their types for the query engine.
// It is derived from the universal models.TrackFields and models.GroupFields.
func getFieldKind(field string) models.FieldKind {
	field = strings.ToLower(field)

	// Handle path-based fields
	if strings.ContainsAny(field, "./-") {
		collection, _, property, stat := ParsePath(field)
		
		// 1. Stats are numeric (except possibly custom ones, but currently all are)
		if stat != "" {
			return models.KindNumeric
		}

		// 2. Look up property in collection definition
		if fields, ok := models.CollectionFields[collection]; ok {
			if kind, ok := fields[property]; ok {
				return kind
			}
		}

		// Default for collections (e.g. hotcues.1 is usually treated as a name)
		return models.KindString
	}

	if def, ok := models.TrackFields[field]; ok {
		return def.Kind
	}
	if def, ok := models.GroupFields[field]; ok {
		return def.Kind
	}
	return models.KindString
}

// Compare executes a comparison between two values based on the operator and field type.
func Compare(field string, fieldValue, targetValue string, op Operator) bool {
	if op == OpRange {
		return matchRange(fieldValue, targetValue)
	}

	if getFieldKind(field) == models.KindNumeric {
		return matchNumeric(fieldValue, targetValue, op)
	}

	return matchString(fieldValue, targetValue, op)
}

func matchRange(fieldValue, rangeValue string) bool {
	parts := strings.Split(rangeValue, "..")
	if len(parts) != 2 { return false }
	return matchNumeric(fieldValue, parts[0], OpGte) && matchNumeric(fieldValue, parts[1], OpLte)
}

func matchNumeric(fieldValue, targetValue string, op Operator) bool {
	f := parseToFloat(fieldValue)
	t := parseToFloat(targetValue)
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

func parseToFloat(s string) float64 {
	var f float64
	fmt.Sscanf(s, "%f", &f)
	return f
}
