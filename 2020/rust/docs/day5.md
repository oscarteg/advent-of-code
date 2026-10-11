# Day 5: Binary Boarding

## Approach

**Part 1:** Parse boarding passes using binary search simulation. Each pass specifies a seat via 10 characters: 7 for row (F=lower, B=upper) and 3 for column (L=lower, R=upper). Convert to binary (F/L=0, B/R=1) to get row and column directly. Seat ID = row × 8 + column. Find the maximum.

**Part 2:** Collect all seat IDs, sort them, and find the gap where two IDs differ by 2 (missing seat between occupied neighbors).

## Complexity

- Time: O(n × 10) parsing + O(n log n) sorting = **O(n log n)**
- Space: O(n) for seat list

## Commands

```bash
cd 2020/rust

# Test
cargo test --bin day5

# Run with real input (the official example has no free seat, so part 2 errors on it)
cargo run -q --bin day5 < input/2020/day5.txt
```
