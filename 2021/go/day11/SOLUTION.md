# Day 11: Dumbo Octopus

Parse the 10x10 grid of energy levels. Each step increments every cell, then
any cell above 9 flashes: it resets to 0 and increments its (up to 8)
neighbors via a work stack, chaining further flashes without double-counting
a cell within the same step.

Part 1 sums flashes over 100 steps. Part 2 keeps stepping past 100 until a
step flashes every one of the 100 cells simultaneously, returning that step
number.

Each step visits every cell and each flash a constant number of neighbors, so
runtime is `O(steps * cells)` and space is `O(cells)`.

From `2021/go`:

```sh
go test -race ./day11
go vet ./day11
go run ./day11 < day11/testdata/example.txt
```

The example prints exactly:

```text
Part 1: 1656
Part 2: 195
```
