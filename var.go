package messageformat

import (
	"bytes"
	"errors"
)

func formatVar(expr Expression, ptr_output *bytes.Buffer, data *map[string]interface{}, _ *MessageFormat, _ string) error {
	value, err := toString(*data, expr.(string))
	if err != nil {
		return err
	}
	ptr_output.WriteString(value)
	return nil
}

func readVar(start, end int, ptr_input *[]rune) (string, rune, int, error) {
	char, pos := whitespace(start, end, ptr_input)
	fc_pos, lc_pos := pos, pos
	input := *ptr_input

	for pos < end {
		switch char {
		default:
			// [_0-9a-zA-Z]+
			if char != '_' && (char < '0' || char > '9') && (char < 'A' || char > 'Z') && (char < 'a' || char > 'z') {
				return "", char, pos, errors.New("InvalidFormat")
			} else if pos != lc_pos { // non continu (inner whitespace)
				return "", char, pos, errors.New("InvalidFormat")
			}

			lc_pos = pos + 1

			pos++

			if pos < end {
				char = input[pos]
			}

		case ' ', '\r', '\n', '\t':
			char, pos = whitespace(pos+1, end, ptr_input)

		case PartChar, CloseChar:
			return string(input[fc_pos:lc_pos]), char, pos, nil

		case OpenChar:
			return "", char, pos, errors.New("InvalidExpr")
		}
	}
	return "", char, pos, errors.New("UnbalancedBraces")
}
