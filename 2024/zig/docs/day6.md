# Day 6: Guard Gallivant (Zig)

A guard patrols a lab grid: step forward, or turn right 90 degrees when an
obstacle (`#`) blocks the way. The guard's start is marked `^`, `>`, `v`, or
`<` for its facing direction.

## Part 1: Distinct Positions

Simulate the patrol from the start until the guard steps off the grid,
recording every distinct cell visited (including the start). The answer is
the count of those cells.

## Part 2: Loop-Causing Obstructions

Count how many single new obstructions (placed on an empty cell other than
the start) would trap the guard in an infinite loop instead of letting it
leave. Only cells on the part 1 path can matter, since the guard never
reaches a cell it wouldn't otherwise visit. For each candidate, re-walk the
grid with that extra wall; the walk loops if it revisits the same
(position, facing) state before leaving the grid.

## Complexity

- **Time:** Part 1 is `O(rows * cols)`. Part 2 re-walks the grid for each of
  up to `O(rows * cols)` candidates, and each walk is bounded by its
  `4 * rows * cols` distinct (position, facing) states, giving
  `O((rows * cols)^2)` overall.
- **Space:** `O(rows * cols)` for the grid and the visited/seen-state sets.

## Run / test

```bash
cd 2024/zig
zig test day6.zig
zig run day6.zig < examples/day6.txt
```
