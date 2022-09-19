package internal

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"unicode"
)

const (
	PATH TokenType = iota
	NUMBER
	BOOL
	STRING
	IF_NOT_FOUND_OP
	MINUS
	PLUS_OP
	SUM_WORD
	LEFT_PAREN
	RIGHT_PAREN
	TIMES_OP
	PRODUCT_WORD
	DIVIDE_OP
	LENGTH_WORD
	NOT_OP
	LESS_THAN_OP
	LESS_THAN_OR_EQUAL_OP
	GREATER_THAN_OP
	GREATER_THAN_OR_EQUAL_OP
	EQUAL_OP
	AND_OP
	AND_WORD
	OR_OP
	OR_WORD
	COMMA
)

type (
	TokenType int

	Token struct {
		Type  TokenType
		Value interface{}
	}
)

var (
	idStart = &unicode.RangeTable{R16: []unicode.Range16{
		{Lo: 'A', Hi: 'Z', Stride: 1},
		{Lo: '_', Hi: '_', Stride: 1},
		{Lo: 'a', Hi: 'z', Stride: 1},
	}}

	idRune = &unicode.RangeTable{R16: []unicode.Range16{
		{Lo: '0', Hi: '9', Stride: 1},
		{Lo: 'A', Hi: 'Z', Stride: 1},
		{Lo: '_', Hi: '_', Stride: 1},
		{Lo: 'a', Hi: 'z', Stride: 1},
	}}

	numRune = &unicode.RangeTable{R16: []unicode.Range16{
		{Lo: '0', Hi: '9', Stride: 1},
	}}
)

func Lex(expr string) ([]Token, error) {
	var tokens []Token
	iter := &stringIterator{runes: []rune(expr)}
	var err error
	for !iter.done() && err == nil {
		r, _ := iter.next()
		switch {
		case r == '$':
			t, e := readPath(iter)
			tokens, err = append(tokens, t), e
		case r == '-':
			if peek, ok := iter.peek(); !ok || !unicode.Is(numRune, peek) {
				tokens = append(tokens, Token{Type: MINUS})
			} else {
				_, _ = iter.next()
				numStr := readNum(iter, r, peek)
				num, e := convertFloat(numStr)
				tokens, err = append(tokens, Token{Type: NUMBER, Value: num}), e
			}
		case unicode.Is(numRune, r):
			numStr := readNum(iter, r)
			num, e := convertFloat(numStr)
			tokens, err = append(tokens, Token{Type: NUMBER, Value: num}), e
		case unicode.Is(idStart, r):
			id := readID(iter, r)
			t, e := idToToken(id)
			tokens, err = append(tokens, t), e
		case r == '\'':
			s, e := readString(iter)
			tokens, err = append(tokens, Token{Type: STRING, Value: s}), e
		case r == '?':
			tokens = append(tokens, Token{Type: IF_NOT_FOUND_OP})
		case r == '+':
			tokens = append(tokens, Token{Type: PLUS_OP})
		case r == '(':
			tokens = append(tokens, Token{Type: LEFT_PAREN})
		case r == ')':
			tokens = append(tokens, Token{Type: RIGHT_PAREN})
		case r == '*':
			tokens = append(tokens, Token{Type: TIMES_OP})
		case r == '/':
			tokens = append(tokens, Token{Type: DIVIDE_OP})
		case r == '<':
			if peek, ok := iter.peek(); ok && peek == '=' {
				_, _ = iter.next()
				tokens = append(tokens, Token{Type: LESS_THAN_OR_EQUAL_OP})
			} else {
				tokens = append(tokens, Token{Type: LESS_THAN_OP})
			}
		case r == '>':
			if peek, ok := iter.peek(); ok && peek == '=' {
				_, _ = iter.next()
				tokens = append(tokens, Token{Type: GREATER_THAN_OR_EQUAL_OP})
			} else {
				tokens = append(tokens, Token{Type: GREATER_THAN_OP})
			}
		case r == '!':
			tokens = append(tokens, Token{Type: NOT_OP})
		case r == '=':
			t, e := readDoubleToken(iter, '=', EQUAL_OP)
			tokens, err = append(tokens, t), e
		case r == '&':
			t, e := readDoubleToken(iter, '&', AND_OP)
			tokens, err = append(tokens, t), e
		case r == '|':
			t, e := readDoubleToken(iter, '|', OR_OP)
			tokens, err = append(tokens, t), e
		case r == ',':
			tokens = append(tokens, Token{Type: COMMA})
		case unicode.IsSpace(r):
		default:
			err = fmt.Errorf("unexpected token %q", r)
		}
	}
	if err != nil {
		return nil, err
	}
	return tokens, nil
}

