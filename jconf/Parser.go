// Copyright 2021 The SSDB-cluster Authors
package jconf

import (
	"fmt"
	"errors"
)

type Parser struct {
	tn *Tokenizer
}

// keep fields in the order as they appear
func (parser *Parser)parse() (ret *Object, err error) {
	token := parser.tn.NextToken()
	if token == '{' {
		ret = NewObject()
		err = parser.parse_object(ret)
		if err != nil {
			return nil, err
		}
	} else if token == '[' {
		ret = NewArray()
		err = parser.parse_array(ret)
		if err != nil {
			return nil, err
		}
	} else if token == TokenQuote {
		str, err := parser.tn.parse_string()
		if err != nil {
			return nil, err
		}
		ret = NewString(str)
	} else {
		return nil, fmt.Errorf("unexpected token: '%c'", token)
	}
	return ret, nil
}

func (parser *Parser)parse_array(arr *Object) error {
	prev := TokenNone
	next := TokenNone

	for {
		next, _ = parser.tn.PeekToken()
		if next == ']' {
			parser.tn.NextToken()
			if prev == TokenComma && arr.Count() == 0 {
				return errors.New("',' in empty array")
			}
			return nil
		} else {
			sub, err := parser.parse()
			if err != nil {
				return err
			}
			arr.Push(sub)

			prev, _ = parser.tn.PeekToken()
			if prev == TokenComma {
				parser.tn.NextToken()
			}
		}
	}
}

func (parser *Parser)parse_object(obj *Object) error {
	prev := TokenNone
	next := TokenNone

	for {
		next = parser.tn.NextToken()
		if next == '}' {
			if prev == TokenComma && obj.Count() == 0 {
				return errors.New("',' in empty array")
			}
			return nil
		} else if next == TokenQuote {
			key, err := parser.tn.parse_string()
			if err != nil {
				return err
			}

			next = parser.tn.NextToken()
			if next != TokenColon {
				fmt.Println(string(parser.tn.buf))
				return errors.New("missing ':'")
			}

			sub, err := parser.parse()
			if err != nil {
				return err
			}
			obj.Set(key, sub)

			prev, _ = parser.tn.PeekToken()
			if prev == TokenComma {
				parser.tn.NextToken()
			}
		} else {
			return fmt.Errorf("unexpected token parsing object: '%c'", next)
		}
	}
}
