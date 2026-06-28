package query

import (
	"fmt"
)

// StatFunction defines a calculation performed on a slice of values.
type StatFunction func(values []string) string

// Stats maps stat names (without the hyphen) to their implementation.
var Stats = map[string]StatFunction{
	"count": func(values []string) string {
		return fmt.Sprintf("%d", len(values))
	},
	"drift": func(values []string) string {
		if len(values) == 0 {
			return "0"
		}
		var min, max float64
		first := true
		for _, v := range values {
			f := parseToFloat(v)
			if first {
				min, max = f, f
				first = false
				continue
			}
			if f < min {
				min = f
			}
			if f > max {
				max = f
			}
		}
		return fmt.Sprintf("%.4f", max-min)
	},
}

// DensityStat calculates markers per minute.
// This is a special stat that needs the track duration.
func DensityStat(count int, durationSeconds int) string {
	if durationSeconds <= 0 {
		return "0"
	}
	minutes := float64(durationSeconds) / 60.0
	density := float64(count) / minutes
	return fmt.Sprintf("%.2f", density)
}
