// @TODO(gotnospirit) add test on parseExpression, parse, NewWithCulture

package messageformat

import (
	"fmt"
	"testing"
)

type Test struct {
	input   string
	expects []Expectation
}

type Expectation struct {
	data   map[string]interface{}
	output string
}

func doTest(t *testing.T, data Test) {
	o := NewParser()
	f := NewFormatter()

	pt, err := o.Parse(data.input)

	if nil != err {
		t.Errorf("`%s` threw <%s>", data.input, err)
	}

	for _, ex := range data.expects {
		result, err := f.FormatMap(pt, ex.data)
		if nil != err {
			t.Errorf("`%s` threw <%s>", data.input, err)
		} else if result != ex.output {
			t.Errorf("Expecting <%v> but got <%v>", ex.output, result)
		} else if testing.Verbose() {
			fmt.Printf("- Got expected value <%s>\n", result)
		}
	}

}

// Executes test code, expecting an exception when calling format.FormatMap
func doTestFormatException(t *testing.T, input string, data map[string]interface{}, expected string) {
	o := NewParser()
	f := NewFormatter()

	pt, err := o.Parse(input)
	if err != nil {
		t.Errorf("`%s` threw <%s>", input, err)
	}

	_, err = f.FormatMap(pt, data)
	doTestCompileError(t, input, expected, err)
}

// Executes test code, expecting an exception when calling parser.Parse
func doTestParseException(t *testing.T, input, expected string) {
	o := NewParser()

	_, err := o.Parse(input)
	doTestCompileError(t, input, expected, err)
}

func doTestCompileError(t *testing.T, input, expected string, err error) {
	if nil == err {
		t.Errorf("`%s` should threw <%s> but got none", input, expected)
	} else if err.Error() != expected {
		t.Errorf("`%s` should threw <%s> but got <%s>", input, expected, err.Error())
	} else if testing.Verbose() {
		fmt.Printf("- Got expected exception <%s>\n", expected)
	}
}

func doTestError(t *testing.T, expected string, err error) {
	if nil == err {
		t.Errorf("Expecting exception <%s> but got none", expected)
	} else if err.Error() != expected {
		t.Errorf("Expecting exception <%s> but got <%s>", expected, err.Error())
	} else if testing.Verbose() {
		fmt.Printf("- Got expected exception <%s>\n", expected)
	}
}

func doBenchmarkExecute(b *testing.B, input, expected string, data map[string]interface{}) (pt *ParseTree) {
	o := NewParser()
	f := NewFormatter()

	pt, _ = o.Parse(input)

	for n := 0; n < b.N; n++ {
		result, _ := f.FormatMap(pt, data)

		if result != expected {
			b.Errorf("Expecting <%s> but got <%s>", expected, result)
		}
	}
	return
}

