package messageformat

import (
	"bytes"
	"fmt"

	"github.com/gotnospirit/makeplural/plural"
	"golang.org/x/text/language"
)

var (
	DefaultCulture = "en"
)

// pluralFunc describes a function used to produce a named key when processing a plural or selectordinal expression.
type pluralFunc func(interface{}, bool) string

type Formatter interface {
	Format(*ParseTree) (string, error)
	FormatMap(*ParseTree, map[string]any) (string, error)
}

type FormatterOption func(f *formatter)

func WithLocale(locale language.Tag) FormatterOption {
	return func(f *formatter) {
		switch locale {
		case language.AmericanEnglish:
			f.date = createAmericanDateFormatter()
			f.locale = locale
		case language.German:
			f.date = createGermanDateFormatter()
			f.locale = locale
		}
	}
}

// TODO(tylermorton): replace this logic to use WithLocale
// WithCulture applies the plurality culture to the formatter
func WithCulture(culture string) FormatterOption {
	return func(f *formatter) {
		f.SetCulture(culture)
	}
}

// NewFormatter creates a new formatter with the given options
func NewFormatter(opts ...FormatterOption) (Formatter, error) {
	f := formatter{}

	for _, opt := range opts {
		opt(&f)
	}

	if f.plural == nil {
		err := f.SetCulture(DefaultCulture)
		if err != nil {
			return nil, err
		}
	}

	return &f, nil
}

type formatter struct {
	date DateFormatter

	// the locale this formatter is configured to
	// output as. this is used for dates and times
	locale language.Tag
	// TODO: replace this with locale-specific logic
	plural pluralFunc
}

func (x *formatter) SetCulture(name string) error {
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
	err := n.ForEach(func(n *Node) error {
		switch n.Type {
		case "date":
			return f.formatDate(n.Expr, buf, data)
		case "literal":
			return f.formatLiteral(n.Expr, buf, value)
		case "plural":
			return f.formatPlural(n.Expr, buf, data)
		case "select":
			return f.formatSelect(n.Expr, buf, data)
		case "selectordinal":
			return f.formatOrdinal(n.Expr, buf, data)
		case "time":
			return fmt.Errorf("formatter not implemented for time")
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
