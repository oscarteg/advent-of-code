# Day 6: Guard Gallivant

Simulate the guard: step forward, or turn right 90 degrees when blocked.

## Part 1

Walk from the start until the guard leaves the grid, recording every distinct
cell visited. The answer is the count of distinct cells.

## Part 2

For every cell on the part 1 path (except the start), re-walk the grid with an
extra obstruction there. A walk loops if it revisits the same
(position, facing) state; count the candidates that loop.

## Complexity

Part 1 is `O(rows * cols)`. Part 2 tries `O(rows * cols)` candidates, each
walk bounded by the `4 * rows * cols` distinct states, so `O((rows * cols)^2)`
overall.

## Commands

```bash
go test -race ./day06
go vet ./day06
go run ./day06 < day06/testdata/example.txt
```

Prints:

```text
Part 1: 41
Part 2: 6
```
