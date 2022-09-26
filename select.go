package messageformat

import (
	"bytes"
	"errors"
)

type selectExpr struct {
	Key     string                `json:"key"`
	Choices map[string]*ParseTree `json:"choices"`
}

func (p *parser) parseSelect(varname string, char rune, start, end int, ptr_input *[]rune) (Expression, int, error) {
	result := new(selectExpr)
	result.Key = varname
	result.Choices = make(map[string]*ParseTree)

	if PartChar != char {
		return nil, start, errors.New("MalformedOption")
	}

	hasOtherChoice := false

	pos := start + 1

	for pos < end {
		key, char, i, err := readKey(pos, end, ptr_input)

		if nil != err {
			return nil, i, err
		} else if char == ColonChar {
			return nil, i, errors.New("UnexpectedExtension")
		}

		if key == "other" {
			hasOtherChoice = true
		}

		choice, char, i, err := p.readChoice(char, i, end, ptr_input)
		if nil != err {
			return nil, i, err
		}

		result.Choices[key] = choice
		pos = i

		if CloseChar == char {
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
func (f *formatter) formatSelect(expr Expression, ptr_output *bytes.Buffer, data map[string]any) error {
	o := expr.(*selectExpr)

	value, err := toString(data, o.Key)
	if nil != err {
		return err
	}

	choice, ok := o.Choices[value]
	if !ok {
		choice = o.Choices["other"]
	}

	return f.format(choice, ptr_output, data, value)
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

func (p *parser) readChoice(char rune, pos, end int, ptr_input *[]rune) (*ParseTree, rune, int, error) {
	if OpenChar != char {
		return nil, char, pos, errors.New("MissingChoiceContent")
	}

	choice := new(ParseTree)

	pos, _, err := p.parse(pos+1, end, ptr_input, choice)
	if nil != err {
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
