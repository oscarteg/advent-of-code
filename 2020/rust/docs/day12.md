# Day 12: Rain Risk

## Part 1

The ship tracks a position `(x, y)` and a heading vector, starting facing
east. `N`/`S`/`E`/`W` move the ship directly; `L`/`R` rotate the heading
vector by a multiple of 90 degrees; `F` moves the ship by the heading vector.
Rotation is a simple 90-degree vector rotation (`(x, y) -> (-y, x)`) applied
`degrees / 90` times. Answer is the Manhattan distance from the origin.

Complexity: O(n) time, O(1) extra space for n instructions.

## Part 2

Same instruction set, but `N`/`S`/`E`/`W` and `L`/`R` move/rotate a waypoint
relative to the ship instead of the ship itself; `F` moves the ship by the
waypoint's offset, repeated by the given value. Same rotation helper is
reused for the waypoint vector.

Complexity: O(n) time, O(1) extra space.

## Commands

```bash
cargo test --bin day12
cargo run -q --bin day12 < input/2020/day12-example.txt
```
