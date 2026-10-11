# Day 14: Regolith Reservoir

Rock paths are parsed into a `AutoHashMap(Point, void)` of occupied cells by
expanding each axis-aligned segment between consecutive points. Sand units
drop one at a time from `(500, 0)`, each falling down, then down-left, then
down-right while that cell is free, settling once none of the three are.

## Part 1

Simulation stops the moment a unit falls past the lowest rock row
(`y > max_y`) with no floor to catch it; the answer is the count of units
settled before that.

**Complexity:** `O(units * max_y)` time (each unit falls at most `max_y`
rows, settled lookups are `O(1)` hashed), `O(rocks + units)` space.

## Part 2

A floor two rows below the lowest rock is treated as an infinite row (never
materialized) and the simulation runs until a unit settles on the source
itself, blocking it.

**Complexity:** same bound as part 1, `O(units * max_y)`.

## Run / Test

```bash
zig run src/day_14.zig < src/examples/day_14.txt
zig test src/day_14.zig
```
