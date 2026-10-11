# Day 12: Christmas Tree Farm

## Problem

Presents are given as 3x3 patterns; each region is `WxH` plus how many of
each present must fit into it. Presents rotate and flip freely, but never
overlap. Count the regions that can hold all their presents. This day has a
single part.

---

*Solution spoilers below.*

## Approach

Settle each region in three steps:

1. Fail when the presents' total area exceeds `W*H`.
2. Succeed when they fit in separate 3x3 blocks: `presents <= (W/3) * (H/3)`.
3. Otherwise search exactly: scanning the grid in row-major order, the first
   undecided cell either stays empty, spending one cell of the waste budget
   (`W*H - area`), or anchors a present that must start there, so every
   packing is reached once. Parsing is `O(n)`; the search is exponential in
   the worst case, bounded by that budget.

## Run / test

```sh
cd 2025
zig test src/day12.zig
zig run src/day12.zig < src/examples/day12.txt
```

The example prints `Part 1: 2`.
