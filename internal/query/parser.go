package query

import (
	"strings"
)

// Parser handles converting a query string into a Query object
type Parser struct{}

func NewParser() *Parser {
	return &Parser{}
}

// Parse converts a string like 'artist:"Four Tet" bpm:120..140' into a Query struct
func (p *Parser) Parse(input string) Query {
	q := Query{}

	if strings.HasPrefix(input, "!") {
		q.Negated = true
		input = strings.TrimPrefix(input, "!")
	}

	// If the input contains a colon but no spaces before it,
	// it might be a single multi-word criterion (like shell passed artist:MJ Cole)
	if strings.Contains(input, ":") && !strings.Contains(input[:strings.Index(input, ":")], " ") {
		if crit, ok := p.parsePart(input); ok {
			q.Criteria = append(q.Criteria, crit)
			return q
		}
	}

	parts := p.splitInput(input)
	for _, part := range parts {
		if criterion, ok := p.parsePart(part); ok {
			q.Criteria = append(q.Criteria, criterion)
		}
	}

	return q
}

// splitInput handles spaces but respects double quotes
func (p *Parser) splitInput(input string) []string {
	var parts []string
	var current strings.Builder
	inQuotes := false

	for i := 0; i < len(input); i++ {
		char := input[i]
		if char == '"' {
			inQuotes = !inQuotes
			continue
		}
		if char == ' ' && !inQuotes {
			if current.Len() > 0 {
				parts = append(parts, current.String())
				current.Reset()
			}
			continue
		}
		current.WriteByte(char)
	}
	if current.Len() > 0 {
		parts = append(parts, current.String())
	}
	return parts
}

func (p *Parser) parsePart(part string) (Criterion, bool) {
	// Standard format: field<operator>value
	// Operators: :: (regex), .. (range), = (exact), : (substring)

	// Check for Range first within the value part, but we need to find the field separator first
	sepIdx := strings.IndexAny(part, ":=")
	if sepIdx == -1 {
		// Default to substring match on 'name' if no separator is found
		return Criterion{
			Field:    "name",
			Operator: OpSubstring,
			Value:    part,
		}, true
	}

	field := part[:sepIdx]
	opChar := string(part[sepIdx])
	value := part[sepIdx+1:]

	// Determine specific operator
	op := OpSubstring
	if opChar == "=" {
		op = OpExact
	}

	// Check for double-colon (Regex)
	if opChar == ":" && strings.HasPrefix(value, ":") {
		op = OpRegex
		value = value[1:]
	}

	// Check if value contains range operator
	if strings.Contains(value, "..") {
		op = OpRange
	}

	return Criterion{
		Field:    strings.ToLower(field),
		Operator: op,
		Value:    value,
	}, true
}
