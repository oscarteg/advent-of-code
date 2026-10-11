# Day 13: Distress Signal

## Approach

**Part 1:** Parse pairs of packets and compare using custom ordering:
- Integers compare numerically
- Lists compare element-by-element; shorter is "less"
- Int vs list: wrap int in a list and compare

Sum the 1-indexed pair numbers in correct order.

**Part 2:** Parse all packets (skip blank lines), add dividers `[[2]]` and `[[6]]`, sort, find product of 1-indexed divider positions.

## Complexity

- **Time:** O(n log n) for sorting; O(n * d) for comparisons where d is max depth
- **Space:** O(n) for packet storage

## Run & Test

```bash
cd 2022/zig
zig fmt --check src/day_13.zig
zig test src/day_13.zig
zig run src/day_13.zig < src/examples/day_13.txt
```

Expected output: `Part 1: 13` and `Part 2: 140`
