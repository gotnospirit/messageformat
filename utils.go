package messageformat

import (
	"fmt"
	"strconv"
	"time"
)

// isWhitespace returns true if the rune is a whitespace.
func isWhitespace(char rune) bool {
	return char == ' ' || char == '\r' || char == '\n' || char == '\t'
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
// It will returns an error if the value's type is not <nil/string/bool/numeric/time.Duration/fmt.Stringer>.
func toString(data map[string]interface{}, key string) (string, error) {
	if v, ok := data[key]; ok {
		switch val := v.(type) {
		default:
			return "", fmt.Errorf("toString: Unsupported type: %T", v)

		case nil:
			return "", nil

		case bool:
			return strconv.FormatBool(val), nil

		case string:
			return val, nil

		case int:
			return fmt.Sprintf("%d", val), nil

		case int8:
			return strconv.FormatInt(int64(val), 10), nil

		case int16:
			return strconv.FormatInt(int64(val), 10), nil

		case int32:
			return strconv.FormatInt(int64(val), 10), nil

		case int64:
			return strconv.FormatInt(val, 10), nil

		case uint:
			return strconv.FormatUint(uint64(val), 10), nil

		case uint8:
			return strconv.FormatUint(uint64(val), 10), nil

		case uint16:
			return strconv.FormatUint(uint64(val), 10), nil

		case uint32:
			return strconv.FormatUint(uint64(val), 10), nil

		case uint64:
			return strconv.FormatUint(val, 10), nil

		case float32:
			return strconv.FormatFloat(float64(val), 'f', -1, 32), nil

		case float64:
			return strconv.FormatFloat(val, 'f', -1, 64), nil

		case complex64:
			return fmt.Sprintf("%g", val), nil

		case complex128:
			return fmt.Sprintf("%g", val), nil

		case uintptr:
			return fmt.Sprintf("%08x", val), nil

		case time.Duration:
			return val.String(), nil

		case fmt.Stringer:
			return val.String(), nil
		}
	}
	return "", nil
}
