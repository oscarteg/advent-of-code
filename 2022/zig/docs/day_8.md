# Day 8: Treetop Tree House

## Part 1: Visible Trees

A tree is visible from the outside if all trees between it and at least one edge are shorter.
We check all four directions (left, right, up, down) for each tree—if any direction is completely clear, the tree is visible.

**Complexity:** O(n²m²) where n×m is grid size (checking all trees from each position).

## Part 2: Scenic Score

The scenic score is the product of visibility distances in all four directions.
For each direction, we count trees until we hit one that's taller or reach the edge.

**Complexity:** O(n²m²) with early termination on tall trees.

## Running

```bash
zig run src/day_8.zig < input.txt
zig test src/day_8.zig
```
