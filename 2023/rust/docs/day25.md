# Day 25: Snowverload

- **Part 1:** The wires form an undirected graph whose global minimum cut is the
  three wires to pull. Store every wire as two unit-capacity arcs, fix one
  component as the source and run Edmonds-Karp towards every other component.
  The cheapest of those maximum flows is the global minimum cut, because the
  fixed component lies on one of the cut's two sides. The components still
  reachable from the source in the residual graph form one group, the rest the
  other; multiply the two sizes.
- Each candidate sink is abandoned as soon as its flow can no longer beat the
  cheapest cut found so far, so a candidate costs `O(cut * wires)` with
  `cut == 3` for puzzle input. In total `O(components * cut * wires)`.
- Single part puzzle: there is no part 2 beyond pushing the button.

From `2023/rust`:

```sh
cargo run -q --bin day25 < examples/day25.txt # 54
cargo test day25::
```

Use the same command with your puzzle input redirected into stdin.
