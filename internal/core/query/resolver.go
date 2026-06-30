package query

import (
	"strconv"
	"strings"
	"time"
)

// ResolveValue normalizes a query value based on the field type.
func ResolveValue(field string, value string) string {
	field = strings.ToLower(field)

	// Date resolution
	if field == "added" || field == "modified" {
		return resolveDateShorthand(value)
	}

	return value
}

func resolveDateShorthand(val string) string {
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
		case "d":
			return now.AddDate(0, 0, -amount).Format("2006-01-02")
		case "m":
			return now.AddDate(0, -amount, 0).Format("2006-01-02")
		case "y":
			return now.AddDate(-amount, 0, 0).Format("2006-01-02")
		}
	}
	return val
}
