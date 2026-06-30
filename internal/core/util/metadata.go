package util

// NormalizeRating scales a rating from a source range (e.g. 0-5) to our 0-255 standard.
func NormalizeRating(val float64, max float64) int {
	if max == 0 {
		return 0
	}
	return int((val / max) * 255)
}
