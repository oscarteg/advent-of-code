# Day 10: Hoof It

The input is a grid of heights 0-9. A trail starts at a 0 and climbs exactly one height
per step (up, down, left or right) until it reaches a 9.

- **Part 1 (score):** for every trailhead, count the distinct 9s it can reach; sum them.
- **Part 2 (rating):** for every trailhead, count the distinct trails; sum them.

One depth-first climb per trailhead answers both: it collects the reachable summits in a set
(score) and returns the number of trails that end at them (rating).

## Complexity

Each trailhead explores at most every uphill path; heights strictly increase, so there are
no cycles. O(trailheads x paths) time, O(grid) memory.

## Run / test

```bash
cd 2024/go
go test -race ./day10
go run ./day10 < day10/testdata/example.txt
```
