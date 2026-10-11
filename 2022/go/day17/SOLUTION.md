# Day 17: Pyroclastic Flow

## Part 1
Simulate rocks falling into a 7-wide chamber: each rock alternates being pushed by the next
jet (wrapping) and falling one row, settling when a fall would collide with the floor or a
previously settled rock. The tower is a slice of row bitmasks. Height after 2022 rocks is the
answer. Complexity: O(rocks x fall steps).

## Part 2
Simulating 1e12 rocks directly is infeasible, so we detect a cycle: after each rock settles we
key the state by (rock shape index, jet index, top 50 rows of the tower). Once a key repeats,
the gap between occurrences gives a repeating (rock count, height) cycle. We then compute how
many full cycles fit in the remaining rocks, add their height, and simulate the small remainder
from already-recorded heights. Complexity: O(rocks until a cycle is found), independent of the
target rock count.

## Run / test
```bash
go test -race -count=1 ./day17
go run ./day17 < day17/testdata/example.txt
```
