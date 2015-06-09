package messageformat

import (
	"testing"
)

func TestSelectOrdinal(t *testing.T) {
	doTest(t, Test{
		"The {FLOOR, selectordinal, one{#st} two{#nd} few{#rd} other{#th}} floor.",
		[]Expectation{
			{map[string]interface{}{"FLOOR": 0}, "The 0th floor."},
			{map[string]interface{}{"FLOOR": 1.0}, "The 1st floor."},
			{map[string]interface{}{"FLOOR": "2"}, "The 2nd floor."},
			{map[string]interface{}{"FLOOR": "3.00"}, "The 3.00rd floor."},
			{map[string]interface{}{"FLOOR": 4}, "The 4th floor."},
			{map[string]interface{}{"FLOOR": 101}, "The 101st floor."},
			{nil, "The #th floor."},
		},
	})

	doTestException(
		t,
		"{VAR,selectordinal,other{succeed}}",
		map[string]interface{}{"VAR": true},
		"toString: Unsupported type: bool",
	)
}

func BenchmarkSelectOrdinal(b *testing.B) {
	doBenchmarkExecute(
		b,
		"The {FLOOR, selectordinal, one{#st} two{#nd} few{#rd} other{#th}} floor.",
		"The 101st floor.",
		map[string]interface{}{"FLOOR": 101},
	)
}
