# Day 24: Lobby Layout

## Approach

**Part 1:** Parse direction sequences and navigate a hexagonal grid using axial coordinates. Flip each tile; count black tiles (flipped odd times).

**Part 2:** Simulate 100 days of tile flipping using cellular automaton rules:
- Black tiles with 0 or > 2 black neighbors flip white
- White tiles with exactly 2 black neighbors flip black

Uses a HashMap to track black tiles. Each step checks only black tiles and their neighbors to avoid iterating the infinite grid.

## Complexity

- Part 1: O(n·m) where n = number of paths, m = average path length
- Part 2: O(100·k) where k = number of tiles to check per step

## Commands

```bash
cd 2020/rust
cargo test --bin day24
cargo run -q --bin day24 < input/2020/day24-example.txt
```
