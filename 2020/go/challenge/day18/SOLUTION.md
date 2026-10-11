# Day 18: Operation Order

## Part A
Tokenize each line into numbers, `+`, `*`, `(`, `)`. A recursive-descent parser evaluates
`+` and `*` with equal precedence, left to right, recursing into parentheses as
sub-expressions. Sum the per-line results. Complexity: O(n) per line.

## Part B
Same tokens and parser, but `+` now binds tighter than `*`: the grammar becomes
`expr := sum ('*' sum)*` and `sum := factor ('+' factor)*`, with parentheses still
recursing into a full `expr`. Complexity: O(n) per line.

## Commands
```bash
cd 2020/go
go test -race ./challenge/day18
go run . 18 < ../rust/input/2020/day18-example.txt
go run . 18 < ../rust/input/2020/day18.txt
```
