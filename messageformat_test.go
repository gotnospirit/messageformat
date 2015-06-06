package messageformat

import (
	"fmt"
	"testing"
)

// doTest(t, Test{
// "You have {NUM_TASKS, plural, zero {no task} one {one task} two {two tasks} few{few tasks} many {many tasks} other {# tasks} =42 {the answer to the life, the universe and everything tasks}} remaining.",
// []Expectation{
// {map[string]interface{}{"NUM_TASKS": -1}, "You have -1 tasks remaining."},
// {map[string]interface{}{"NUM_TASKS": 0}, "You have no task remaining."},
// {map[string]interface{}{"NUM_TASKS": 1}, "You have one task remaining."},
// {map[string]interface{}{"NUM_TASKS": 2}, "You have two tasks remaining."},
// {map[string]interface{}{"NUM_TASKS": 3}, "You have few tasks remaining."},
// {map[string]interface{}{"NUM_TASKS": 6}, "You have many tasks remaining."},
// {map[string]interface{}{"NUM_TASKS": 15}, "You have 15 tasks remaining."},
// {map[string]interface{}{"NUM_TASKS": 42}, "You have the answer to the life, the universe and everything tasks remaining."},
// },
// })

func doParse(input string) *MessageFormat {
	o, _ := New()
	mf, err := o.Parse(input)
	if nil != err {
		return nil
	}
	return mf
}

func TestSetPluralFunction(t *testing.T) {
	mf := doParse(`{N,plural,one{1}other{2}}`)
	if nil == mf {
		t.Errorf("Unexpected parse failure")
	} else {
		// checks we can't unset the default plural function
		err := mf.SetPluralFunction(nil)
		doTestError(t, "PluralFunctionRequired", err)

		data := map[string]interface{}{"N": 1}
		result, _ := mf.FormatMap(data)
		if "1" != result {
			t.Errorf("Unexpected format result")
		} else {
			// checks we can set a new plural function and get a different result
			err := mf.SetPluralFunction(func(float64, bool) string {
				return "other"
			})

			if nil != err {
				t.Errorf("Unexpected error <%s>", err)
			} else {
				result, _ := mf.FormatMap(data)

				if "2" != result {
					t.Errorf("Unexpected format result")
				} else if testing.Verbose() {
					fmt.Printf("- Got expected value <%s>\n", result)
				}
			}
		}
	}
}
