# Day 3: Toboggan Trajectory

## Problem

Navigate a repeating grid of trees following different slopes. The grid wraps horizontally using modulo arithmetic.

**Part 1:** Count trees on slope (right 3, down 1). Answer: 7

**Part 2:** Check multiple slopes and multiply tree counts: (1,1), (3,1), (5,1), (7,1), (1,2). Answer: 336

## Solution

Parse the input into a vector of line references. For each slope, iterate through rows incrementing by `down`, column by `right` (modulo line width), and count `#` characters.

**Complexity:** O(n) per slope, where n = number of rows.

## Commands

```bash
cd 2020/rust
cargo test --bin day3
cargo run -q --bin day3 < input/2020/day3-example.txt
cargo run -q --bin day3 < input/2020/day3.txt
```
