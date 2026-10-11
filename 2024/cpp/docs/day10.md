# Day 10: Hoof It (C++)

## Part 1: Trail Scores

A hiking trail connects elevation 0 to 9, increasing by exactly 1 at each step. For each trailhead (elevation 0), count the number of distinct nines reachable via valid hiking trails. A "trail score" is the count of distinct peaks reachable from a trailhead.

Sum all trail scores across all trailheads.

## Part 2: Trail Ratings

A "trail rating" is the number of distinct paths from a trailhead to any peak (9). Count all possible paths, not just unique destinations.

Sum all trail ratings across all trailheads.

## Approach

1. Parse the topographic grid (each cell is a single digit, elevation 0-9).
2. For each cell with elevation 0 (trailhead):
   - Use DFS to explore all valid next steps (adjacent cells with elevation exactly 1 higher).
   - Track unique destination peaks (for part 1).
   - Count total path completions (for part 2).
3. Aggregate results across all trailheads.

## Complexity

- **Time:** O(rows × cols × avg_path_length). In practice, the grid is small and paths are short.
- **Space:** O(rows × cols) for the grid, O(path_depth) for recursion stack.

## Run / test

```bash
cd 2024/cpp
make test day=day10
make -s run day=day10 < examples/day10.txt
```
