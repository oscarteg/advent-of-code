# Day 6: Tuning Trouble

## Problem

Find the start-of-packet and start-of-message markers in a datastream.

- **Part 1:** Find the position where the first sequence of 4 distinct characters ends.
- **Part 2:** Find the position where the first sequence of 14 distinct characters ends.

## Approach

Both parts use a sliding window algorithm:
1. Maintain a window of the required size (4 for part 1, 14 for part 2)
2. For each position, check if all characters in the window are distinct
3. Return the position when all characters are unique

A HashSet tracks distinct characters efficiently in O(1) average case.

## Complexity

- Time: O(n * window_size) where n is the input length
- Space: O(window_size) for the HashSet

## Commands

```bash
cd 2022/rust
cargo test --bin day6
cargo run -q --bin day6 < input/day6-example.txt
```
