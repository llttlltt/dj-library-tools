package query

import "fmt"

// Operator defines the type of match to perform
type Operator string

const (
	OpSubstring Operator = ":"  // artist:Four
	OpExact     Operator = "="  // artist=Four Tet
	OpRegex     Operator = "::" // artist::"^Four"
	OpRange     Operator = ".." // bpm:120..140
	OpGt        Operator = ">"  // rating:>3
	OpGte       Operator = ">=" // bpm:>=128
	OpLt        Operator = "<"  // bpm:<100
	OpLte       Operator = "<=" // bitrate:<=320
	OpNeq       Operator = "!=" // genre:!=House
)

// Expression is the interface for all nodes in the query tree
type Expression interface {
	isExpression()
}

// Comparison represents a single filter (e.g., artist:Four)
type Comparison struct {
	Field    string
	Operator Operator
	Value    string
}

func (c Comparison) isExpression() {}

// Logical represents a boolean operation (AND, OR)
type Logical struct {
	Op    string // "AND", "OR"
	Left  Expression
	Right Expression
}

func (l Logical) isExpression() {}

// Not represents a negation
type Not struct {
	Expr Expression
}

func (n Not) isExpression() {}

// Query is the top-level container
type Query struct {
	Root Expression
}

// Validate checks if the query is valid and returns a helpful error if not.
func (q Query) Validate() error {
	if q.Root == nil {
		return nil
	}
	return q.validateExpr(q.Root)
}

func (q Query) validateExpr(expr Expression) error {
	switch v := expr.(type) {
	case Comparison:
		if v.Field == "" {
			return fmt.Errorf("query must specify a field (e.g. title:%q). Bare values are not supported", v.Value)
		}
	case Logical:
		if err := q.validateExpr(v.Left); err != nil {
			return err
		}
		return q.validateExpr(v.Right)
	case Not:
		return q.validateExpr(v.Expr)
	}
	return nil
}
