# Day 18: Boiling Boulders

## Approach

**Part 1:** Count exposed faces of lava droplets.
- Parse coordinates into a set.
- For each cube, check 6 neighbors; count faces not adjacent to another cube.

**Part 2:** Count only external surface area.
- Expand bounds by 1 to start outside all cubes.
- Flood fill from outside using BFS to mark all external space.
- Count cube faces adjacent to external space only.

## Complexity

- **Time:** O(n + k) where n is cube count, k is external space volume.
- **Space:** O(n + k) for hash maps and queue.

## Run

```bash
zig run src/day_18.zig < input.txt
```

## Test

```bash
zig test src/day_18.zig
zig run src/day_18.zig < src/examples/day_18.txt
```
