package messageformat

import (
	"testing"
	"time"
)

func TestParseDate(t *testing.T) {
	doTest(t, Test{
		"{varname, date, long}",
		[]Expectation{
			{data: map[string]any{"varname": time.Now()}, output: ""},
		},
	})
}
