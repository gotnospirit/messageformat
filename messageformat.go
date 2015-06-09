// Package messageformat implements ICU messages formatting for Go.
// see http://userguide.icu-project.org/formatparse/messages
// inspired by https://github.com/SlexAxton/messageformat.js
package messageformat

import (
	"bytes"
	"fmt"
	"github.com/gotnospirit/makeplural/plural"
)

type MessageFormat struct {
	root       node
	formatters map[string]formatFunc
	plural     pluralFunc
}

func (x *MessageFormat) SetCulture(name string) error {
	fn, err := plural.GetFunc(name)
	if nil != err {
		return err
	}
	x.plural = fn
	return nil
}

func (x *MessageFormat) SetPluralFunction(fn pluralFunc) error {
	if nil == fn {
		return fmt.Errorf("PluralFunctionRequired")
	}
	x.plural = fn
	return nil
}

func (x *MessageFormat) Format() (string, error) {
	return x.FormatMap(nil)
}

func (x *MessageFormat) FormatMap(data map[string]interface{}) (string, error) {
	var buf bytes.Buffer

	err := x.root.format(&buf, &data, x, "")
	if nil != err {
		return "", err
	}
	return buf.String(), nil
}

func (x *MessageFormat) getNamedKey(value interface{}, ordinal bool) (string, error) {
	if nil == x.plural {
		return "", fmt.Errorf("UndefinedPluralFunc")
	}
	return x.plural(value, ordinal), nil
}

func (x *MessageFormat) getFormatter(key string) (formatFunc, error) {
	fn, ok := x.formatters[key]
	if !ok {
		return nil, fmt.Errorf("UnknownType: `%s`", key)
	} else if nil == fn {
		return nil, fmt.Errorf("UndefinedFormatFunc: `%s`", key)
	}
	return fn, nil
}
