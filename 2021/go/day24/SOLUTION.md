# Day 24: Arithmetic Logic Unit

Lines are parsed into typed instructions and an interpreter runs them over the
registers `w`, `x`, `y` and `z`, rejecting the crashes the puzzle warns about:
`div` by zero and `mod` with a negative left or non-positive right operand.

MONAD repeats one 18-instruction block per digit, parameterised by `div` (1 or
26), `addX` and `addY`; each block is matched against that template, so other
shapes are rejected. `z` is a base-26 stack: a `div 1` block pushes
`digit + addY` (its `addX` exceeds every digit, so its comparison always fails),
and a `div 26` block pops it, leaving `z` small only when
`digits[pop] == digits[push] + addY[push] + addX[pop]`. Part 1 takes the largest
digit pair satisfying each constraint, part 2 the smallest, and the interpreter
validates both: `O(n)` time and space, not 9^14 candidates.

`testdata/example.txt` is a synthetic 14-digit MONAD: the puzzle ships only
small ALU examples (negation, triple check, binary split; all tested), no MONAD.

From `2021/go`:

```sh
go run ./day24 < day24/testdata/example.txt
# Part 1: 99999596891916
# Part 2: 49411151121911
go test -race ./day24
go vet ./day24
```
