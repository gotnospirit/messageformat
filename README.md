# messageformat

The messageformat package is a [Go](http://golang.org/) implementation of [ICU](http://site.icu-project.org/) [messages formatting](http://userguide.icu-project.org/formatparse/messages).

## Dependencies

The messageformat package depends on the [makeplural](http://github.com/gotnospirit/makeplural) package to compute named keys based on ICU rules.

## Tests

	> go test -v ./...
<details>
  <summary>Details</summary>

	=== RUN   TestLiteral
	- Got expected value <>
	- Got expected value <\>
	- Got expected value <\\>
	- Got expected value <\\\>
	- Got expected value <\q\>
	- Got expected value <test\>
	- Got expected value <
	>
	- Got expected value <\n>
	- Got expected value < This is
	a string">
	- Got expected value <日本語>
	- Got expected value <Hello, 世界>
	--- PASS: TestLiteral (0.00s)
	=== RUN   TestSetPluralFunction
	- Got expected exception <PluralFunctionRequired>
	- Got expected value <2>
	--- PASS: TestSetPluralFunction (0.00s)
	=== RUN   TestSelectOrdinal
	- Got expected value <The 0th floor.>
	- Got expected value <The 1st floor.>
	- Got expected value <The 2nd floor.>
	- Got expected value <The 3.00rd floor.>
	- Got expected value <The 4th floor.>
	- Got expected value <The 101st floor.>
	- Got expected value <The #th floor.>
	- Got expected exception <toString: Unsupported type: struct {}>
	--- PASS: TestSelectOrdinal (0.00s)
	=== RUN   TestParseException
	- Got expected exception <ParseError: `UnbalancedBraces` at 1>
	- Got expected exception <ParseError: `UnbalancedBraces` at 3>
	- Got expected exception <ParseError: `InvalidFormat` at 1>
	- Got expected exception <ParseError: `UnbalancedBraces` at 2>
	- Got expected exception <ParseError: `UnbalancedBraces` at 12>
	- Got expected exception <ParseError: `UnbalancedBraces` at 17>
	- Got expected exception <ParseError: `UnbalancedBraces` at 20>
	- Got expected exception <ParseError: `UnbalancedBraces` at 20>
	- Got expected exception <ParseError: `UnbalancedBraces` at 17>
	- Got expected exception <ParseError: `UnbalancedBraces` at 19>
	- Got expected exception <ParseError: `InvalidExpr` at 1>
	- Got expected exception <ParseError: `InvalidExpr` at 3>
	- Got expected exception <ParseError: `MissingVarName` at 1>
	- Got expected exception <ParseError: `MissingVarName` at 8>
	- Got expected exception <ParseError: `MissingVarName` at 5>
	- Got expected exception <ParseError: `MissingVarName` at 2>
	- Got expected exception <ParseError: `InvalidFormat` at 3>
	- Got expected exception <ParseError: `InvalidFormat` at 3>
	- Got expected exception <ParseError: `InvalidFormat` at 4>
	- Got expected exception <ParseError: `InvalidFormat` at 4>
	- Got expected exception <ParseError: `InvalidFormat` at 1>
	- Got expected exception <ParseError: `InvalidFormat` at 1>
	- Got expected exception <ParseError: `InvalidFormat` at 9>
	- Got expected exception <ParseError: `UnknownType: `SELECT`` at 11>
	- Got expected exception <ParseError: `MissingChoiceName` at 12>
	- Got expected exception <ParseError: `MissingChoiceName` at 22>
	- Got expected exception <ParseError: `MissingChoiceName` at 19>
	- Got expected exception <ParseError: `MissingChoiceName` at 29>
	- Got expected exception <ParseError: `MissingChoiceName` at 12>
	- Got expected exception <ParseError: `MissingChoiceName` at 22>
	- Got expected exception <ParseError: `MissingChoiceName` at 20>
	- Got expected exception <ParseError: `MissingChoiceName` at 21>
	- Got expected exception <ParseError: `MissingChoiceName` at 31>
	- Got expected exception <ParseError: `MalformedOption` at 10>
	- Got expected exception <ParseError: `MalformedOption` at 17>
	- Got expected exception <ParseError: `MalformedOption` at 10>
	- Got expected exception <ParseError: `MissingChoiceContent` at 16>
	- Got expected exception <ParseError: `MissingChoiceContent` at 23>
	- Got expected exception <ParseError: `MissingChoiceContent` at 16>
	- Got expected exception <ParseError: `MissingMandatoryChoice` at 28>
	- Got expected exception <ParseError: `MissingMandatoryChoice` at 35>
	- Got expected exception <ParseError: `MissingMandatoryChoice` at 28>
	- Got expected exception <ParseError: `UnexpectedExtension` at 18>
	- Got expected exception <ParseError: `UnexpectedExtension` at 25>
	- Got expected exception <ParseError: `UnsupportedExtension: `factor`` at 18>
	- Got expected exception <ParseError: `MissingOffsetValue` at 19>
	- Got expected exception <ParseError: `BadCast` at 23>
	- Got expected exception <ParseError: `BadCast` at 20>
	- Got expected exception <ParseError: `BadCast` at 22>
	- Got expected exception <ParseError: `InvalidOffsetValue` at 21>
	--- PASS: TestParseException (0.00s)
	=== RUN   TestNested
	- Got expected value <1>
	- Got expected value <deep in the heart.>
	- Got expected value <deep in the heart.>
	- Got expected value <I have 0 friends but one enemy..>
	- Got expected value <I have # friends but # enemies..>
	--- PASS: TestNested (0.00s)
	=== RUN   TestEscaped
	- Got expected value <#>
	- Got expected value <\\>
	- Got expected value <{>
	- Got expected value <}>
	- Got expected value <{ 5 is a # }>
	- Got expected value <{{{4}}}>
	- Got expected value <日{本}語>
	- Got expected value <he\\#ll\\\{o\\} ##!>
	--- PASS: TestEscaped (0.00s)
	=== RUN   TestNonAscii
	- Got expected value <猫 キティ。。。>
	--- PASS: TestNonAscii (0.00s)
	=== RUN   TestMultiline
	- Got expected value <He>
	- Got expected value <She>
	- Got expected value <They>
	- Got expected value <He found 1 result.>
	- Got expected value <She found 1 result.>
	- Got expected value <He found
	    2 results in 1 category !.>
	- Got expected value <They found
	    2 results in 2 categories !.>
	- Got expected value <1 result, 2 categories.>
	- Got expected value <2 results, 1 category.>
	- Got expected value <2 results, 2 categories.>
	- Got expected value <He found 1 result in 2 categories.>
	- Got expected value <She found 1 result in 2 categories.>
	- Got expected value <He found 2 results in 1 category.>
	- Got expected value <They found 2 results in 2 categories.>
	--- PASS: TestMultiline (0.00s)
	=== RUN   TestRegister
	- Got expected exception <ParserAlreadyRegistered>
	- Got expected exception <ParserAlreadyRegistered>
	- Got expected exception <ParserAlreadyRegistered>
	- Got expected exception <ParserAlreadyRegistered>
	- Got expected exception <ParseError: `UndefinedParseFunc: `noparse`` at 10>
	- Got expected exception <UndefinedFormatFunc: `noeval`>
	--- PASS: TestRegister (0.00s)
	=== RUN   TestPlural
	- Got expected value <You have -1 tasks remaining.>
	- Got expected value <You have one task remaining.>
	- Got expected value <You have the answer to the life, the universe and everything tasks remaining.>
	- Got expected value <b>
	- Got expected value <a>
	- Got expected value <a>
	- Got expected exception <toString: Unsupported type: struct {}>
	--- PASS: TestPlural (0.00s)
	=== RUN   TestPluralOffsetExtension
	- Got expected value <You didnt add this to your profile.>
	- Got expected value <You added this to your profile.>
	- Got expected value <You and one other person added this to their profile.>
	- Got expected value <You and 2 others added this to their profiles.>
	--- PASS: TestPluralOffsetExtension (0.00s)
	=== RUN   TestSelect
	- Got expected value <He liked this.>
	- Got expected value <She liked this.>
	- Got expected value <They liked this.>
	- Got expected value <He liked this.>
	- Got expected value <She liked this.>
	- Got expected value <They liked this.>
	- Got expected value <!black, and mortimer!>
	- Got expected value <black, and mortimer!>
	- Got expected value <#black\, and mortimer#>
	- Got expected value <Hello Kitty>
	- Got expected value <Hello World>
	- Got expected value <True>
	- Got expected value <False>
	- Got expected exception <toString: Unsupported type: struct {}>
	--- PASS: TestSelect (0.00s)
	=== RUN   TestIsWhitespace
	--- PASS: TestIsWhitespace (0.00s)
	=== RUN   TestWhitespace
	Successfully returns `h`, 2
	Successfully returns `h`, 2
	Successfully returns ` `, 2
	--- PASS: TestWhitespace (0.00s)
	=== RUN   TestToString
	Successfully returns the expected value: ``
	Successfully returns the expected value: ``
	Successfully returns the expected value: `I am a string`
	Successfully returns the expected value: `42`
	Successfully returns the expected value: `0.305`
	Successfully returns the expected value: `true`
	--- PASS: TestToString (0.00s)
	=== RUN   TestToStringNumericTypes
	Successfully returns the expected value: `255`
	Successfully returns the expected value: `123456`
	Successfully returns the expected value: `255`
	Successfully returns the expected value: `65535`
	Successfully returns the expected value: `4294967295`
	Successfully returns the expected value: `18446744073709551615`
	Successfully returns the expected value: `-123456`
	Successfully returns the expected value: `-128`
	Successfully returns the expected value: `127`
	Successfully returns the expected value: `-32768`
	Successfully returns the expected value: `32767`
	Successfully returns the expected value: `-2147483648`
	Successfully returns the expected value: `2147483647`
	Successfully returns the expected value: `-9223372036854775808`
	Successfully returns the expected value: `9223372036854775807`
	Successfully returns the expected value: `3.14`
	Successfully returns the expected value: `0.000000000314`
	Successfully returns the expected value: `(1.23+9.87i)`
	Successfully returns the expected value: `(1.23+9.87i)`
	Successfully returns the expected value: `(1.23+9.87i)`
	Successfully returns the expected value: `97`
	Successfully returns the expected value: `075bcd15`
	--- PASS: TestToStringNumericTypes (0.00s)
	=== RUN   TestToStringBool
	Successfully returns the expected value: `true`
	Successfully returns the expected value: `false`
	--- PASS: TestToStringBool (0.00s)
	=== RUN   TestToStringTimeDuration
	Successfully returns the expected value: `87672h0m0s`
	--- PASS: TestToStringTimeDuration (0.00s)
	=== RUN   TestToStringStringer
	Successfully returns the expected value: `1`
	Successfully returns the expected value: `2`
	--- PASS: TestToStringStringer (0.00s)
	=== RUN   TestVar
	- Got expected value <Hello キティ>
	- Got expected value <leila>
	- Got expected value <>
	- Got expected value <>
	- Got expected value <My name is yoda>
	- Got expected value <My name is chewy...>
	- Got expected value <Hey luke, i'm your father!>
	- Got expected value <chewy is my name>
	- Got expected exception <toString: Unsupported type: struct {}>
	--- PASS: TestVar (0.00s)
</details>

	PASS
	ok      github.com/gotnospirit/messageformat    0.252s

## Test benchmarks

	> go test -bench=.
<details>
  <summary>Details</summary>

	goos: windows
	goarch: amd64
	pkg: github.com/gotnospirit/messageformat
	cpu: Intel(R) Core(TM) i7-4710HQ CPU @ 2.50GHz
	BenchmarkLiteral-8                       6423349               184.6 ns/op
	BenchmarkSelectOrdinal-8                 2440293               473.0 ns/op
	BenchmarkNested-8                        4074882               296.0 ns/op
	BenchmarkEscaped-8                       5898448               200.1 ns/op
	BenchmarkPluralNonInteger-8              2866503               419.0 ns/op
	BenchmarkPluralLiteralize-8              2360140               506.1 ns/op
	BenchmarkPluralOther-8                   2338374               511.0 ns/op
	BenchmarkPluralExactValue-8              2636702               458.4 ns/op
	BenchmarkPluralOffsetExtension-8         2514459               456.6 ns/op
	BenchmarkSelect-8                        4575456               263.1 ns/op
	BenchmarkVar-8                           5351002               222.3 ns/op
</details>

  PASS
  ok      github.com/gotnospirit/messageformat    17.449s

## Test coverage

	> go test -coverprofile cover.out
	PASS
	coverage: 93.5% of statements
	ok      github.com/gotnospirit/messageformat    0.255s

## Golang versions

This package was originally made using Go v1.4, and was successfully tested with v1.18.1

## Todo

* More unit tests
* More doc
* `FormatArray(data []interface{}) (string, error)` ?
