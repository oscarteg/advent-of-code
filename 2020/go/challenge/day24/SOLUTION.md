# Day 24: Lobby Layout

## Approach

**Part A:** Parse hex direction sequences and navigate a hexagonal grid using axial coordinates. Track tiles by their final position; count those flipped an odd number of times (black).

**Part B:** Simulate 100 days of tile flipping using cellular automaton rules:
- Black tiles with != 1 or 2 black neighbors flip white
- White tiles with exactly 2 black neighbors flip black

Each simulation step only processes black tiles and their neighbors to avoid checking an infinite grid.

## Complexity

- Part A: O(n·m) where n = number of paths, m = average path length
- Part B: O(100·k) where k = number of tiles to check per step

## Commands

```bash
cd 2020/go
go test -race ./challenge/day24
go run . 24 < ../rust/input/2020/day24-example.txt
```
