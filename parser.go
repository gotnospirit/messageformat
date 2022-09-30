package messageformat

import (
	"bytes"
	"fmt"
	"github.com/gotnospirit/makeplural/plural"
)

const (
	EscapeChar = '\\'
	OpenChar   = '{'
	CloseChar  = '}'
	PartChar   = ','
	PoundChar  = '#'
)

type (
	// parseFunc describes a function used to parse a subset of the input string into an expression.
	parseFunc func(string, *Parser, rune, int, int, *[]rune) (Expression, int, error)
	// formatFunc describes a function used to format an expression into the output buffer.
	formatFunc func(Expression, *bytes.Buffer, *map[string]interface{}, *MessageFormat, string) error
	// pluralFunc describes a function used to produce a named key when processing a plural or selectordinal expression.
	pluralFunc func(interface{}, bool) string

	Parser struct {
		parsers    map[string]parseFunc
		formatters map[string]formatFunc
		plural     pluralFunc
	}
)

func (x *Parser) Parse(input string) (*MessageFormat, error) {
	runes := []rune(input)
	pos, end := 0, len(runes)

	root := node{}
	for pos < end {
		i, level, err := x.parse(pos, end, &runes, &root)
		if err != nil {
			return nil, parseError{err.Error(), i}
		} else if level != 0 {
			return nil, parseError{"UnbalancedBraces", i}
		}

		pos = i
	}
	return &MessageFormat{root, x.formatters, x.plural}, nil
}

func (x *Parser) Register(key string, p parseFunc, f formatFunc) error {
	if _, ok := x.parsers[key]; ok {
		return fmt.Errorf("ParserAlreadyRegistered")
	}
	x.parsers[key] = p
	x.formatters[key] = f
	return nil
}

func (x *Parser) parseExpression(start, end int, ptr_input *[]rune) (string, Expression, int, error) {
	varname, char, pos, err := readVar(start, end, ptr_input)
	if err != nil {
		return "", nil, pos, err
	} else if varname == "" {
		return "", nil, pos, fmt.Errorf("MissingVarName")
	} else if char == CloseChar {
		return "var", varname, pos, nil
	}

	ctype, char, pos, err := readVar(pos+1, end, ptr_input)
	if err != nil {
		return "", nil, pos, err
	}

	fn, ok := x.parsers[ctype]
	if !ok {
		return "", nil, pos, fmt.Errorf("UnknownType: `%s`", ctype)
	} else if fn == nil {
		return "", nil, pos, fmt.Errorf("UndefinedParseFunc: `%s`", ctype)
	}

	expr, pos, err := fn(varname, x, char, pos, end, ptr_input)
	if err != nil {
		return "", nil, pos, err
	}

	if pos >= end || (*ptr_input)[pos] != CloseChar {
		return "", nil, pos, fmt.Errorf("UnbalancedBraces")
	}
	return ctype, expr, pos, nil
}

func (x *Parser) parse(start, end int, ptr_input *[]rune, parent *node) (int, int, error) {
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
					parent.add("literal", parseLiteral(start, pos, ptr_input))
				}

				ctype, child, i, err := x.parseExpression(pos+1, end, ptr_input)
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
		parent.add("literal", parseLiteral(start, pos, ptr_input))
	}
	return pos, level, nil
}

func NewWithCulture(name string) (*Parser, error) {
	fn, err := plural.GetFunc(name)
	if err != nil {
		return nil, err
	}

	result := new(Parser)

	result.parsers = make(map[string]parseFunc)
	result.formatters = make(map[string]formatFunc)
	result.plural = fn

	result.Register("literal", nil, formatLiteral)
	result.Register("var", nil, formatVar)
	result.Register("select", parseSelect, formatSelect)
	result.Register("selectordinal", parseSelect, formatOrdinal)
	result.Register("plural", parsePlural, formatPlural)
	return result, nil
}

func New() (*Parser, error) {
	return NewWithCulture("en")
}
