# Day 10: Pipe Maze

- **Parsing:** each character becomes a `Tile`: a pipe stores the two
  `Direction`s it connects, `.` is ground, `S` is the start.
- **Finding the loop:** check which neighbors of `S` have a pipe that
  connects back toward it (exactly two must); walk the loop one step at a
  time, always exiting through the connection that isn't where we came from,
  until we return to `S`. This also replaces `S` conceptually with its real
  shape, since we only ever use its inferred directions.
- **Part 1:** the farthest tile is half the loop's length (`loop.len() / 2`),
  since it is a simple cycle on a grid.
- **Part 2:** apply the shoelace formula to the loop's vertices to get its
  area, then Pick's theorem (`area = interior + boundary/2 - 1`) to recover
  the enclosed tile count directly, without scanning the whole grid.
- Both parts are `O(rows * cols)`: one pass to parse, one to walk the loop
  (bounded by the grid size) and a linear scan over the loop for the area.

From `2023/rust`:

```sh
cargo run -q --bin day10 < examples/day10.txt   # 4
cargo run -q --bin day10 < examples/day10-2.txt  # 8
cargo run -q --bin day10 < examples/day10-3.txt  # 4
cargo run -q --bin day10 < examples/day10-4.txt  # 8
cargo test day10::
```
