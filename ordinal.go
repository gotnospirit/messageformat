package messageformat

import (
	"bytes"
	"fmt"
	"strconv"
)

// formatOrdinal is the format function associated with the "selectordinal" type.
//
// It will returns an error if :
// - the associated value can't be convert to string or to an int (i.e. bool, ...)
// - the pluralFunc is not defined (MessageFormat.getNamedKey)
//
// It will falls back to the "other" choice if :
// - its key can't be found in the given map
// - the computed named key (MessageFormat.getNamedKey) is not a key of the given map
func formatOrdinal(expr Expression, ptr_output *bytes.Buffer, data *map[string]interface{}, ptr_mf *MessageFormat, _ string) error {
	o := expr.(*selectExpr)

	value, err := toString(*data, o.key)
	if nil != err {
		return err
	}

	var choice *node

	if v, ok := (*data)[o.key]; ok {
		float_value, err := toFloat(v)
		if nil != err {
			return err
		}

		key, err := ptr_mf.getNamedKey(float_value, true)
		if nil != err {
			return err
		}
		choice = o.choices[key]
	}

	if nil == choice {
		choice = o.choices["other"]
	}
	return choice.format(ptr_output, data, ptr_mf, value)
}

// toFloat converts an interface{} value to a float64.
//
// It will returns an error if the value's type is not <string/int/float64>.
func toFloat(v interface{}) (float64, error) {
	switch v.(type) {
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

	default:
		return 0, fmt.Errorf("toFloat: Unsupported type: %T", v)
	}
}
