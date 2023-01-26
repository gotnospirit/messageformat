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

type ExpressionType = string

const (
	DateExpression    ExpressionType = "date"
	PluralExpression  ExpressionType = "plural"
	SelectExpression  ExpressionType = "select"
	OrdinalExpression ExpressionType = "selectordinal"
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
	cursor, end := 0, len(runes)

	root := ParseTree{}
	for cursor < end {
		i, level, err := x.parse(cursor, end, &runes, &root)
		if nil != err {
			return nil, parseError{err.Error(), i}
		} else if level != 0 {
			return nil, parseError{"UnbalancedBraces", i}
		}

		cursor = i
	}

	return &root, nil
}

func (p *parser) parseExpression(start, end int, ptr_input *[]rune) (string, Expression, int, error) {
	var cursor int
	var expr Expression

	// This is the start of an expression like { var, expr, params}
	varName, nextChar, cursor, err := readVar(start, end, ptr_input)
	if err != nil {
		return "", nil, cursor, err
	}

	if varName == "" {
		return "", nil, cursor, fmt.Errorf("MissingVarName")
	}

	// if its just { var } return a VarExpr
	if nextChar == CloseChar {
		return "var", VarExpr{varName}, cursor, nil
	}

	exprType, nextChar, cursor, err := readVar(cursor+1, end, ptr_input)
	if nil != err {
		return "", nil, cursor, err
	}

	switch exprType {
	case DateExpression:
		expr, cursor, err = p.parseDate(varName, nextChar, cursor, end, ptr_input)
	case PluralExpression:
		expr, cursor, err = p.parsePlural(varName, nextChar, cursor, end, ptr_input)
	case SelectExpression:
		fallthrough
	case OrdinalExpression:
		expr, cursor, err = p.parseSelect(varName, nextChar, cursor, end, ptr_input)
	default:
		return "", nil, cursor, fmt.Errorf("UnknownType: `%s`", exprType)
	}
	if err != nil {
		return "", nil, cursor, err
	}

	if cursor >= end || CloseChar != (*ptr_input)[cursor] {
		return "", nil, cursor, fmt.Errorf("UnbalancedBraces")
	}

	return exprType, expr, cursor, nil
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
				if nil != err {
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
