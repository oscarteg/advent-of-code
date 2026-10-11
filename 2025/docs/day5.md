# Day 5: Cafeteria

## Problem

The input lists inclusive fresh ID ranges `start-end`, then the available
ingredient IDs.

**Part 1:** count how many available IDs fall within at least one fresh range.

**Part 2:** count the total number of distinct IDs covered by all ranges combined.

---

*Solution spoilers below.*

## Approach

**Part 1:** for each ID, check if it's contained in any of the inclusive ranges.

**Part 2:** sort ranges by start position, then merge adjacent or overlapping
ranges by extending the end point, finally sum the coverage of each merged range.

## Complexity

- Part 1: O(n·m) where n = number of IDs, m = number of ranges
- Part 2: O(m log m) for sorting + O(m) for merging

## Commands

```bash
cd 2025
zig fmt --check src/day5.zig
zig test src/day5.zig
zig run src/day5.zig < src/examples/day5.txt
```
