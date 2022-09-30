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
	if err != nil {
		return err
	}

	var choice *node

	if v, ok := (*data)[o.key]; ok {
		switch t := v.(type) {
		default:
			return fmt.Errorf("Ordinal: Unsupported type for named key: %T", v)

		case int, float64:

		case string:
			_, err := strconv.ParseFloat(t, 64)
			if err != nil {
				return err
			}
		}

		key, err := ptr_mf.getNamedKey(v, true)
		if err != nil {
			return err
		}
		choice = o.choices[key]
	}

	if choice == nil {
		choice = o.choices["other"]
	}
	return choice.format(ptr_output, data, ptr_mf, value)
}
