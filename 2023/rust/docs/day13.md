# Day 13: Point of Incidence

- **Part 1:** For each pattern, test every possible horizontal and vertical
  mirror line. A line is valid when the rows (or columns) reflected across it
  match exactly on the overlapping span. Sum `100 * rows above` for horizontal
  mirrors and `columns to the left` for vertical ones.
- **Part 2:** Same search, but a line is valid when the reflected cells differ
  by exactly one ("smudge") across the whole pattern instead of zero.
- Runtime is `O(patterns * cells)`: each candidate line compares at most the
  full grid once, and there are at most `rows + columns` candidate lines.

From `2023/rust`:

```sh
cargo run -q --bin day13 < examples/day13.txt # 405
cargo run -q --bin day13 < examples/day13.txt # 400
cargo test day13::
```

Use the same commands with your puzzle input redirected into stdin.
