# Day 20: Grove Positioning System

## Part 1

Parse the numbers, then mix the circular list once: for each number in its
original order, remove it from its current position and reinsert it `value`
positions forward (or backward), wrapping modulo `len - 1` since removing an
element from a list of `len` temporarily shortens the cycle it moves within.
Track positions by original index (not by value) because values can repeat.
Sum the values 1000, 2000, and 3000 positions after 0, wrapping around.

Complexity: O(n^2) per mix round (linear scan to locate each index), O(n)
rounds for part 2, so O(n^2) overall for the example-sized and real inputs.

## Part 2

Multiply every number by the decryption key (811589153) before mixing, then
mix 10 times instead of 1, reusing the same mixing and summing logic.

## Run

```bash
cargo run -q --bin day20 < input/day20-example.txt
```

## Test

```bash
cargo test --bin day20
```
