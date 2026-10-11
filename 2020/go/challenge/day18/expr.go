package day18

import (
	"fmt"
	"strconv"
)

type tokenKind int

const (
	tokenNum tokenKind = iota
	tokenPlus
	tokenStar
	tokenLParen
	tokenRParen
)

type token struct {
	kind tokenKind
	num  int
}

func tokenize(line string) ([]token, error) {
	var tokens []token
	i := 0
	for i < len(line) {
		c := line[i]
		switch {
		case c == ' ':
			i++
		case c == '+':
			tokens = append(tokens, token{kind: tokenPlus})
			i++
		case c == '*':
			tokens = append(tokens, token{kind: tokenStar})
			i++
		case c == '(':
			tokens = append(tokens, token{kind: tokenLParen})
			i++
		case c == ')':
			tokens = append(tokens, token{kind: tokenRParen})
			i++
		case c >= '0' && c <= '9':
			start := i
			for i < len(line) && line[i] >= '0' && line[i] <= '9' {
				i++
			}
			n, err := strconv.Atoi(line[start:i])
			if err != nil {
				return nil, fmt.Errorf("invalid number %q in %q: %w", line[start:i], line, err)
			}
			tokens = append(tokens, token{kind: tokenNum, num: n})
		default:
			return nil, fmt.Errorf("unexpected character %q in %q", c, line)
		}
	}
	if len(tokens) == 0 {
		return nil, fmt.Errorf("empty expression in %q", line)
	}
	return tokens, nil
}

// parser is a recursive-descent evaluator. When additionBindsTighter is set,
// '+' binds more tightly than '*' (part 2); otherwise both operators have
// equal precedence and are applied left to right (part 1).
type parser struct {
	tokens               []token
	pos                  int
	additionBindsTighter bool
}

func (p *parser) peekKind() (tokenKind, bool) {
	if p.pos >= len(p.tokens) {
		return 0, false
	}
	return p.tokens[p.pos].kind, true
}

func (p *parser) nextToken() (token, error) {
	if p.pos >= len(p.tokens) {
		return token{}, fmt.Errorf("unexpected end of expression, expected a token")
	}
	t := p.tokens[p.pos]
	p.pos++
	return t, nil
}

func (p *parser) parseExpr() (int, error) {
	if p.additionBindsTighter {
		value, err := p.parseSum()
		if err != nil {
			return 0, err
		}
		for {
			kind, ok := p.peekKind()
			if !ok || kind != tokenStar {
				break
			}
			p.pos++
			rhs, err := p.parseSum()
			if err != nil {
				return 0, err
			}
			value *= rhs
		}
		return value, nil
	}
	return p.parseFlat()
}

// parseSum parses a chain of '+' only, used when addition binds tighter.
func (p *parser) parseSum() (int, error) {
	value, err := p.parseFactor()
	if err != nil {
		return 0, err
	}
	for {
		kind, ok := p.peekKind()
		if !ok || kind != tokenPlus {
			break
		}
		p.pos++
		rhs, err := p.parseFactor()
		if err != nil {
			return 0, err
		}
		value += rhs
	}
	return value, nil
}

// parseFlat parses a left-to-right chain of '+'/'*' with equal precedence,
// used for part 1.
func (p *parser) parseFlat() (int, error) {
	value, err := p.parseFactor()
	if err != nil {
		return 0, err
	}
	for {
		kind, ok := p.peekKind()
		if !ok {
			break
		}
		switch kind {
		case tokenPlus:
			p.pos++
			rhs, err := p.parseFactor()
			if err != nil {
				return 0, err
			}
			value += rhs
		case tokenStar:
			p.pos++
			rhs, err := p.parseFactor()
			if err != nil {
				return 0, err
			}
			value *= rhs
		default:
			return value, nil
		}
	}
	return value, nil
}

func (p *parser) parseFactor() (int, error) {
	t, err := p.nextToken()
	if err != nil {
		return 0, err
	}
	switch t.kind {
	case tokenNum:
		return t.num, nil
	case tokenLParen:
		value, err := p.parseExpr()
		if err != nil {
			return 0, err
		}
		closing, err := p.nextToken()
		if err != nil || closing.kind != tokenRParen {
			return 0, fmt.Errorf("expected closing parenthesis")
		}
		return value, nil
	default:
		return 0, fmt.Errorf("unexpected token, expected a number or '('")
	}
}

func evalLine(tokens []token, additionBindsTighter bool) (int, error) {
	p := &parser{tokens: tokens, additionBindsTighter: additionBindsTighter}
	value, err := p.parseExpr()
	if err != nil {
		return 0, err
	}
	if p.pos != len(tokens) {
		return 0, fmt.Errorf("trailing tokens after a complete expression")
	}
	return value, nil
}

func evalExpression(line string, additionBindsTighter bool) (int, error) {
	tokens, err := tokenize(line)
	if err != nil {
		return 0, err
	}
	return evalLine(tokens, additionBindsTighter)
}
