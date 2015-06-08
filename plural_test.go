// @TODO(gotnospirit) add test on readOffset

package messageformat

import (
	"testing"
)

func TestPlural(t *testing.T) {
	doTest(t, Test{
		"You have {NUM_TASKS, plural, one {one task} other {# tasks} =42 {the answer to the life, the universe and everything tasks}} remaining.",
		[]Expectation{
			{map[string]interface{}{"NUM_TASKS": -1}, "You have -1 tasks remaining."},
			{map[string]interface{}{"NUM_TASKS": 1}, "You have one task remaining."},
			{map[string]interface{}{"NUM_TASKS": 42}, "You have the answer to the life, the universe and everything tasks remaining."},
		},
	})

	doTest(t, Test{
		"{NUM_TASKS, plural, one {a} =1 {b} other {c}}",
		[]Expectation{
			{map[string]interface{}{"NUM_TASKS": 1}, "b"},
		},
	})

	doTest(t, Test{
		`{NUM, plural, one{a} other{b}}`,
		[]Expectation{
			{map[string]interface{}{"NUM": 1}, "a"},
			{map[string]interface{}{"NUM": "1"}, "a"},
		},
	})

	doTestException(
		t,
		"{NUM,plural,other{b}}",
		map[string]interface{}{"NUM": true},
		"toString: Unsupported type: bool",
	)
}

func BenchmarkPluralNonInteger(b *testing.B) {
	doBenchmarkExecute(
		b,
		"This is a {A, plural, other{benchmark}}",
		"This is a benchmark",
		map[string]interface{}{"A": "1"},
	)
}

func BenchmarkPluralLiteralize(b *testing.B) {
	doBenchmarkExecute(
		b,
		"This is a {A, plural, one{benchmark} other{}}",
		"This is a benchmark",
		map[string]interface{}{"A": 1},
	)
}

func BenchmarkPluralOther(b *testing.B) {
	doBenchmarkExecute(
		b,
		"This is a {A, plural, other{benchmark}}",
		"This is a benchmark",
		map[string]interface{}{"A": 1},
	)
}

func BenchmarkPluralExactValue(b *testing.B) {
	doBenchmarkExecute(
		b,
		"This is a {A, plural, =2{benchmark} other{}}",
		"This is a benchmark",
		map[string]interface{}{"A": 2},
	)
}

func TestPluralOffsetExtension(t *testing.T) {
	doTest(t, Test{
		"You {NUM_ADDS, plural, offset:1 =0{didnt add this to your profile} =1{added this to your profile} one{and one other person added this to their profile} other{and # others added this to their profiles}}.",
		[]Expectation{
			{map[string]interface{}{"NUM_ADDS": 0}, "You didnt add this to your profile."},
			{map[string]interface{}{"NUM_ADDS": 1}, "You added this to your profile."},
			{map[string]interface{}{"NUM_ADDS": 2}, "You and one other person added this to their profile."},
			{map[string]interface{}{"NUM_ADDS": 3}, "You and 2 others added this to their profiles."},
		},
	})
}

func BenchmarkPluralOffsetExtension(b *testing.B) {
	doBenchmarkExecute(
		b,
		"This is a {A, plural, offset:1 =2{benchmark} other{}}",
		"This is a benchmark",
		map[string]interface{}{"A": 2},
	)
}
