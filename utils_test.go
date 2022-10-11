package messageformat

import (
	"fmt"
	"testing"
	"time"
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

func toStringError(t *testing.T, data map[string]interface{}, key string, expected string) {
	_, err := toString(data, key)

	if err == nil {
		t.Errorf("Expecting an error but got none")
	} else if expected != err.Error() {
		t.Errorf("Expecting error `%s` but got `%s`", expected, err)
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
		"N": nil,
		"B": true,
	}

	// should returns an empty string when the key does not exists
	toStringResult(t, data, "NAME", "")

	// should returns an empty string when the value is nil
	toStringResult(t, data, "N", "")

	// should otherwise returns a string representation (string, int, float)
	toStringResult(t, data, "S", "I am a string")
	toStringResult(t, data, "I", "42")
	toStringResult(t, data, "F", "0.305")
	toStringResult(t, data, "B", "true")

}

func TestToStringNumericTypes(t *testing.T) {
	data := map[string]interface{}{
		"byteMax":    byte(255),
		"uint":       uint(123456),
		"uint8Max":   uint8(255),
		"uint16Max":  uint16(65535),
		"uint32Max":  uint32(4294967295),
		"uint64Max":  uint64(18446744073709551615),
		"int":        int(-123456),
		"int8Min":    int8(-128),
		"int8Max":    int8(127),
		"int16Min":   int16(-32768),
		"int16Max":   int16(32767),
		"int32Min":   int32(-2147483648),
		"int32Max":   int32(2147483647),
		"int64Min":   int64(-9223372036854775808),
		"int64Max":   int64(9223372036854775807),
		"float32":    float32(3.14),
		"float64":    float64(3.14e-10),
		"complex":    complex(1.23, 9.87),
		"complex64":  complex64(1.23 + 9.87i),
		"complex128": complex128(1.23 + 9.87i),
		"rune":       rune('a'),
		"uintptr":    uintptr(0x75bcd15),
	}

	toStringResult(t, data, "byteMax", "255")
	toStringResult(t, data, "uint", "123456")
	toStringResult(t, data, "uint8Max", "255")
	toStringResult(t, data, "uint16Max", "65535")
	toStringResult(t, data, "uint32Max", "4294967295")
	toStringResult(t, data, "uint64Max", "18446744073709551615")
	toStringResult(t, data, "int", "-123456")
	toStringResult(t, data, "int8Min", "-128")
	toStringResult(t, data, "int8Max", "127")
	toStringResult(t, data, "int16Min", "-32768")
	toStringResult(t, data, "int16Max", "32767")
	toStringResult(t, data, "int32Min", "-2147483648")
	toStringResult(t, data, "int32Max", "2147483647")
	toStringResult(t, data, "int64Min", "-9223372036854775808")
	toStringResult(t, data, "int64Max", "9223372036854775807")
	toStringResult(t, data, "float32", "3.14")
	toStringResult(t, data, "float64", "0.000000000314")
	toStringResult(t, data, "complex", "(1.23+9.87i)")
	toStringResult(t, data, "complex64", "(1.23+9.87i)")
	toStringResult(t, data, "complex128", "(1.23+9.87i)")
	toStringResult(t, data, "rune", "97")
	toStringResult(t, data, "uintptr", "075bcd15")
}

func TestToStringBool(t *testing.T) {
	data := map[string]interface{}{
		"boolTrue":  bool(true),
		"boolFalse": bool(false),
	}

	toStringResult(t, data, "boolTrue", "true")
	toStringResult(t, data, "boolFalse", "false")
}

func TestToStringTimeDuration(t *testing.T) {
	du := time.Date(1970, 0, 1, 0, 0, 0, 0, time.UTC).Sub(time.Date(1960, 0, 1, 0, 0, 0, 0, time.UTC))
	data := map[string]interface{}{
		"duration": du,
	}
	toStringResult(t, data, "duration", "87672h0m0s")
}

type testStruct struct {
	Value int
}

// Implement the fmt.Stringer interface
func (t testStruct) String() string {
	return fmt.Sprintf("%d", t.Value)
}

func TestToStringStringer(t *testing.T) {
	data := map[string]interface{}{
		"struct":    testStruct{Value: 1},
		"structPtr": &testStruct{Value: 2},
	}

	toStringResult(t, data, "struct", "1")
	toStringResult(t, data, "structPtr", "2")
}

func TestToStringException(t *testing.T) {
	data := map[string]any{
		"emptyStruct": struct{}{},
	}

	toStringError(t, data, "emptyStruct", "toString: Unsupported type: struct {}")
}
