# Day 8: Handheld Halting

## Part A
Parse each line into an `acc`/`jmp`/`nop` instruction with a signed argument.
Run the program while tracking visited instruction indices; the moment the
program counter is about to re-execute a visited instruction it has entered
its infinite loop, so return the accumulator at that point. O(n) time/space.

## Part B
Exactly one `jmp`/`nop` is corrupted. Try swapping each one in turn and
re-run; the first swap that makes the program terminate (pc runs past the
last instruction) gives the answer. O(n) swaps, each re-run in O(n): O(n^2).

## Run and test
```bash
go run . 8 < ../rust/input/2020/day8.txt
go test -race ./challenge/day8
```
