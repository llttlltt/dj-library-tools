package query

import (
	"fmt"
	"math"
)

// StatFunction defines a calculation performed on a slice of values.
type StatFunction func(values []string) string

// Stats maps stat names (without the hyphen) to their implementation.
var Stats = map[string]StatFunction{
	"count": func(values []string) string {
		return fmt.Sprintf("%d", len(values))
	},
	"min": func(values []string) string {
		if len(values) == 0 {
			return "0"
		}
		min := parseToFloat(values[0])
		for _, v := range values[1:] {
			f := parseToFloat(v)
			if f < min {
				min = f
			}
		}
		return fmt.Sprintf("%.4f", min)
	},
	"max": func(values []string) string {
		if len(values) == 0 {
			return "0"
		}
		max := parseToFloat(values[0])
		for _, v := range values[1:] {
			f := parseToFloat(v)
			if f > max {
				max = f
			}
		}
		return fmt.Sprintf("%.4f", max)
	},
	"avg": func(values []string) string {
		if len(values) == 0 {
			return "0"
		}
		sum := 0.0
		for _, v := range values {
			sum += parseToFloat(v)
		}
		return fmt.Sprintf("%.4f", sum/float64(len(values)))
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
	"jitter": func(values []string) string {
		if len(values) < 2 {
			return "0"
		}
		diffSum := 0.0
		prev := parseToFloat(values[0])
		for _, v := range values[1:] {
			curr := parseToFloat(v)
			diffSum += math.Abs(curr - prev)
			prev = curr
		}
		return fmt.Sprintf("%.4f", diffSum/float64(len(values)-1))
	},
	"redundancy": func(values []string) string {
		if len(values) < 2 {
			return "0"
		}
		matches := 0
		prev := parseToFloat(values[0])
		for _, v := range values[1:] {
			curr := parseToFloat(v)
			if curr == prev {
				matches++
			}
			prev = curr
		}
		return fmt.Sprintf("%.0f", (float64(matches)/float64(len(values)-1))*100.0)
	},
	"stability": func(values []string) string {
		if len(values) < 2 {
			return "100"
		}
		// Composite: 100 - (drift * 5) - (jitter * 20)
		// We'll cap it at 0-100.
		var min, max float64
		diffSum := 0.0
		first := true
		prev := 0.0
		for _, v := range values {
			f := parseToFloat(v)
			if first {
				min, max, prev = f, f, f
				first = false
				continue
			}
			if f < min { min = f }
			if f > max { max = f }
			diffSum += math.Abs(f - prev)
			prev = f
		}
		drift := max - min
		jitter := diffSum / float64(len(values)-1)
		
		score := 100.0 - (drift * 10.0) - (jitter * 50.0)
		if score < 0 { score = 0 }
		return fmt.Sprintf("%.0f", score)
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
