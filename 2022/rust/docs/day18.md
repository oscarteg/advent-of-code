# Day 18: Boiling Boulders

## Approach

**Part 1:** Calculate the surface area of a 3D lava droplet by counting exposed faces.
- Parse coordinates into a set of cubes.
- For each cube, check its 6 neighbors; count faces not adjacent to another cube.
- Total: 6 × cubes − 2 × adjacent pairs.

**Part 2:** Find only the external surface area, excluding internal air pockets.
- Expand the bounding box by 1 to ensure we start outside.
- Flood fill from outside using BFS to mark all external space.
- Count only cube faces adjacent to external space (not to internal air).

## Complexity

- **Time:** O(n + k) where n is cube count, k is external space volume.
- **Space:** O(n + k) for set and BFS queue.

## Run

```bash
cargo run -q --bin day18 < input/day18.txt
```

## Test

```bash
cargo test --bin day18
```
