# Day 15: Beacon Exclusion Zone

Every reading gives a sensor and its closest beacon; their Manhattan distance
is the radius of a diamond that can hold no other beacon.

## Part 1
On a single row a diamond covers `[sx - spread, sx + spread]` with
`spread = radius - |sy - row|`, and nothing when that is negative. Sort those
ranges, merge overlapping and touching ones, sum their lengths and subtract
the distinct known beacons on the row. Complexity: O(n log n) for n readings.

## Part 2
Scan the rows of the search square and merge the ranges again: walking them
from x = 0 leaves the first x no range covers. The first row with such a gap
holds the distress beacon, whose tuning frequency is `x * 4000000 + y`.
Complexity: O(limit * n log n), stopping at the first gap.
`main` uses the real parameters (row 2000000, square 0...4000000); the tests
use the example's row 10 and square 0...20 through `solveWith`.

## Run / test
```bash
cd 2022/zig
zig test src/day_15.zig
zig run src/day_15.zig < src/examples/day_15.txt
```
