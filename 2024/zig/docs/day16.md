# Day 16: Reindeer Maze (Zig)

The reindeer starts on `S` facing east and has to reach `E`: a step forward scores 1,
a 90 degree turn 1000. Part 1 asks for the lowest score, part 2 for the number
of tiles that lie on at least one lowest-scoring path.

## Approach

A search state is a tile plus a facing, which turns both turning and stepping
into ordinary weighted edges, so Dijkstra over `tiles * 4` states answers part 1:
the cheapest state at `E` over the four facings.

Part 2 runs the same search a second time on the reverse graph (step against
the facing; turns are symmetric) from every facing at `E`. A state lies on a
best path exactly when its score from the start plus its score to the end
equals the best score, so part 2 counts the tiles with at least one such state.

## Complexity

O(tiles * log tiles) time and O(tiles) space - four states per tile, each
relaxing at most three edges.

## Run / test

```bash
cd 2024/zig
zig test day16.zig
zig run day16.zig < examples/day16.txt
```
