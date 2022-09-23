package messageformat

import (
	"bytes"
	"fmt"

	"github.com/gotnospirit/makeplural/plural"
)

type Formatter interface {
	Format(*ParseTree) (string, error)
	FormatMap(*ParseTree, map[string]any) (string, error)
}

func NewFormatter() Formatter {
	return NewFormatterWithCulture("en")
}

func NewFormatterWithCulture(culture string) Formatter {
	f := &formatter{}
	err := f.SetCulture("en")

	// TODO: return an error here?
	// refactor to not throw
	if err != nil {
		panic(err)
	}

	return f
}

type formatter struct {
	plural pluralFunc
}

func (x *formatter) SetCulture(name string) error {
	// TODO: refactor to not throw
	fn, err := plural.GetFunc(name)
	if nil != err {
		return err
	}

	x.plural = fn
	return nil
}

func (x *formatter) SetPluralFunction(fn pluralFunc) error {
	if nil == fn {
		return fmt.Errorf("PluralFunctionRequired")
	}
	x.plural = fn

	return nil
}

func (f *formatter) Format(n *ParseTree) (string, error) {
	return f.FormatMap(n, nil)
}

func (f *formatter) FormatMap(n *ParseTree, data map[string]any) (string, error) {
	var buf bytes.Buffer

	err := f.format(n, &buf, data, "")
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}

func (f *formatter) format(n *ParseTree, buf *bytes.Buffer, data map[string]any, value string) error {
	err := n.forEach(func(n *Node) error {
		switch n.Type {
		case "literal":
			return f.formatLiteral(n.Expr, buf, value)
		case "plural":
			return f.formatPlural(n.Expr, buf, data)
		case "select":
			return f.formatSelect(n.Expr, buf, data)
		case "selectordinal":
			return f.formatOrdinal(n.Expr, buf, data)
		case "var":
			return f.formatVar(n.Expr, buf, data)
		default:
			return fmt.Errorf("formatter not implemented for expression of type %s", n.Type)
		}
	})
	if nil != err {
		return err
	}

	return nil
}

func (f *formatter) getNamedKey(value interface{}, ordinal bool) (string, error) {
	if nil == f.plural {
		return "", fmt.Errorf("UndefinedPluralFunc")
	}

	return f.plural(value, ordinal), nil
}
