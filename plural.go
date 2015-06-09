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
	if PartChar != char {
		return nil, start, fmt.Errorf("MalformedOption")
	}

	hasOtherChoice := false

	result := new(pluralExpr)
	result.key = varname
	result.choices = make(map[string]*node)

	pos := start + 1

	for pos < end {
		key, char, i, err := readKey(char, pos, end, ptr_input)

		if nil != err {
			return nil, i, err
		}

		if ':' == char {
			if "offset" != key {
				return nil, i, fmt.Errorf("UnsupportedExtension: `%s`", key)
			}

			offset, c, j, err := readOffset(i+1, end, ptr_input)
			if nil != err {
				return nil, j, err
			}

			result.offset = offset

			if isWhitespace(c) {
				j++
			}

			k, c, j, err := readKey(c, j, end, ptr_input)

			if "" == k {
				return nil, j, fmt.Errorf("MissingChoiceName")
			}

			key, char, i = k, c, j
		}

		if "other" == key {
			hasOtherChoice = true
		}

		choice, c, i, err := readChoice(ptr_compiler, char, i, end, ptr_input)
		if nil != err {
			return nil, i, err
		}

		result.choices[key] = choice
		pos, char = i, c

		if CloseChar == char {
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
	if nil != err {
		return err
	}

	var choice *node

	if v, ok := (*data)[key]; ok {
		switch v.(type) {
		default:
			return fmt.Errorf("Plural: Unsupported type for named key: %T", v)

		case int:
			key = fmt.Sprintf("=%d", v.(int))

		case float64:
			key = "=" + strconv.FormatFloat(v.(float64), 'f', -1, 64)

		case string:
			key = "=" + v.(string)
		}

		if choice = o.choices[key]; nil == choice {
			switch v.(type) {
			case int:
				if 0 != offset {
					offset_value := v.(int) - offset
					value = fmt.Sprintf("%d", offset_value)
					key, err = ptr_mf.getNamedKey(offset_value, false)
				} else {
					key, err = ptr_mf.getNamedKey(v.(int), false)
				}

			case float64:
				if 0 != offset {
					offset_value := v.(float64) - float64(offset)
					value = strconv.FormatFloat(offset_value, 'f', -1, 64)
					key, err = ptr_mf.getNamedKey(offset_value, false)
				} else {
					key, err = ptr_mf.getNamedKey(v.(float64), false)
				}

			case string:
				if 0 != offset {
					offset_value, fError := strconv.ParseFloat(value, 64)
					if nil != fError {
						return fError
					}
					offset_value -= float64(offset)
					value = strconv.FormatFloat(offset_value, 'f', -1, 64)
					key, err = ptr_mf.getNamedKey(offset_value, false)
				} else {
					key, err = ptr_mf.getNamedKey(value, false)
				}
			}

			if nil != err {
				return err
			}
			choice = o.choices[key]
		}
	}

	if nil == choice {
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
			if 0 != buf.Len() {
				result, err := strconv.Atoi(buf.String())
				if nil != err {
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
