package messageformat

import (
	"fmt"
)

const (
	EscapeChar = '\\'
	OpenChar   = '{'
	CloseChar  = '}'
	PartChar   = ','
	PoundChar  = '#'
	ColonChar  = ':'
)

type Parser interface {
	Parse(string) (*ParseTree, error)
}

type parser struct{}

func NewParser() Parser {
	return &parser{}
}

func (x *parser) Parse(input string) (*ParseTree, error) {
	runes := []rune(input)
	pos, end := 0, len(runes)

	root := ParseTree{}
	for pos < end {
		i, level, err := x.parse(pos, end, &runes, &root)
		if err != nil {
			return nil, parseError{err.Error(), i}
		} else if level != 0 {
			return nil, parseError{"UnbalancedBraces", i}
		}

		pos = i
	}

	return &root, nil
}

func (x *parser) parseExpression(start, end int, ptr_input *[]rune) (string, Expression, int, error) {
	var pos int
	var expr Expression

	varname, char, pos, err := readVar(start, end, ptr_input)
	if err != nil {
		return "", nil, pos, err
	}

	if varname == "" {
		return "", nil, pos, fmt.Errorf("MissingVarName")
	}

	if char == CloseChar {
		return "var", VarExpr{
			Name: varname,
		}, pos, nil
	}

	ctype, char, pos, err := readVar(pos+1, end, ptr_input)
	if err != nil {
		return "", nil, pos, err
	}

	switch ctype {
	case "plural":
		expr, pos, err = x.parsePlural(varname, char, pos, end, ptr_input)
	case "select":
		fallthrough
	case "selectordinal":
		expr, pos, err = x.parseSelect(varname, char, pos, end, ptr_input)
	default:
		return "", nil, pos, fmt.Errorf("UnknownType: `%s`", ctype)
	}
	if err != nil {
		return "", nil, pos, err
	}

	if pos >= end || CloseChar != (*ptr_input)[pos] {
		return "", nil, pos, fmt.Errorf("UnbalancedBraces")
	}

	return ctype, expr, pos, nil
}

func (p *parser) parse(start, end int, ptr_input *[]rune, parent *ParseTree) (int, int, error) {
	pos := start
	level := 0
	escaped := false
	input := *ptr_input

loop:
	for pos < end {
		char := input[pos]

		switch char {
		default:
			pos++
			escaped = false

		case EscapeChar:
			pos++
			escaped = true

		case CloseChar:
			if !escaped {
				level--
				break loop
			}
			pos++
			escaped = false

		case OpenChar:
			if !escaped {
				level++

				if pos > start {
					parent.add("literal", p.parseLiteral(start, pos, ptr_input))
				}

				ctype, child, i, err := p.parseExpression(pos+1, end, ptr_input)
				if err != nil {
					return i, level, err
				}

				parent.add(ctype, child)

				level--

				pos = i
				start = pos + 1
			}

			pos++
			escaped = false
		}
	}

	if pos > start {
		parent.add("literal", p.parseLiteral(start, pos, ptr_input))
	}

	return pos, level, nil
}