// readPath is called after a '$' rune is read, which starts a path.
func readPath(iter *stringIterator) (Token, error) {
	var path []interface{}
outer:
	for {
		peek, ok := iter.peek()
		if !ok {
			break
		}
		switch peek {
		case '.':
			_, _ = iter.next()
			next, ok := iter.next()
			if !ok {
				return Token{}, errors.New("unexpected end of input")
			}
			if !unicode.Is(idStart, next) {
				return Token{}, fmt.Errorf("key in path cannot start with %q", next)
			}
			key := readID(iter, next)
			path = append(path, key)
		case '[':
			_, _ = iter.next()
			index, err := readIndex(iter)
			if err != nil {
				return Token{}, err
			}
			path = append(path, index)
		default:
			break outer
		}
	}
	return Token{Type: PATH, Value: path}, nil
}

func readIndex(iter *stringIterator) (interface{}, error) {
	next, ok := iter.next()
	if !ok {
		return nil, errors.New("unexpected end of input")
	}

	var index interface{}
	switch {
	case next == '\'':
		str, err := readString(iter)
		if err != nil {
			return nil, err
		}
		index = str
	case unicode.Is(numRune, next):
		numStr := readNum(iter, next)
		num, err := convertInt(numStr)
		if err != nil {
			return nil, err
		}
		index = num
	default:
		return nil, fmt.Errorf("unexpected token %q", next)
	}

	next, ok = iter.next()
	if !ok {
		return nil, errors.New("unexpected end of input")
	}
	if next != ']' {
		return nil, fmt.Errorf("unexpected token %q", next)
	}

	return index, nil
}

func readNum(iter *stringIterator, r ...rune) string {
	var seenDecimal bool
	sb := strings.Builder{}
	for _, r := range r {
		sb.WriteRune(r)
	}
outer:
	for {
		peek, ok := iter.peek()
		if !ok {
			break
		}
		switch {
		case unicode.Is(numRune, peek):
			_, _ = iter.next()
			sb.WriteRune(peek)
		case peek == '.':
			if !seenDecimal {
				_, _ = iter.next()
				sb.WriteRune(peek)
				seenDecimal = true
			} else {
				break outer
			}
		default:
			break outer
		}
	}
	return sb.String()
}

func convertFloat(numStr string) (float64, error) {
	num, err := strconv.ParseFloat(numStr, 64)
	if err != nil {
		return 0, fmt.Errorf("error parsing number %s", numStr)
	}
	return num, nil
}

func convertInt(numStr string) (int, error) {
	num, err := strconv.Atoi(numStr)
	if err != nil {
		return 0, fmt.Errorf("error parsing number %s", numStr)
	}
	return num, nil
}

func readID(iter *stringIterator, r rune) string {
	sb := strings.Builder{}
	sb.WriteRune(r)
	for {
		peek, ok := iter.peek()
		if !ok || !unicode.Is(idRune, peek) {
			break
		}
		_, _ = iter.next()
		sb.WriteRune(peek)
	}
	return sb.String()
}

func idToToken(id string) (Token, error) {
	switch id {
	case "true":
		return Token{Type: BOOL, Value: true}, nil
	case "false":
		return Token{Type: BOOL, Value: false}, nil
	case "sum":
		return Token{Type: SUM_WORD}, nil
	case "product":
		return Token{Type: PRODUCT_WORD}, nil
	case "length":
		return Token{Type: LENGTH_WORD}, nil
	case "and":
		return Token{Type: AND_WORD}, nil
	case "or":
		return Token{Type: OR_WORD}, nil
	default:
		return Token{}, fmt.Errorf("unexpected identifier %q", id)
	}
}

// readDoubleToken is called after the lexer reads one rune and expects to see
// the same rune again in order to add an operator token
func readDoubleToken(iter *stringIterator, want rune, tokenType TokenType) (Token, error) {
	next, ok := iter.next()
	if !ok {
		return Token{}, errors.New("unexpected end of input")
	}
	if next != want {
		return Token{}, fmt.Errorf("unexpected token %q after %q", next, want)
	}
	return Token{Type: tokenType}, nil
}

// readString is called after the lexer reads a '\” rune, which starts a
// string. This functions reads runes from the iterator until it finds an
// unescaped '\” rune. It will read an escaped rune after a '\\' rune.
//
// It will error if the iterator ends before an unescaped '\” rune, or if it
// reads a rune other than '\” or '\\' after a '\\' rune.
func readString(iter *stringIterator) (string, error) {
	sb := strings.Builder{}
loop:
	for {
		r, ok := iter.next()
		if !ok {
			return "", errors.New("unexpected end of input")
		}
		switch r {
		case '\'':
			break loop
		case '\\':
			escaped, ok := iter.next()
			if !ok {
				return "", errors.New("unexpected end of input")
			}
			if escaped != '\'' && escaped != '\\' {
				return "", fmt.Errorf("unexpected escaped token %q in string", escaped)
			}
			sb.WriteRune(escaped)
		default:
			sb.WriteRune(r)
		}
	}
	return sb.String(), nil
}
