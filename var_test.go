// @TODO(gotnospirit) add test on readVar

package messageformat

import (
	"testing"
)

func TestVar(t *testing.T) {
	doTest(t, Test{
		"Hello {NAME}",
		[]Expectation{
			{map[string]interface{}{"NAME": "キティ"}, "Hello キティ"},
		},
	})

	doTest(t, Test{
		"{NAME}",
		[]Expectation{
			{map[string]interface{}{"NAME": "leila"}, "leila"},
			{map[string]interface{}{"NAME": nil}, ""},
			{nil, ""},
		},
	})

	doTest(t, Test{
		"My name is { NAME}",
		[]Expectation{
			{map[string]interface{}{"NAME": "yoda"}, "My name is yoda"},
		},
	})

	doTest(t, Test{
		"My name is { NAME  }...",
		[]Expectation{
			{map[string]interface{}{"NAME": "chewy"}, "My name is chewy..."},
		},
	})

	doTest(t, Test{
		"Hey {A}, i'm your {B}!",
		[]Expectation{
			{map[string]interface{}{"A": "luke", "B": "father"}, "Hey luke, i'm your father!"},
		},
	})

	doTest(t, Test{
		`{
        NAME
        } is my name`,
		[]Expectation{
			{map[string]interface{}{"NAME": "chewy"}, "chewy is my name"},
		},
	})

	doTestException(
		t,
		"{VAR}",
		map[string]interface{}{"VAR": true},
		"toString: Unsupported type: bool",
	)
}

func BenchmarkVar(b *testing.B) {
	doBenchmarkExecute(
		b,
		"This is a {A}",
		"This is a benchmark",
		map[string]interface{}{"A": "benchmark"},
	)
}
