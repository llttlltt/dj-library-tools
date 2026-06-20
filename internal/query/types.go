package query

// Operator defines the type of match to perform
type Operator string

const (
	OpSubstring  Operator = ":"  // artist:Four
	OpExact      Operator = "="  // artist=Four Tet
	OpRegex      Operator = "::" // artist::"^Four"
	OpRange      Operator = ".." // bpm:120..140
)

// Criterion represents a single filter in a query
type Criterion struct {
	Field    string
	Operator Operator
	Value    string
}

// Query is a collection of criteria that must all be met (AND logic)
type Query struct {
	Criteria []Criterion
	Negated  bool // !artist:Four Tet
}
