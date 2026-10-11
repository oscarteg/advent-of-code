# Day 18: Operation Order

## Part 1
Tokenize each line into numbers, `+`, `*`, `(`, `)`. A recursive-descent parser evaluates
`+` and `*` with equal precedence, left to right, and recurses into parentheses as
sub-expressions. Sum the per-line results. Complexity: O(n) per line.

## Part 2
Same tokens and parser, but `+` now binds tighter than `*`: the grammar becomes
`expr := sum ('*' sum)*` and `sum := factor ('+' factor)*`, with parentheses still
recursing into a full `expr`. Complexity: O(n) per line.

## Run
```bash
cargo run -q --bin day18 < input/2020/day18-example.txt
```

## Test
```bash
cargo test --bin day18
```
