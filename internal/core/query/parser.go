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
	Kind   TokenKind
	Value  string
	Quoted bool
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
	// First split by : if present to separate field from value
	sepIdx := strings.IndexAny(word, ":=")
	if sepIdx == -1 {
		return []Token{{Kind: TokenValue, Value: word}}
	}

	field := word[:sepIdx]
	rest := word[sepIdx:] // starts with : or =

	// Ensure we don't treat quoted values inside the word as part of the operator
	if strings.ContainsAny(field, "\"'") {
		return []Token{{Kind: TokenValue, Value: word}}
	}

	op := ":"
	val := rest[1:]

	switch {
	case strings.HasPrefix(rest, ":="):
		op = ":="
		val = rest[2:]
	case strings.HasPrefix(rest, "::"):
		op = "::"
		val = rest[2:]
	case strings.HasPrefix(rest, "!="):
		op = "!="
		val = rest[2:]
	case strings.HasPrefix(rest, "=="):
		op = ":="
		val = rest[2:]
	case strings.HasPrefix(rest, ":>="):
		op = ">="
		val = strings.TrimPrefix(rest[1:], ">=")
	case strings.HasPrefix(rest, ":<="):
		op = "<="
		val = strings.TrimPrefix(rest[1:], "<=")
	case strings.HasPrefix(rest, ":>"):
		op = ">"
		val = strings.TrimPrefix(rest[1:], ">")
	case strings.HasPrefix(rest, ":<"):
		op = "<"
		val = strings.TrimPrefix(rest[1:], "<")
	case strings.HasPrefix(rest, "gte") || strings.HasPrefix(rest, ":gte"):
		op = ">="
		val = strings.TrimPrefix(strings.TrimPrefix(rest, ":"), "gte")
	case strings.HasPrefix(rest, "lte") || strings.HasPrefix(rest, ":lte"):
		op = "<="
		val = strings.TrimPrefix(strings.TrimPrefix(rest, ":"), "lte")
	case strings.HasPrefix(rest, "gt") || strings.HasPrefix(rest, ":gt"):
		op = ">"
		val = strings.TrimPrefix(strings.TrimPrefix(rest, ":"), "gt")
	case strings.HasPrefix(rest, "lt") || strings.HasPrefix(rest, ":lt"):
		op = "<"
		val = strings.TrimPrefix(strings.TrimPrefix(rest, ":"), "lt")
	case strings.HasPrefix(rest, "neq") || strings.HasPrefix(rest, ":neq"):
		op = "!="
		val = strings.TrimPrefix(strings.TrimPrefix(rest, ":"), "neq")
	case strings.HasPrefix(rest, "ne") || strings.HasPrefix(rest, ":ne"):
		op = "!="
		val = strings.TrimPrefix(strings.TrimPrefix(rest, ":"), "ne")
	case strings.HasPrefix(rest, "eq") || strings.HasPrefix(rest, ":eq"):
		op = ":="
		val = strings.TrimPrefix(strings.TrimPrefix(rest, ":"), "eq")
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
		case '"', '\'':
			quote := r
			val, end := p.readQuoted(runes, i, quote)
			// Check if this was preceded by a field name (e.g., name:"Value")
			if len(tokens) > 0 && tokens[len(tokens)-1].Kind == TokenField && tokens[len(tokens)-1].Value != "" {
				// The field is already there, just add the value
				tokens = append(tokens, Token{Kind: TokenValue, Value: val, Quoted: true})
			} else if i > 0 && runes[i-1] == ':' {
				// This is the value part of a field:value pair where the colon was already tokenized
				tokens = append(tokens, Token{Kind: TokenValue, Value: val, Quoted: true})
			} else {
				// Bare quoted value
				tokens = append(tokens, Token{Kind: TokenValue, Value: val, Quoted: true})
			}
			i = end
		case '!', '-':
			// Look ahead: only treat as NOT if it's start of token AND not preceded by an operator
			// This prevents - inside field values (like name:Sorting/1 - Inbox) from being seen as NOT.
			isStartOfToken := i == 0 || unicode.IsSpace(runes[i-1]) || runes[i-1] == '('
			if isStartOfToken && i+1 < len(runes) && !unicode.IsSpace(runes[i+1]) && runes[i+1] != '=' {
				tokens = append(tokens, Token{Kind: TokenNot, Value: "!"})
			} else {
				word, end := p.readWord(runes, i)
				tokens = append(tokens, Token{Kind: TokenValue, Value: word})
				i = end
			}
		default:
			word, end := p.readWord(runes, i)
			wordLower := strings.ToLower(word)
			if wordLower == "and" {
				tokens = append(tokens, Token{Kind: TokenAnd, Value: "AND"})
			} else if wordLower == "or" {
				tokens = append(tokens, Token{Kind: TokenOr, Value: "OR"})
			} else if wordLower == "not" {
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

func (p *Parser) readQuoted(runes []rune, start int, quote rune) (string, int) {
	var sb strings.Builder
	for i := start + 1; i < len(runes); i++ {
		if runes[i] == quote {
			// If followed by another quote, it's an escaped quote
			if i+1 < len(runes) && runes[i+1] == quote {
				sb.WriteRune(quote)
				i++
				continue
			}
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
		// Special handling for operators: keep them separate if they are syntactically meaningful.
		// We allow ' inside words for contractions, but it still functions as a quote if it starts the word.
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
	return kind == TokenValue || kind == TokenLParen || kind == TokenNot || kind == TokenField
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
		case ":=":
			op = OpExact
		case "!=":
			op = OpNeq
		case "::":
			op = OpRegex
		case "..":
			op = OpRange
		case ">":
			op = OpGt
		case ">=":
			op = OpGte
		case "<":
			op = OpLt
		case "<=":
			op = OpLte
		}

		val := ""
		quoted := false
		if p.peek().Kind == TokenValue {
			tok := p.next()
			val = tok.Value
			quoted = tok.Quoted
		}

		for {
			peek := p.peek()
			if peek.Kind != TokenValue {
				break
			}
			if val == "" || val == `""` {
				val = peek.Value
				quoted = peek.Quoted
			} else {
				val += " " + peek.Value
				if peek.Quoted {
					quoted = true
				}
			}
			p.next()
		}

		if strings.Contains(val, "..") {
			op = OpRange
		}

		if val == `""` {
			val = ""
		}

		// Trim surrounding quotes if present
		if len(val) >= 2 {
			if (val[0] == '"' && val[len(val)-1] == '"') || (val[0] == '\'' && val[len(val)-1] == '\'') {
				val = val[1 : len(val)-1]
			}
		}

		return Comparison{Field: token.Value, Operator: op, Value: val, Quoted: quoted}
	case TokenValue:
		val := token.Value
		if len(val) >= 2 {
			if (val[0] == '"' && val[len(val)-1] == '"') || (val[0] == '\'' && val[len(val)-1] == '\'') {
				val = val[1 : len(val)-1]
			}
		}
		return Comparison{Field: "", Operator: OpSubstring, Value: val, Quoted: token.Quoted}
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
