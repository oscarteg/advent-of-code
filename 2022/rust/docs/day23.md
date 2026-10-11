# Day 23: Unstable Diffusion

## Problem

Simulate elf movement on a 2D grid. Each round, elves propose moves based on checking neighbors in a rotating sequence of directions: North, South, West, East. An elf only moves if it has neighbors and can find a clear direction. If multiple elves propose the same destination, none move.

**Part 1:** Count empty ground tiles in the bounding rectangle after 10 rounds.
**Part 2:** Find the first round where no elf moves.

## Approach

- **Data Structure:** Use a `HashSet<(i32, i32)>` for elf positions (sparse representation).
- **Movement Logic:** For each elf with neighbors, check directions in rotating order. Return the first direction with no blocking elves in its 3-tile span.
- **Conflict Resolution:** Collect proposals in a HashMap, move only elves with unique destinations.
- **Round Tracking:** Rotate direction index after each round; continue until `round()` returns false (no movement).

## Complexity

- **Time:** O(r * e * 1) where r = rounds, e = number of elves; checking 8 neighbors and proposing moves is O(1) per elf.
- **Space:** O(e) for the HashSet and proposal map.

## Commands

```bash
cd 2022/rust
cargo test --bin day23
cargo run -q --bin day23 < input/day23-example.txt
```

## Example Output

```
Part 1: 110
Part 2: 20
```
