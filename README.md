# messageformat

The messageformat package is a [Go](http://golang.org/) implementation of [ICU](http://site.icu-project.org/) [messages formatting](http://userguide.icu-project.org/formatparse/messages).

## Dependencies

The messageformat package depends on the [makeplural](http://github.com/gotnospirit/makeplural) package to compute named keys based on ICU rules.

## Test benchmarks

	Platform  Windows 10
	Processor Intel i7-4710HQ
	CPU Cores 4 @ 2.50Ghz
	Memory    8GB RAM

	> go test --bench=./
	BenchmarkLiteral-8                      10000000               189 ns/op
	BenchmarkSelectOrdinal-8                 2000000               625 ns/op
	BenchmarkNested-8                        5000000               341 ns/op
	BenchmarkEscaped-8                      10000000               213 ns/op
	BenchmarkPluralNonInteger-8              3000000               524 ns/op
	BenchmarkPluralLiteralize-8              2000000               764 ns/op
	BenchmarkPluralOther-8                   2000000               781 ns/op
	BenchmarkPluralExactValue-8              2000000               639 ns/op
	BenchmarkPluralOffsetExtension-8         2000000               640 ns/op
	BenchmarkSelect-8                        5000000               306 ns/op
	BenchmarkVar-8                           5000000               263 ns/op
	PASS
	ok      github.com/gotnospirit/messageformat    22.594s

## Todo

* More unit tests
* More doc
* `FormatArray(data []interface{}) (string, error)` ?
