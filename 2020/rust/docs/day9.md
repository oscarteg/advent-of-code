# Day 9: Encoding Error

## Approach

**Part 1:** The XMAS cipher requires each number after the preamble to be a sum of two different numbers from the preceding preamble window. Find the first number that violates this rule by checking all pairs in the current window.

**Part 2:** Find a contiguous sequence of at least 2 numbers that sum to the invalid number from part 1. Return the sum of the minimum and maximum values in that sequence.

## Complexity

- **Part 1:** O(n × p²) where n is total numbers and p is preamble size (25). For each of n numbers, check p² pairs.
- **Part 2:** O(n²) with early termination. For each start position, iterate forward summing until target is found or exceeded.

## Commands

```bash
cd 2020/rust
cargo test --bin day9
cargo run -q --bin day9 -- 5 < input/2020/day9-example.txt   # the example uses a preamble of 5
cargo run -q --bin day9 < input/2020/day9.txt                # real inputs default to 25
```
