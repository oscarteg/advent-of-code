# Day 4: Ceres Search (C++)

## Part 1
From every `X` in the grid, check all 8 directions for the literal sequence `XMAS`,
bounds-checked cell by cell. Count the matches.
Complexity: O(rows * cols * 8 * 4) time, O(rows * cols) space for the parsed grid.

## Part 2
From every `A` in the grid, check both diagonals through it: each must read `MAS` or `SAM`
forwards or backwards (i.e. the two endpoints are `M`/`S` in either order). Count positions
where both diagonals match, forming an X-shape "X-MAS".
Complexity: O(rows * cols) time, O(rows * cols) space.

## Run / test

```bash
cd 2024/cpp
make test day=day4
make -s run day=day4 < examples/day4.txt
```
