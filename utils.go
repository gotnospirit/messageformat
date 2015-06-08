package messageformat

import (
	"fmt"
	"strconv"
)

// isWhitespace returns true if the rune is a whitespace.
func isWhitespace(char rune) bool {
	return ' ' == char || '\r' == char || '\n' == char || '\t' == char
}

// whitespace traverses the input until a non-whitespace char is encountered.
func whitespace(start, end int, ptr_input *[]rune) (rune, int) {
	input := *ptr_input

	pos := start
	for pos < end {
		char := input[pos]

		switch char {
		default:
			return char, pos

		case ' ', '\r', '\n', '\t':
		}

		pos++
	}
	return 0, pos
}

// toString retrieves a value from the given map and tries to return a string representation.
//
// It will returns an error if the value's type is not <nil/string/int/float64>.
func toString(data map[string]interface{}, key string) (string, error) {
	if v, ok := data[key]; ok {
		switch v.(type) {
		default:
			return "", fmt.Errorf("toString: Unsupported type: %T", v)

		case nil:
			return "", nil

		case string:
			return v.(string), nil

		case int:
			return fmt.Sprintf("%d", v.(int)), nil

		case float64:
			return fmt.Sprintf("%.2f", v.(float64)), nil
		}
	}
	return "", nil
}

// toFloat converts an interface{} value to a float64.
//
// It will returns an error if the value's type is not <string/int/float64>.
func toFloat(v interface{}) (float64, error) {
	switch v.(type) {
	default:
		return 0, fmt.Errorf("toFloat: Unsupported type: %T", v)

	case int:
		return float64(v.(int)), nil

	case float64:
		return v.(float64), nil

	case string:
		value, err := strconv.ParseFloat(v.(string), 64)
		if nil != err {
			return 0, err
		}
		return value, nil
	}
}
