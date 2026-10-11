# Day 25: Full of Hot Air

## Part 1

Each line is a base-5 "SNAFU" number with digits `2`, `1`, `0`, `-` (-1) and
`=` (-2). Parsing folds left to right: `value = value * 5 + digit`. Sum all
parsed values as `i64`, then convert the total back to SNAFU by repeatedly
taking `value mod 5`, mapping remainders 3 and 4 to `=`/`-` with a carry,
and reversing the collected digits.

Complexity: O(n) in total input length for parsing and summing; O(log sum)
for the final conversion. Single part puzzle.

## Run & test

```bash
zig fmt --check src/day_25.zig
zig test src/day_25.zig
zig run src/day_25.zig < src/examples/day_25.txt
```
