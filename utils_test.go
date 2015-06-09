package messageformat

import (
	"fmt"
	"testing"
)

func toStringResult(t *testing.T, data map[string]interface{}, key, expected string) {
	result, err := toString(data, key)

	if nil != err {
		t.Errorf("Expecting `%s` but got an error `%s`", expected, err.Error())
	} else if expected != result {
		t.Errorf("Expecting `%s` but got `%s`", expected, result)
	} else if testing.Verbose() {
		fmt.Printf("Successfully returns the expected value: `%s`\n", expected)
	}
}

func toStringError(t *testing.T, data map[string]interface{}, key string) {
	result, err := toString(data, "B")

	if nil == err {
		t.Errorf("Expecting an error but got `%s`", result)
	} else if testing.Verbose() {
		fmt.Printf("Successfully returns an error `%s`\n", err.Error())
	}
}

func whitespaceResult(t *testing.T, start, end int, ptr_input *[]rune, expected_char rune, expected_pos int) int {
	char, pos := whitespace(start, end, ptr_input)
	if expected_pos != pos {
		t.Errorf("Expecting first non-whitespace found at %d but got %d", expected_pos, pos)
	} else if expected_char != char {
		t.Errorf("Expecting first non-whitespace was `%s` but got `%s`", string(expected_char), string(char))
	} else if testing.Verbose() {
		fmt.Printf("Successfully returns `%s`, %d\n", string(char), pos)
	}
	return pos
}

func TestIsWhitespace(t *testing.T) {
	for _, char := range []rune{' ', '\r', '\n', '\t'} {
		if true != isWhitespace(char) {
			t.Errorf("Do not returns true when receiving `%s`", string(char))
		}
	}

	if false != isWhitespace('a') {
		t.Errorf("Do not returns false when receiving `%s`", string('a'))
	}
}

func TestWhitespace(t *testing.T) {
	input := []rune(`  hello world`)
	start, end := 0, len(input)

	// should traverses the input, from "start" to "end"
	// until a non-whitespace char is encountered
	// and returns that char and its position
	pos := whitespaceResult(t, start, end, &input, 'h', 2)

	// should returns the same char and position because the char
	// at "start" position is not a whitespace
	whitespaceResult(t, pos, end, &input, 'h', 2)

	// should returns a \0 char when the end position is reached
	whitespaceResult(t, pos, pos, &input, 0, 2)
}

func TestToString(t *testing.T) {
	data := map[string]interface{}{
		"S": "I am a string",
		"I": 42,
		"F": 0.305,
		"B": true,
		"N": nil,
	}

	// should returns an empty string when the key does not exists
	toStringResult(t, data, "NAME", "")

	// should returns an empty string when the value is nil
	toStringResult(t, data, "N", "")

	// should returns an error when the value's type is not supported
	toStringError(t, data, "B")

	// should otherwise returns a string representation (string, int, float)
	toStringResult(t, data, "S", "I am a string")
	toStringResult(t, data, "I", "42")
	toStringResult(t, data, "F", "0.305")
}
