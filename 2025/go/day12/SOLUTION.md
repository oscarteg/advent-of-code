# Day 12: Christmas Tree Farm

Presents are 3x3 patterns; each region is `WxH` plus how many of each present
must fit into it. Presents rotate and flip freely, but never overlap.

## Part 1

Count the regions that hold all their presents, settling each in three steps:

1. Fail when the presents' total area exceeds `W*H`.
2. Succeed when they fit in separate 3x3 blocks: `presents <= (W/3) * (H/3)`.
3. Otherwise search exactly: scanning the grid in row-major order, the first
   undecided cell either stays empty, spending one cell of the waste budget
   (`W*H - area`), or anchors a present that must start there, so every
   packing is reached once. Parsing is `O(n)`; the search is exponential in
   the worst case, bounded by that budget.

## Run / test

```bash
cd 2025/go
go test -race -count=1 ./day12
go run ./day12 < day12/testdata/example.txt
```

The example prints `Part 1: 2`. This day has a single part.
