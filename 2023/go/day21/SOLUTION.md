# Day 21: Step Counter

## Part 1
Breadth-first flood from `S` over the map as handed out. The elf can step off a
plot and back onto it, so every plot whose distance has the same parity as the
step budget is reachable. **Complexity:** O(width * height).

## Part 2
The same flood on the infinitely tiled map reproduces the example's counts, but
26501365 steps are far too many to flood. The real map is square with an open
border and an open row and column through `S`, so the flood swallows one more
map copy in every direction per `width` steps: the counts after
`offset + tiles*width` steps are a quadratic in `tiles`. Flood the first three
tile-aligned counts, fit the quadratic, evaluate it. **Complexity:** O(width^3)
for those three floods, then constant time.

`solve` uses the real parameters (64 and 26501365 steps), so on the example it
prints those applied to the example. `solveWith` takes explicit step counts and
the tests use the example's own: 6, 10, 50, 100 and 500 steps give 16, 50, 1594,
6536 and 167004 plots.

## Run and test
```bash
go run ./day21 < input.txt
go test -race ./day21
```
