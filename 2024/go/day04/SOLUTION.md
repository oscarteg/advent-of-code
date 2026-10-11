# Day 4: Ceres Search

## Part 1
From every `X` in the grid, check all 8 directions for the literal sequence `XMAS`,
bounds-checked cell by cell. Count the matches.
Complexity: O(rows * cols) time and space.

## Part 2
From every `A`, check both diagonals through it: each must read `MAS` or `SAM`
(the two endpoints match `M`/`S` in either order). Count positions where both diagonals
match, forming an X-shape "X-MAS".
Complexity: O(rows * cols) time and space.

## Run / test
```bash
cd 2024/go
go run ./day04 < day04/testdata/example.txt
go test ./day04
```
