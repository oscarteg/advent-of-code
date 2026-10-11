# Day 5: Supply Stacks

Parse a drawing of crate stacks (crate letters sit at fixed columns, every 4th
character starting at index 1) followed by a blank line and a list of
`move N from A to B` instructions.

## Part 1
Simulate each move with a `CrateMover 9000`: pick crates off stack `A` one at a
time and push each directly onto stack `B`, which reverses the order of the
moved group. Answer is the top crate of every stack, read left to right.

## Part 2
Same moves, but a `CrateMover 9001` lifts the whole group of `N` crates at
once and places it on stack `B` preserving order (no reversal).

## Complexity
`O(n)` time and space in the total number of crate-moving operations; each
move touches at most `N` crates once.

## Commands
```bash
cargo test --bin day5
cargo run -q --bin day5 < input/day5-example.txt
```
