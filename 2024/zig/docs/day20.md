# Day 20: Race Condition (Zig)

The racetrack is a single path from `S` to `E` through a walled maze. Once per race a program
may ignore walls for up to 2 picoseconds (part 1) or 20 (part 2), and both parts count the
cheats saving at least a threshold - 100 for the real input, a parameter here so the published
example counts stay testable.

## Approach

BFS once from `S` and once from `E` over the parsed wall bitmap, so every track cell knows its
distance from the start and to the finish; `fromStart[E]` is the uncheated time. Because a
cheat moves freely through walls, landing on a track cell at Manhattan distance `d` costs
exactly `d` picoseconds and the cheated race takes `fromStart[begin] + d + toEnd[landing]`.
Both parts count the cell pairs whose saving reaches the threshold, differing only in the
largest allowed `d`.

## Complexity

`O(T * L^2)` time for `T` track cells and cheat length `L`, since each cell scans the diamond
of cells within `L`, and `O(W * H)` space.

## Run / test

```bash
cd 2024/zig
zig test day20.zig
zig run day20.zig < examples/day20.txt
```
