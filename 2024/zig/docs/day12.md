# Day 12: Garden Groups (Zig)

The map is a grid of plant letters; orthogonally touching plots of the same
letter form a region. Part 1 prices a region as area * perimeter, part 2 as
area * number of straight sides.

## Approach

Parse the map into a fixed-width grid of letters, then flood-fill from every
plot not yet visited to collect one region at a time. Inside the fill, a step
that leaves the region (different letter or grid edge) adds one to its
perimeter. A region has as many sides as corners, so the same fill counts
corners per plot: for each of the four diagonals, both adjacent plots outside
the region is a convex corner, and both inside while the diagonal plot is
outside is a concave corner.

## Complexity

O(width * height) time and space - every plot is visited once and does a
constant amount of perimeter and corner work.

## Run / test

```bash
cd 2024/zig
zig test day12.zig
zig run day12.zig < examples/day12.txt
```
