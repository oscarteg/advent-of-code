# Day 17: Chronospatial Computer

A three-register virtual machine with eight 3-bit opcodes. `parse` returns the
registers and program as a `machine`; `run` interprets it with any starting A.

## Approach

Part 1 interprets the program once with the given A and joins the output with
commas.

Part 2 finds the lowest A that makes the program print itself: such a program
emits one value per loop and then drops the low three bits of A, so A is built
three bits at a time from the last value backwards -- for every surviving
prefix, try all eight next digits and keep those whose run equals the program's
remaining tail. It reports 0 when no A reproduces the program, as for the
part 1 example.

## Complexity

Part 1 is O(outputs); part 2 runs the interpreter eight times per surviving
candidate per program value.

## Run / test

```bash
cd 2024/go && go test -race -count=1 ./day17
go run ./day17 < day17/testdata/example.txt
```
