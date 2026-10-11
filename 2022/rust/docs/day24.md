# Day 24: Blizzard Basin

## Approach

Parse the grid into dimensions, start/end column, and a list of blizzards
(position + direction). A blizzard's position at any time `t` is computed
directly with modular arithmetic (`rem_euclid`), so the full occupancy map
for time `t` is derived rather than simulated step by step.

Since each blizzard's position is periodic in `width` (horizontal) or
`height` (vertical), the whole occupancy pattern repeats with period
`lcm(width, height)`. Occupancy sets for every time in one period are
precomputed once.

**Part 1**: BFS over `(position, time mod period)` states from the start
opening to the end opening, trying "stay", "up", "down", "left", "right"
each minute and skipping cells occupied by a blizzard or outside the grid
(except the two openings).

**Part 2**: run the same BFS three times - start to end, end back to start,
start to end again - each starting from the elapsed time of the previous
leg. Part 1's answer is the first leg; part 2's answer is the total of all
three legs.

## Complexity

Precomputing occupancy: `O(period * blizzards)`. Each BFS visits at most
`O(width * height * period)` states, each with 5 neighbors.

## Commands

```bash
cargo test --bin day24
cargo run -q --bin day24 < input/day24-example.txt
```
