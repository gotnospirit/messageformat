package messageformat

import (
	"bytes"
	"fmt"
	"strconv"
)

type PluralExpr struct {
	Select SelectExpr `json:"select"`
	Offset int        `json:"offset"`
}

func (p *parser) parsePlural(varname string, char rune, start, end int, ptr_input *[]rune) (Expression, int, error) {
	if char != PartChar {
		return nil, start, fmt.Errorf("MalformedOption")
	}

	hasOtherChoice := false

	result := &PluralExpr{
		Select: SelectExpr{
			Key:     varname,
			Choices: make(map[string]*ParseTree),
		},
	}

	pos := start + 1

	for pos < end {
		key, char, i, err := readKey(pos, end, ptr_input)

		if nil != err {
			return nil, i, err
		}

		if char == ColonChar {
			if key != "offset" {
				return nil, i, fmt.Errorf("UnsupportedExtension: `%s`", key)
			}

			offset, c, j, err := readOffset(i+1, end, ptr_input)
			if nil != err {
				return nil, j, err
			}

			result.Offset = offset

			if isWhitespace(c) {
				j++
			}

			k, c, j, err := readKey(j, end, ptr_input)
			if err != nil || k == "" {
				return nil, j, fmt.Errorf("MissingChoiceName")
			}

			key, char, i = k, c, j
		}

		choice, c, i, err := p.readChoice(char, i, end, ptr_input)
		if nil != err {
			return nil, i, err
		}

		if key == "other" {
			hasOtherChoice = true
		}

		result.Select.Choices[key] = choice
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
func (f *formatter) formatPlural(expr Expression, ptr_output *bytes.Buffer, data map[string]any) error {
	o, ok := expr.(*PluralExpr)
	if !ok {
		return fmt.Errorf("expression is not a plural")
	}

	key := o.Select.Key
	offset := o.Offset

	value, err := toString(data, key)
	if nil != err {
		return err
	}

	var choice *ParseTree

	if v, ok := data[key]; ok {
		switch val := v.(type) {
		case int:
			key = fmt.Sprintf("=%d", val)

		case float64:
			key = "=" + strconv.FormatFloat(val, 'f', -1, 64)

		case string:
			key = "=" + val

		default:
			return fmt.Errorf("unsupported type for named plural key: %T", v)

		}

		if choice = o.Select.Choices[key]; nil == choice {
			switch val := v.(type) {
			case int:
				if offset != 0 {
					offset_value := val - offset
					value = fmt.Sprintf("%d", offset_value)
					key, err = f.getNamedKey(offset_value, false)
				} else {
					key, err = f.getNamedKey(v.(int), false)
				}

			case float64:
				if offset != 0 {
					offset_value := val - float64(offset)
					value = strconv.FormatFloat(offset_value, 'f', -1, 64)
					key, err = f.getNamedKey(offset_value, false)
				} else {
					key, err = f.getNamedKey(v.(float64), false)
				}

			case string:
				if offset != 0 {
					offset_value, fError := strconv.ParseFloat(value, 64)
					if nil != fError {
						return fError
					}
					offset_value -= float64(offset)
					value = strconv.FormatFloat(offset_value, 'f', -1, 64)
					key, err = f.getNamedKey(offset_value, false)
				} else {
					key, err = f.getNamedKey(value, false)
				}
			}

			if nil != err {
				return err
			}
			choice = o.Select.Choices[key]
		}
	}

	if choice == nil {
		choice = o.Select.Choices["other"]
	}

	return f.format(choice, ptr_output, data, value)
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
