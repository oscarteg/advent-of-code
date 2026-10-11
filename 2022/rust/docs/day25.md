# Day 25: Full of Hot Air

## Part 1

Parse each line as a base-5 "SNAFU" number whose digits are `2`, `1`, `0`,
`-` (-1) and `=` (-2). Fold left to right, multiplying the running total by
5 and adding the digit's value. Sum all numbers as plain `i64`, then convert
back to SNAFU by repeatedly taking the value mod 5, mapping remainders 3 and
4 to `=`/`-` with a carry into the next digit, and reversing the collected
digits.

Complexity: O(n) in total input length for parsing and summing; O(log sum)
for the final conversion. Single part puzzle.

## Run & test

```bash
cargo test --bin day25
cargo run -q --bin day25 < input/day25-example.txt
```
