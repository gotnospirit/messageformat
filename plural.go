package messageformat

import (
	"bytes"
	"fmt"
	"strconv"
)

type pluralExpr struct {
	selectExpr
	offset int
}

func parsePlural(varname string, ptr_compiler *Parser, char rune, start, end int, ptr_input *[]rune) (Expression, int, error) {
	if char != PartChar {
		return nil, start, fmt.Errorf("MalformedOption")
	}

	hasOtherChoice := false

	result := new(pluralExpr)
	result.key = varname
	result.choices = make(map[string]*node)

	pos := start + 1

	for pos < end {
		key, char, i, err := readKey(pos, end, ptr_input)

		if err != nil {
			return nil, i, err
		}

		if char == ':' {
			if key != "offset" {
				return nil, i, fmt.Errorf("UnsupportedExtension: `%s`", key)
			}

			offset, c, j, err := readOffset(i+1, end, ptr_input)
			if err != nil {
				return nil, j, err
			}

			result.offset = offset

			if isWhitespace(c) {
				j++
			}

			k, c, j, err := readKey(j, end, ptr_input)

			if err != nil {
				return nil, j, err
			} else if k == "" {
				return nil, j, fmt.Errorf("MissingChoiceName")
			}

			key, char, i = k, c, j
		}

		if key == "other" {
			hasOtherChoice = true
		}

		choice, c, i, err := readChoice(ptr_compiler, char, i, end, ptr_input)
		if err != nil {
			return nil, i, err
		}

		result.choices[key] = choice
		pos, char = i, c

		if char == CloseChar {
			break
		}
	}

	if !hasOtherChoice {
		return nil, pos, fmt.Errorf("MissingMandatoryChoice")
	}
	return result, pos, nil
}

// formatPlural is the format function associated with the "plural" type.
//
// It will returns an error if :
// - the associated value can't be convert to string or to an int (i.e. bool, ...)
// - the pluralFunc is not defined (MessageFormat.getNamedKey)
//
// It will falls back to the "other" choice if :
// - its key can't be found in the given map
// - the computed named key (MessageFormat.getNamedKey) is not a key of the given map
func formatPlural(expr Expression, ptr_output *bytes.Buffer, data *map[string]interface{}, ptr_mf *MessageFormat, _ string) error {
	o := expr.(*pluralExpr)
	key := o.key
	offset := o.offset

	value, err := toString(*data, key)
	if err != nil {
		return err
	}

	var choice *node

	if v, ok := (*data)[key]; ok {
		switch t := v.(type) {
		default:
			return fmt.Errorf("Plural: Unsupported type for named key: %T", v)

		case int:
			key = fmt.Sprintf("=%d", t)

		case float64:
			key = "=" + strconv.FormatFloat(t, 'f', -1, 64)

		case string:
			key = "=" + t
		}

		if choice = o.choices[key]; choice == nil {
			switch t := v.(type) {
			case int:
				if offset != 0 {
					offset_value := t - offset
					value = fmt.Sprintf("%d", offset_value)
					key, err = ptr_mf.getNamedKey(offset_value, false)
				} else {
					key, err = ptr_mf.getNamedKey(t, false)
				}

			case float64:
				if offset != 0 {
					offset_value := t - float64(offset)
					value = strconv.FormatFloat(offset_value, 'f', -1, 64)
					key, err = ptr_mf.getNamedKey(offset_value, false)
				} else {
					key, err = ptr_mf.getNamedKey(t, false)
				}

			case string:
				if offset != 0 {
					offset_value, fError := strconv.ParseFloat(value, 64)
					if fError != nil {
						return fError
					}
					offset_value -= float64(offset)
					value = strconv.FormatFloat(offset_value, 'f', -1, 64)
					key, err = ptr_mf.getNamedKey(offset_value, false)
				} else {
					key, err = ptr_mf.getNamedKey(value, false)
				}
			}

			if err != nil {
				return err
			}
			choice = o.choices[key]
		}
	}

	if choice == nil {
		choice = o.choices["other"]
	}
	return choice.format(ptr_output, data, ptr_mf, value)
}

func readOffset(start, end int, ptr_input *[]rune) (int, rune, int, error) {
	var buf bytes.Buffer

	char, pos := whitespace(start, end, ptr_input)
	input := *ptr_input

	for pos < end {
		switch char {
		default:
			buf.WriteRune(char)
			pos++

			if pos < end {
				char = input[pos]
			}

		case ' ', '\r', '\n', '\t', OpenChar, CloseChar:
			if buf.Len() != 0 {
				result, err := strconv.Atoi(buf.String())
				if err != nil {
					return 0, char, pos, fmt.Errorf("BadCast")
				} else if result < 0 {
					return 0, char, pos, fmt.Errorf("InvalidOffsetValue")
				}
				return result, char, pos, nil
			}
			return 0, char, pos, fmt.Errorf("MissingOffsetValue")
		}
	}
	return 0, char, pos, fmt.Errorf("UnbalancedBraces")
}