func TestParseException(t *testing.T) {
	doTestParseException(t, "{", "ParseError: `UnbalancedBraces` at 1")
	doTestParseException(t, "{N}}", "ParseError: `UnbalancedBraces` at 3")
	doTestParseException(t, `{\}`, "ParseError: `InvalidFormat` at 1")
	doTestParseException(t, `{N`, "ParseError: `UnbalancedBraces` at 2")
	doTestParseException(t, `{ N , select`, "ParseError: `UnbalancedBraces` at 12")
	doTestParseException(t, `{N, select, other`, "ParseError: `UnbalancedBraces` at 17")
	doTestParseException(t, `{N, plural, other{#}`, "ParseError: `UnbalancedBraces` at 20")
	doTestParseException(t, `{N, plural, other{#{`, "ParseError: `UnbalancedBraces` at 20")
	doTestParseException(t, `{N, plural, other`, "ParseError: `UnbalancedBraces` at 17")
	doTestParseException(t, "{N, plural, offset:", "ParseError: `UnbalancedBraces` at 19")

	doTestParseException(t, `{{}`, "ParseError: `InvalidExpr` at 1")
	doTestParseException(t, `{N,{}`, "ParseError: `InvalidExpr` at 3")

	doTestParseException(t, `{}`, "ParseError: `MissingVarName` at 1")
	doTestParseException(t, `{       }`, "ParseError: `MissingVarName` at 8")
	doTestParseException(t, `{    ,   }`, "ParseError: `MissingVarName` at 5")
	doTestParseException(t, `{ , , }`, "ParseError: `MissingVarName` at 2")

	doTestParseException(t, `{NA-ME}`, "ParseError: `InvalidFormat` at 3")
	doTestParseException(t, `{N A M E}`, "ParseError: `InvalidFormat` at 3")
	doTestParseException(t, `{NAMé}`, "ParseError: `InvalidFormat` at 4")
	doTestParseException(t, `{NAMÉ}`, "ParseError: `InvalidFormat` at 4")
	doTestParseException(t, `{\}NAME`, "ParseError: `InvalidFormat` at 1")
	doTestParseException(t, `{なまえ}`, "ParseError: `InvalidFormat` at 1")

	doTestParseException(t, `{ N, sel ect, other {#} }`, "ParseError: `InvalidFormat` at 9")
	doTestParseException(t, `{ N, SELECT, other {#} }`, "ParseError: `UnknownType: `SELECT`` at 11")

	doTestParseException(t, `{N, select, {#} other {#}}`, "ParseError: `MissingChoiceName` at 12")
	doTestParseException(t, `{N, select, other {#} {#}}`, "ParseError: `MissingChoiceName` at 22")
	doTestParseException(t, `{N, selectordinal, {#} other {#}}`, "ParseError: `MissingChoiceName` at 19")
	doTestParseException(t, `{N, selectordinal, other {#} {#}}`, "ParseError: `MissingChoiceName` at 29")
	doTestParseException(t, `{N, plural, {#} other {#}}`, "ParseError: `MissingChoiceName` at 12")
	doTestParseException(t, `{N, plural, other {#} {#}}`, "ParseError: `MissingChoiceName` at 22")
	doTestParseException(t, `{N, plural, offset:1{#} other {#}}`, "ParseError: `MissingChoiceName` at 20")
	doTestParseException(t, `{N, plural, offset:1 {#} other {#}}`, "ParseError: `MissingChoiceName` at 21")
	doTestParseException(t, `{N, plural, offset:1 other {#} {#}}`, "ParseError: `MissingChoiceName` at 31")

	doTestParseException(t, `{N, select}`, "ParseError: `MalformedOption` at 10")
	doTestParseException(t, `{N, selectordinal}`, "ParseError: `MalformedOption` at 17")
	doTestParseException(t, `{N, plural}`, "ParseError: `MalformedOption` at 10")

	doTestParseException(t, `{N, select, one two{She} other{Other}}`, "ParseError: `MissingChoiceContent` at 16")
	doTestParseException(t, `{N, selectordinal, one two{She} other{Other}}`, "ParseError: `MissingChoiceContent` at 23")
	doTestParseException(t, `{N, plural, one two{She} other{Other}}`, "ParseError: `MissingChoiceContent` at 16")

	doTestParseException(t, `{N, select, one{He} two{She}}`, "ParseError: `MissingMandatoryChoice` at 28")
	doTestParseException(t, `{N, selectordinal, one{He} two{She}}`, "ParseError: `MissingMandatoryChoice` at 35")
	doTestParseException(t, `{N, plural, one{He} two{She}}`, "ParseError: `MissingMandatoryChoice` at 28")

	doTestParseException(t, "{N, select, offset:1 one{#} other {#}}", "ParseError: `UnexpectedExtension` at 18")
	doTestParseException(t, "{N, selectordinal, offset:1 one{#} other {#}}", "ParseError: `UnexpectedExtension` at 25")
	doTestParseException(t, "{N, plural, factor:1 one{#} other {#}}", "ParseError: `UnsupportedExtension: `factor`` at 18")

	doTestParseException(t, "{N, plural, offset:}", "ParseError: `MissingOffsetValue` at 19")
	doTestParseException(t, "{N, plural, offset: one{#} other {#}}", "ParseError: `BadCast` at 23")
	doTestParseException(t, "{N, plural, offset:A one{#} other {#}}", "ParseError: `BadCast` at 20")
	doTestParseException(t, "{N, plural, offset:1.0 one{#} other {#}}", "ParseError: `BadCast` at 22")
	doTestParseException(t, "{N, plural, offset:-1 one{#} other {#}}", "ParseError: `InvalidOffsetValue` at 21")
}

func Test_All(t *testing.T) {
	doTest(t, Test{
		"Hello World. My name is Tyler and I have {count} dogs",
		[]Expectation{
			{map[string]interface{}{
				"count": 3,
			}, "Hello World. My name is Tyler and I have 3 dogs"},
		},
	})

}

func TestNested(t *testing.T) {
	doTest(t, Test{
		"{PLUR1, plural, one {1} other {{SEL2, select, other {deep in the heart.}}}}",
		[]Expectation{
			{map[string]interface{}{"PLUR1": 1}, "1"},
			{map[string]interface{}{"SEL2": 1}, "deep in the heart."},
			{nil, "deep in the heart."},
		},
	})

	doTest(t, Test{
		"I have {FRIENDS, plural, one{one friend} other{# friends but {ENEMIES, plural, one{one enemy} other{# enemies}}.}}.",
		[]Expectation{
			{map[string]interface{}{"FRIENDS": 0, "ENEMIES": "1"}, "I have 0 friends but one enemy.."},
			{nil, "I have # friends but # enemies.."},
		},
	})
}

