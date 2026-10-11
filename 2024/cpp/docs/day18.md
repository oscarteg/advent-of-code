# Day 18: RAM Run (C++)

Bytes fall one at a time onto a grid at given `x,y` coordinates, corrupting
that cell. Part 1 asks for the shortest path from the top-left to the
bottom-right corner after a fixed number of bytes have fallen; part 2 asks
for the coordinate of the first byte whose fall makes the exit unreachable.

## Approach

Parse every coordinate into a typed list. Part 1 runs a BFS over the grid,
treating fallen-byte cells as walls, and returns the distance to the exit.
Part 2 binary-searches over the byte index: the exit is reachable for a
prefix of byte counts and unreachable for every count after the first
blocker, so the smallest blocking index can be found in O(log n) BFS runs
instead of one BFS per byte.

## Complexity

Grid is `size x size` (7x7 for the example, 71x71 for real input). Each BFS
is O(size^2). Part 1: one BFS. Part 2: O(log n) BFS runs for `n` bytes.
Space: O(size^2) per BFS.

## Run / test

```bash
cd 2024/cpp
make test day=day18
make -s run day=day18 < examples/day18.txt
```
