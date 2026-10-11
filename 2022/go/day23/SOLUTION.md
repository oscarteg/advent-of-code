# Day 23: Unstable Diffusion

## Problem

Simulate elf movement on a 2D grid. Each round, elves propose moves based on checking neighbors in a rotating sequence of directions: North, South, West, East. An elf only moves if it has neighbors and can find a clear direction. If multiple elves propose the same destination, none move.

**Part 1:** Count empty ground tiles in the bounding rectangle after 10 rounds.
**Part 2:** Find the first round where no elf moves.

## Approach

- **Data Structure:** Use `map[Pos]bool` for elf positions (sparse representation with O(1) lookups).
- **Movement Logic:** For each elf with neighbors, check directions in rotating order. Return the first direction with no blocking elves in its 3-tile span.
- **Conflict Resolution:** Collect proposals in a map, move only elves with unique destinations.
- **Round Tracking:** Rotate direction index after each round; track first round with no movement.

## Complexity

- **Time:** O(r * e) where r = rounds, e = number of elves; map lookups for neighbor checking.
- **Space:** O(e) for the elf map and proposal map.

## Commands

```bash
cd 2022/go
go test ./day23
go run ./day23 < day23/testdata/example.txt
```

## Example Output

```
Part 1: 110
Part 2: 20
```
