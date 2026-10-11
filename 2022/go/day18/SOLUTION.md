# Day 18: Boiling Boulders

## Approach

**Part 1:** Calculate surface area of lava droplets.
- Parse coordinates into a map set.
- For each cube, count its 6 faces not adjacent to other cubes.

**Part 2:** Calculate external surface area only.
- Expand bounding box by 1 to ensure we start outside.
- BFS flood fill from outside to mark all external space.
- Count cube faces adjacent to external space.

## Complexity

- **Time:** O(n + k) where n is cube count, k is external space volume.
- **Space:** O(n + k) for map and queue.

## Run

```bash
go run ./day18 < input.txt
```

## Test

```bash
go test ./day18
```
