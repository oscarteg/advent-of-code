# Day 8: Handheld Halting

## Part 1
Parse each line into an `acc`/`jmp`/`nop` instruction with a signed argument.
Interpret the program, tracking which instructions have already run. The
moment the program counter is about to re-execute a visited instruction, it
has entered its infinite loop: return the accumulator at that point.
Complexity: O(n) time and space for a program of n instructions.

## Part 2
Exactly one `jmp`/`nop` is corrupted. Try swapping each one in turn (skipping
`acc` instructions) and re-run the patched program; the first swap that makes
it terminate (program counter runs past the last instruction) gives the
answer. Complexity: O(n) candidate swaps, each re-run in O(n), so O(n^2)
overall.

## Run and test
```bash
cargo run -q --bin day8 < input/2020/day8-example.txt
cargo test --bin day8
```
