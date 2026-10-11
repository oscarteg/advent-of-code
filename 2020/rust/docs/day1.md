# Day 1: Report Repair

## Problem Summary

Part 1: Find two entries in the expense report that sum to 2020 and multiply them together.

Part 2: Find three entries that sum to 2020 and multiply them together.

## Approach

Both parts use a HashSet for O(1) average-case lookups.

**Part 1:** Iterate through each entry x. Check if (2020 - x) exists in the set. Time: O(n).

**Part 2:** Iterate through each pair (x, y) and check if (2020 - x - y) exists. Time: O(n²).

## Complexity

- Time: O(n) for part 1, O(n²) for part 2
- Space: O(n) for the HashSet

## Commands

```bash
cd 2020/rust
cargo test --bin day1
cargo run -q --bin day1 < input/2020/day1-example.txt
cargo run -q --bin day1 < input/2020/day1.txt
```
