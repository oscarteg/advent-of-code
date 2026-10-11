# Day 18: RAM Run

Bytes fall onto a grid at given coordinates, corrupting cells.

## Part 1

BFS from the top-left to the bottom-right corner over the grid after a fixed
number of bytes have fallen, treating fallen-byte cells as walls. Complexity:
O(size^2).

## Part 2

Binary-search the byte index: the exit is reachable for a prefix of byte
counts and unreachable for every count after the first blocker. Each probe
re-runs the part 1 BFS. Complexity: O(size^2 log n) for n bytes.

## Run / test

```bash
go test -race -count=1 ./day18
go run ./day18 < day18/testdata/example.txt
```
