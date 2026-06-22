package query

import (
	"strings"
	"unicode"
)

type TokenKind int

const (
	TokenValue TokenKind = iota
	TokenField
	TokenOp
	TokenLParen
	TokenRParen
	TokenAnd
	TokenOr
	TokenNot
	TokenEOF
)

type Token struct {
	Kind  TokenKind
	Value string
}

type Parser struct {
	tokens []Token
	pos    int
}

func NewParser() *Parser {
	return &Parser{}
}

func (p *Parser) Parse(input string) Query {
	p.tokens = p.tokenize(input)
	p.pos = 0
	return Query{Root: p.parseExpression()}
}

func (p *Parser) parsePart(word string) []Token {
	// Split by Field separator first
	sepIdx := strings.IndexAny(word, ":=")
	if sepIdx == -1 {
		return []Token{{Kind: TokenValue, Value: word}}
	}

	field := word[:sepIdx]
	rest := word[sepIdx:] // starts with : or =
	
	op := ":"
	val := rest[1:]
	
	if strings.HasPrefix(rest, "::") {
		op = "::"
		val = rest[2:]
	} else if strings.HasPrefix(rest, ":>=") || strings.HasPrefix(rest, ">= ") {
		op = ">="
		val = strings.TrimPrefix(rest[1:], ">=")
	} else if strings.HasPrefix(rest, ":<=") || strings.HasPrefix(rest, "<= ") {
		op = "<="
		val = strings.TrimPrefix(rest[1:], "<=")
	} else if strings.HasPrefix(rest, ":>") || strings.HasPrefix(rest, "> ") {
		op = ">"
		val = strings.TrimPrefix(rest[1:], ">")
	} else if strings.HasPrefix(rest, ":<") || strings.HasPrefix(rest, "< ") {
		op = "<"
		val = strings.TrimPrefix(rest[1:], "<")
	} else if strings.HasPrefix(rest, "gte") || strings.HasPrefix(rest, ":gte") {
		op = ">="
		val = strings.TrimPrefix(strings.TrimPrefix(rest, ":"), "gte")
	} else if strings.HasPrefix(rest, "lte") || strings.HasPrefix(rest, ":lte") {
		op = "<="
		val = strings.TrimPrefix(strings.TrimPrefix(rest, ":"), "lte")
	} else if strings.HasPrefix(rest, "gt") || strings.HasPrefix(rest, ":gt") {
		op = ">"
		val = strings.TrimPrefix(strings.TrimPrefix(rest, ":"), "gt")
	} else if strings.HasPrefix(rest, "lt") || strings.HasPrefix(rest, ":lt") {
		op = "<"
		val = strings.TrimPrefix(strings.TrimPrefix(rest, ":"), "lt")
	} else if strings.HasPrefix(rest, "=") {
		op = "="
		val = rest[1:]
	}

	if strings.Contains(val, "..") {
		op = ".."
	}

	return []Token{
		{Kind: TokenField, Value: field},
		{Kind: TokenOp, Value: op},
		{Kind: TokenValue, Value: val},
	}
}

func (p *Parser) tokenize(input string) []Token {
	var tokens []Token
	runes := []rune(input)
	for i := 0; i < len(runes); i++ {
		r := runes[i]
		if unicode.IsSpace(r) {
			continue
		}

		switch r {
		case '(':
			tokens = append(tokens, Token{Kind: TokenLParen, Value: "("})
		case ')':
			tokens = append(tokens, Token{Kind: TokenRParen, Value: ")"})
		case '!':
			tokens = append(tokens, Token{Kind: TokenNot, Value: "!"})
		case '&':
			if i+1 < len(runes) && runes[i+1] == '&' {
				tokens = append(tokens, Token{Kind: TokenAnd, Value: "&&"})
				i++
			}
		case '|':
			if i+1 < len(runes) && runes[i+1] == '|' {
				tokens = append(tokens, Token{Kind: TokenOr, Value: "||"})
				i++
			}
		case '"':
			val, end := p.readQuoted(runes, i)
			if val == "" {
				tokens = append(tokens, Token{Kind: TokenValue, Value: `""`})
			} else {
				tokens = append(tokens, Token{Kind: TokenValue, Value: val})
			}
			i = end
		default:
			word, end := p.readWord(runes, i)
			if strings.EqualFold(word, "AND") {
				tokens = append(tokens, Token{Kind: TokenAnd, Value: "AND"})
			} else if strings.EqualFold(word, "OR") {
				tokens = append(tokens, Token{Kind: TokenOr, Value: "OR"})
			} else if strings.EqualFold(word, "NOT") {
				tokens = append(tokens, Token{Kind: TokenNot, Value: "!"})
			} else if strings.ContainsAny(word, ":=") {
				tokens = append(tokens, p.parsePart(word)...)
			} else {
				tokens = append(tokens, Token{Kind: TokenValue, Value: word})
			}
			i = end
		}
	}
	return tokens
}

