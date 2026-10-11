# Day 17: Clumsy Crucible

- **Part 1 & 2:** Dijkstra over the state `(x, y, last direction, run length)` instead of plain
  `(x, y)`, since the heat loss to reach a cell depends on how many consecutive blocks were
  already travelled in the current direction. From a state, the crucible may continue the same
  direction while `run < max_run`, or turn 90 degrees once `run >= min_run`; it never reverses.
  The goal is reached only once `run >= min_run`. Part 1 uses `min_run = 0, max_run = 3`; part 2
  (the "ultra crucible") uses `min_run = 4, max_run = 10`.
- **Complexity:** `O(width * height * 4 * max_run * log(...))` for the heap-based search.

From `2023/rust`:

```sh
cargo run -q --bin day17 < examples/day17.txt  # 102
cargo run -q --bin day17 < examples/day17.txt  # 94
cargo run -q --bin day17 < examples/day17-2.txt # 71
cargo test day17::
```

Use the same commands with your puzzle input redirected into stdin.
