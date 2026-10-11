# Day 14: Restroom Redoubt (Zig)

## Approach

**Part 1:** Each robot moves in a straight line and wraps around the grid
edges. Compute each robot's position after 100 seconds directly via modular
arithmetic (`(px + vx*t) mod width`, `(py + vy*t) mod height`), bucket robots
into the four quadrants (excluding the middle row/column), and multiply the
quadrant counts.

**Part 2:** The robots form a recognizable picture at the first second when
no two of them occupy the same tile. Simulate seconds from 0 upward, checking
for duplicate positions, until a collision-free arrangement is found.
Positions repeat with period `width*height`, which bounds the search.

The grid size (101x103 for the real puzzle, 11x7 for the example) is not
given in the input, so the runner uses 11x7 when every robot fits the example
grid and the official 101x103 otherwise.

## Complexity

- Part 1: O(n) for n robots.
- Part 2: O(n * width * height) worst case.

## Run / test

```bash
cd 2024/zig
zig test day14.zig
zig run day14.zig < examples/day14.txt
```
