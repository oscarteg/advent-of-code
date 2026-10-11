# Day 24: Blizzard Basin

## Approach

Parse the grid into dimensions, start/end column, and a list of blizzards
(position + direction). A blizzard's position at time `t` is computed
directly with modular arithmetic, so occupancy at any time is derived
rather than simulated minute by minute.

Since each blizzard repeats with period `width` or `height`, the whole
occupancy pattern repeats with period `lcm(width, height)`. Occupancy sets
for every time in one period are precomputed once.

**Part 1**: BFS over `(position, time mod period)` states from the start to
the end opening, trying "stay" and the four directions each minute,
skipping cells occupied by a blizzard or outside the grid (except the two
openings).

**Part 2**: run the same BFS three times - start to end, end to start,
start to end again - each starting at the elapsed time of the previous leg.
Part 1 is the first leg; part 2 is the sum of all three legs.

## Complexity

Precomputing occupancy: `O(period * blizzards)`. Each BFS visits at most
`O(width * height * period)` states, each with 5 neighbors.

## Commands

```bash
go test -race -count=1 ./day24
go run ./day24 < day24/testdata/example.txt
```
