// @TODO(gotnospirit) add test on readKey, readChoice

package messageformat

import (
	"testing"
)

func TestSelect(t *testing.T) {
	doTest(t, Test{
		"{GENDER, select, male{He} female {She} other{They}} liked this.",
		[]Expectation{
			{map[string]interface{}{"GENDER": "male"}, "He liked this."},
			{map[string]interface{}{"GENDER": "female"}, "She liked this."},
			{nil, "They liked this."},
		},
	})

	doTest(t, Test{
		"{GENDER,select,male{He}female{She}other{They}} liked this.",
		[]Expectation{
			{map[string]interface{}{"GENDER": "male"}, "He liked this."},
			{map[string]interface{}{"GENDER": "female"}, "She liked this."},
			{nil, "They liked this."},
		},
	})

	doTest(t, Test{
		"{A, select, other{!#}}, and {B, select, other{#!}}",
		[]Expectation{
			{map[string]interface{}{"A": "black", "B": "mortimer"}, "!black, and mortimer!"},
		},
	})

	doTest(t, Test{
		"{A,select,other{#, and {B,select,other{#}}}}!",
		[]Expectation{
			{map[string]interface{}{"A": "black", "B": "mortimer"}, "black, and mortimer!"},
		},
	})

	doTest(t, Test{
		`{A,select,other{\##\, and {B,select,other{#\#}}}}`,
		[]Expectation{
			{map[string]interface{}{"A": "black", "B": "mortimer"}, `#black\, and mortimer#`},
		},
	})

	doTest(t, Test{
		`Hello {A,select,キティ{Kitty}other{World}}`,
		[]Expectation{
			{map[string]interface{}{"A": "キティ"}, `Hello Kitty`},
			{map[string]interface{}{"A": "世界"}, `Hello World`},
		},
	})

	doTest(t, Test{
		`{isTrueOrFalse, select, true {True} other {False}}`,
		[]Expectation{
			{map[string]interface{}{"isTrueOrFalse": true}, `True`},
			{map[string]interface{}{"isTrueOrFalse": false}, `False`},
		},
	})

	doTestFormatException(
		t,
		"{VAR,select,other{succeed}}",
		map[string]interface{}{"VAR": struct{}{}},
		"toString: Unsupported type: struct {}",
	)
}

func BenchmarkSelect(b *testing.B) {
	doBenchmarkExecute(
		b,
		"This is a {A, select, other{#}}",
		"This is a benchmark",
		map[string]interface{}{"A": "benchmark"},
	)
}
