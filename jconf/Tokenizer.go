package jconf

import (
	"fmt"
	"errors"
	"bytes"
)

/*
tokens: '{', '}', '[', ']', ':', '"', ','
*/
type Tokenizer struct {
	buf []byte
}

type TokenType byte

const (
	TokenNone  TokenType = 'N'
	TokenError TokenType = 'E'
	TokenComma TokenType = ','
	TokenQuote TokenType = '"'
	TokenColon TokenType = ':'
)

func (t *Tokenizer)PeekToken() (ret TokenType, parsed int) {
	for s := 0; s < len(t.buf); s ++ {
		switch c := t.buf[s]; c {
		case ' ', '\t', '\r', '\n':
			continue
		case '{', '}', '[', ']', ':', '"', ',':
			return TokenType(c), s
		default:
			return TokenError, 0
		}
	}
	return TokenError, 0
}

func (t *Tokenizer)NextToken() TokenType {
	ret, n := t.PeekToken()
	if ret == TokenError {
		return ret
	}
	t.buf = t.buf[n + 1 : ]
	return ret
}

func (t *Tokenizer)parse_string() (string, error) {
	var ret bytes.Buffer
	var s int

Loop:
	for s = 0; s < len(t.buf); s ++ {
		switch c := t.buf[s]; c {
		case '\\':
			s += 1
			if s == len(t.buf) {
				return "", errors.New("error")
			}
			switch n := t.buf[s]; n {
			case '\\':
				ret.WriteByte('\\')
			case 't':
				ret.WriteByte('\t')
			case 'r':
				ret.WriteByte('\r')
			case 'n':
				ret.WriteByte('\n')
			default: // \f, \uxxxx, \x
				return "", fmt.Errorf("error character '%c' after \\", n)
			}
		case '"':
			break Loop
		default:
			ret.WriteByte(c)
		}
	}

	t.buf = t.buf[s + 1 : ]
	return ret.String(), nil
}
