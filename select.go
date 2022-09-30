package messageformat

import (
	"bytes"
	"errors"
)

type selectExpr struct {
	key     string
	choices map[string]*node
}

func parseSelect(varname string, ptr_compiler *Parser, char rune, start, end int, ptr_input *[]rune) (Expression, int, error) {
	result := new(selectExpr)
	result.key = varname
	result.choices = make(map[string]*node)

	if char != PartChar {
		return nil, start, errors.New("MalformedOption")
	}

	hasOtherChoice := false

	pos := start + 1

	for pos < end {
		key, char, i, err := readKey(pos, end, ptr_input)

		if err != nil {
			return nil, i, err
		} else if char == ':' {
			return nil, i, errors.New("UnexpectedExtension")
		}

		if key == "other" {
			hasOtherChoice = true
		}

		choice, char, i, err := readChoice(ptr_compiler, char, i, end, ptr_input)
		if err != nil {
			return nil, i, err
		}

		result.choices[key] = choice
		pos = i

		if char == CloseChar {
			break
		}
	}

	if !hasOtherChoice {
		return nil, pos, errors.New("MissingMandatoryChoice")
	}
	return result, pos, nil
}

// formatSelect is the format function associated with the "select" type.
//
// It will falls back to the "other" choice if :
// - its key can't be found in the given map
// - its string representation is not a key of the given map
//
// It will returns an error if :
// - the associated value can't be convert to string (i.e. struct {}, ...)
func formatSelect(expr Expression, ptr_output *bytes.Buffer, data *map[string]interface{}, ptr_mf *MessageFormat, _ string) error {
	o := expr.(*selectExpr)

	value, err := toString(*data, o.key)
	if err != nil {
		return err
	}

	choice, ok := o.choices[value]
	if !ok {
		choice = o.choices["other"]
	}
	return choice.format(ptr_output, data, ptr_mf, value)
}

func readKey(start, end int, ptr_input *[]rune) (string, rune, int, error) {
	char, pos := whitespace(start, end, ptr_input)
	fc_pos, lc_pos := pos, pos

	input := *ptr_input

	for pos < end {
		switch char {
		default:
			lc_pos = pos + 1

		case ' ', '\r', '\n', '\t':
			char, pos = whitespace(pos+1, end, ptr_input)
			return string(input[fc_pos:lc_pos]), char, pos, nil

		case ':', PartChar, CloseChar, OpenChar:
			if fc_pos != lc_pos {
				return string(input[fc_pos:lc_pos]), char, pos, nil
			}
			return "", char, pos, errors.New("MissingChoiceName")
		}

		pos++

		if pos < end {
			char = input[pos]
		}
	}
	return "", char, pos, errors.New("UnbalancedBraces")
}

func readChoice(ptr_compiler *Parser, char rune, pos, end int, ptr_input *[]rune) (*node, rune, int, error) {
	if char != OpenChar {
		return nil, char, pos, errors.New("MissingChoiceContent")
	}

	choice := new(node)

	pos, _, err := ptr_compiler.parse(pos+1, end, ptr_input, choice)
	if err != nil {
		return nil, char, pos, err
	}

	pos++
	if pos < end {
		char = (*ptr_input)[pos]
	}

	if isWhitespace(char) {
		char, pos = whitespace(pos+1, end, ptr_input)
	}
	return choice, char, pos, nil
}
