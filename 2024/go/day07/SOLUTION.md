# Day 7: Bridge Repair

Parse each line into a target and its operands. Recursively try every
operator assignment (`+`, `*` for part 1; also `||` concatenation for part
2), left to right with no precedence, pruning once the running total exceeds
the target. Sum the targets of every solvable equation.

## Complexity

O(equations x branches) time, where branches is up to `2^(n-1)` (part 1) or
`3^(n-1)` (part 2) for `n` operands; O(n) recursion depth per equation.

## Run / test

```bash
cd 2024/go && go test -race -count=1 ./day07
go run ./day07 < day07/testdata/example.txt
```