func (p *Parser) readQuoted(runes []rune, start int) (string, int) {
	var sb strings.Builder
	for i := start + 1; i < len(runes); i++ {
		if runes[i] == '"' {
			return sb.String(), i
		}
		sb.WriteRune(runes[i])
	}
	return sb.String(), len(runes)
}

func (p *Parser) readWord(runes []rune, start int) (string, int) {
	var sb strings.Builder
	i := start
	for ; i < len(runes); i++ {
		r := runes[i]
		if unicode.IsSpace(r) || r == '(' || r == ')' || r == '!' || r == '&' || r == '|' || r == '"' {
			break
		}
		sb.WriteRune(r)
	}
	return sb.String(), i - 1
}

func (p *Parser) parseExpression() Expression {
	return p.parseOr()
}

func (p *Parser) parseOr() Expression {
	left := p.parseAnd()
	if left == nil {
		return nil
	}
	for p.peek().Kind == TokenOr {
		p.next()
		right := p.parseAnd()
		if right != nil {
			left = Logical{Op: "OR", Left: left, Right: right}
		}
	}
	return left
}

func (p *Parser) parseAnd() Expression {
	left := p.parsePrimary()
	if left == nil {
		return nil
	}
	for {
		kind := p.peek().Kind
		if kind == TokenAnd {
			p.next()
			right := p.parsePrimary()
			if right != nil {
				left = Logical{Op: "AND", Left: left, Right: right}
			}
		} else if p.isNextImplicitAnd() {
			right := p.parsePrimary()
			if right != nil {
				left = Logical{Op: "AND", Left: left, Right: right}
			} else {
				break
			}
		} else {
			break
		}
	}
	return left
}

func (p *Parser) isNextImplicitAnd() bool {
	kind := p.peek().Kind
	return kind == TokenValue || kind == TokenLParen || kind == TokenNot
}

func (p *Parser) parsePrimary() Expression {
	token := p.next()
	if token.Kind == TokenEOF {
		return nil
	}
	
	switch token.Kind {
	case TokenNot:
		return Not{Expr: p.parsePrimary()}
	case TokenLParen:
		expr := p.parseExpression()
		if p.peek().Kind == TokenRParen {
			p.next() // consume ')'
		}
		return expr
	case TokenField:
		opToken := p.next()
		
		op := OpSubstring
		switch opToken.Value {
		case "=": op = OpExact
		case "::": op = OpRegex
		case "..": op = OpRange
		case ">": op = OpGt
		case ">=": op = OpGte
		case "<": op = OpLt
		case "<=": op = OpLte
		}
		
		val := ""
		if p.peek().Kind == TokenValue {
			val = p.next().Value
		}
		
		isCueField := strings.HasPrefix(strings.ToLower(token.Value), "hotcue") || strings.HasPrefix(strings.ToLower(token.Value), "memorycue")

		for {
			peek := p.peek()
			if peek.Kind != TokenValue {
				break
			}
			if val == "" || val == `""` {
				val = peek.Value
			} else if isCueField {
				val += ":" + peek.Value
			} else {
				val += " " + peek.Value
			}
			p.next()
		}

		if strings.Contains(val, "..") {
			op = OpRange
		}
		
		if val == `""` {
			val = ""
		}
		
		return Comparison{Field: token.Value, Operator: op, Value: val}
	case TokenValue:
		val := token.Value
		if val == `""` {
			val = ""
		}
		return Comparison{Field: "name", Operator: OpSubstring, Value: val}
	}
	return nil
}

func (p *Parser) peek() Token {
	if p.pos >= len(p.tokens) {
		return Token{Kind: TokenEOF}
	}
	return p.tokens[p.pos]
}

func (p *Parser) next() Token {
	token := p.peek()
	p.pos++
	return token
}
