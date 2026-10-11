# Day 17: Conway Cubes

## Part 1

Simulate Conway's Game of Life in 3D. Starting with a 2D grid on the z=0 plane, simulate 6 generations. A cube becomes active if it has exactly 3 active neighbors (from 26 neighboring positions in 3D). An active cube stays active with 2 or 3 neighbors, otherwise goes inactive.

**Approach:** Use a HashSet to track active cubes. Each generation, collect all cubes to check (active cubes + their neighbors), then apply the rules.

**Complexity:** O(n * m^3) where n is number of active cubes and m is the spread per generation. With 6 generations, spread is bounded.

## Part 2

Same as Part 1, but in 4D (x, y, z, w). Each cube has 80 neighbors in 4D (3^4 - 1).

**Approach:** Identical algorithm, but with 4D coordinates and 80 neighbors instead of 26.

**Complexity:** O(n * m^4) per generation in 4D space.

## Run

```bash
cargo run -q --bin day17 < input/2020/day17-example.txt
```

## Test

```bash
cargo test --bin day17
```
