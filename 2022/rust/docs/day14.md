# Day 14: Regolith Reservoir

Rock paths are parsed into a set of occupied `(x, y)` points, expanding each
segment of a path (segments are always axis-aligned). Sand units are then
dropped one at a time from `(500, 0)`, each falling straight down, then down-left,
then down-right while a space is free; it settles as soon as none of the three
are free.

## Part 1
Sand falls into the abyss once a unit drops below the lowest rock row
(`y > max_y`) with nothing left to catch it; simulation stops there and the
answer is the number of units that settled before that happens.

**Complexity:** each unit can traverse at most `O(max_y)` cells, and settled
sand plus rocks are looked up in a `HashSet`, so total time is
`O(units * max_y)`, space `O(rocks + units)`.

## Part 2
A floor is added two rows below the lowest rock (treated as infinite, not
materialized as points) and the simulation runs until a unit settles on the
source `(500, 0)` itself, which blocks it.

**Complexity:** same shape as part 1, bounded by `O(units * max_y)` since the
floor caps the fall depth.

## Run / Test
```bash
cargo run -q --bin day14 < input/day14-example.txt
cargo test --bin day14
```
