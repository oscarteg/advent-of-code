# Day 22: Sand Slabs

- **Settling:** parse every line into an axis-aligned brick with inclusive
  bounds, sort by lowest `z` and drop the bricks in that order. A height map
  keyed by `(x, y)` column stores the highest occupied `z` and the brick that
  owns it, so a falling brick lands one cube above the tallest column under its
  footprint; the bricks owning those top cubes are its supporters. That gives
  `supports` and `supported_by` for the settled stack.
- **Part 1:** a brick is safe when every brick resting on it has at least one
  other supporter.
- **Part 2:** for each brick, walk up the support graph with a queue and mark a
  brick as falling once all of its supporters are falling; the marks minus the
  removed brick are its chain reaction. Summing those gives the answer.
- Settling is `O(n log n)` plus the total footprint area; part 2 traverses the
  graph once per brick, `O(n * edges)`.

From `2023/rust`:

```sh
cargo run -q --bin day22 < examples/day22.txt # 5
cargo run -q --bin day22 < examples/day22.txt # 7
cargo test day22::
```

Redirect your own puzzle input into stdin the same way.
