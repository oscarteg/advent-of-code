# Day 20: Race Condition

The racetrack is a single path from `S` to `E` through a maze. Once per race a
program may ignore walls for a few picoseconds; both parts count the cheats
that save at least a threshold, 100 picoseconds for the real input.

## Part 1

One BFS from `S` and one from `E` give every track cell its distance from the
start and to the finish, so `fromStart[E]` is the uncheated time. A 2-picosecond
cheat lands on a track cell at Manhattan distance `d <= 2`, which costs exactly
`d`, so the cheated race takes `fromStart[begin] + d + toEnd[landing]`. Count
the cell pairs whose saving reaches the threshold. Complexity: `O(cells)`.

## Part 2

The same count with cheats up to 20 picoseconds long, so every track cell scans
the diamond of cells within Manhattan distance 20. Complexity: `O(cells * 400)`.

## Run / test

```bash
go test -race -count=1 ./day20
go run ./day20 < day20/testdata/example.txt
```

The example's best cheat saves 76, so at the real threshold it prints zeros. The
tests pin the published counts at lower thresholds.