func BenchmarkNested(b *testing.B) {
	doBenchmarkExecute(
		b,
		"This is a {A,select,other{{B,select,other{benchmark}}}}",
		"This is a benchmark",
		map[string]interface{}{},
	)
}

func TestEscaped(t *testing.T) {
	doTest(t, Test{
		`\#`,
		[]Expectation{
			{output: `#`},
		},
	})

	doTest(t, Test{
		`\\`,
		[]Expectation{
			{output: `\\`},
		},
	})

	doTest(t, Test{
		`\{`,
		[]Expectation{
			{output: `{`},
		},
	})

	doTest(t, Test{
		`\}`,
		[]Expectation{
			{output: `}`},
		},
	})

	doTest(t, Test{
		`\{ {S, select, other{# is a \#}} \}`,
		[]Expectation{
			{map[string]interface{}{"S": 5}, "{ 5 is a # }"},
		},
	})

	doTest(t, Test{
		`\{\{\{{test, plural, other{#}}\}\}\}`,
		[]Expectation{
			{map[string]interface{}{"test": 4}, "{{{4}}}"},
		},
	})

	doTest(t, Test{
		`日\{本\}語`,
		[]Expectation{
			{output: "日{本}語"},
		},
	})

	doTest(t, Test{
		`he\\#ll\\\{o\\} \##!`,
		[]Expectation{
			{output: `he\\#ll\\\{o\\} ##!`},
		},
	})
}

func BenchmarkEscaped(b *testing.B) {
	doBenchmarkExecute(
		b,
		`日\{本\}語`,
		`日{本}語`,
		map[string]interface{}{},
	)
}

func TestNonAscii(t *testing.T) {
	doTest(t, Test{
		`猫 {N}。。。`,
		[]Expectation{
			{map[string]interface{}{"N": "キティ"}, "猫 キティ。。。"},
		},
	})
}

func TestMultiline(t *testing.T) {
	doTest(t, Test{
		`{GENDER, select,
    male {He}
  female {She}
   other {They}
}`,
		[]Expectation{
			{map[string]interface{}{"GENDER": "male"}, "He"},
			{map[string]interface{}{"GENDER": "female"}, "She"},
			{nil, "They"},
		},
	})

	doTest(t, Test{
		`{GENDER, select,
    male {He}
  female {She}
   other {They}
} found {NUM_RESULTS, plural,
            one
            {1 result}
          other {
          # results in {NUM_CATEGORIES, plural,
                  one {1 category}
                other {# categories}
             } !}
        }.`,
		[]Expectation{
			{map[string]interface{}{"GENDER": "male", "NUM_RESULTS": 1, "NUM_CATEGORIES": 2}, "He found 1 result."},
			{map[string]interface{}{"GENDER": "female", "NUM_RESULTS": 1, "NUM_CATEGORIES": 2}, "She found 1 result."},
			{map[string]interface{}{"GENDER": "male", "NUM_RESULTS": 2, "NUM_CATEGORIES": 1}, "He found \n          2 results in 1 category !."},
			{map[string]interface{}{"NUM_RESULTS": 2, "NUM_CATEGORIES": 2}, "They found \n          2 results in 2 categories !."},
		},
	})

	doTest(t, Test{
		`{NUM_RESULTS, plural,
            one
            {1 result}
          other {# results}
        }, {NUM_CATEGORIES, plural,
                  one {1 category}
                other {# categories}
             }.`,
		[]Expectation{
			{map[string]interface{}{"NUM_RESULTS": 1, "NUM_CATEGORIES": 2}, "1 result, 2 categories."},
			{map[string]interface{}{"NUM_RESULTS": 2, "NUM_CATEGORIES": 1}, "2 results, 1 category."},
			{map[string]interface{}{"NUM_RESULTS": 2, "NUM_CATEGORIES": 2}, "2 results, 2 categories."},
		},
	})

	doTest(t, Test{
		`{GENDER, select,
    male {He}
  female {She}
   other {They}
} found {NUM_RESULTS, plural,
            one
            {1 result}
          other {# results}
        } in {NUM_CATEGORIES, plural,
                  one {1 category}
                other {# categories}
             }.`,
		[]Expectation{
			{map[string]interface{}{"GENDER": "male", "NUM_RESULTS": 1, "NUM_CATEGORIES": 2}, "He found 1 result in 2 categories."},
			{map[string]interface{}{"GENDER": "female", "NUM_RESULTS": 1, "NUM_CATEGORIES": 2}, "She found 1 result in 2 categories."},
			{map[string]interface{}{"GENDER": "male", "NUM_RESULTS": 2, "NUM_CATEGORIES": 1}, "He found 2 results in 1 category."},
			{map[string]interface{}{"NUM_RESULTS": 2, "NUM_CATEGORIES": 2}, "They found 2 results in 2 categories."},
		},
	})
}
