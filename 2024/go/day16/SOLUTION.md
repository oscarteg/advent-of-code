# Day 16: Reindeer Maze

The reindeer starts on `S` facing east and has to reach `E`: a step forward scores
1, a 90 degree turn 1000.

## Part 1

A search state is a tile plus a facing, which turns both turning and stepping
into ordinary weighted edges. Dijkstra over `tiles * 4` states gives the answer as
the cheapest state at `E` over the four facings. Complexity:
O(tiles * log tiles).

## Part 2

The same search runs again on the reverse graph (step against the facing; turns
are symmetric) from every facing at `E`. A state lies on a best path exactly
when its score from the start plus its score to the end equals the best score,
so the answer counts the tiles with at least one such state. Complexity:
O(tiles * log tiles).

## Run / test

```bash
go test -race -count=1 ./day16
go run ./day16 < day16/testdata/example.txt
```
