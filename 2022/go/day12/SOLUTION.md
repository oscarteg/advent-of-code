# Day 12: Hill Climbing Algorithm

## Part 1
Parse the grid, mapping `S` to elevation `a` and `E` to elevation `z`, and
record their positions. Run a breadth-first search over the grid graph where
an edge `u -> v` exists if `elevation(v) <= elevation(u) + 1`. The answer is
the BFS distance from `S` to `E`.

## Part 2
Any lowest point (elevation `a`) is a valid start. Instead of running BFS
from every candidate start, run a single BFS from `E` over the *reversed*
graph (edge `u -> v` exists if `elevation(u) <= elevation(v) + 1`, i.e. the
original edge reversed). The answer is the minimum recorded distance among
all cells at elevation `a`.

## Complexity
BFS visits every cell and edge once: `O(width * height)` time and space.

## Run / test
```bash
go test -race -count=1 ./day12
go run ./day12 < day12/testdata/example.txt
```
