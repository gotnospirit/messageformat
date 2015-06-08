package messageformat

import (
	"testing"
)

func doTestLiteral(t *testing.T, input string) {
	doTest(t, Test{
		input,
		[]Expectation{
			{output: input},
		},
	})
}

func TestLiteral(t *testing.T) {
	doTestLiteral(t, "")
	doTestLiteral(t, `\`)
	doTestLiteral(t, `\\`)
	doTestLiteral(t, `\\\`)
	doTestLiteral(t, `\q\`)
	doTestLiteral(t, `test\`)
	doTestLiteral(t, "\n")
	doTestLiteral(t, `\n`)
	doTestLiteral(t, " This is \n a string\"")
	doTestLiteral(t, "日本語")
	doTestLiteral(t, "Hello, 世界")
}

func BenchmarkLiteral(b *testing.B) {
	doBenchmarkExecute(b, "This is a benchmark", "This is a benchmark", nil)
}
