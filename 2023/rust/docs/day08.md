# Day 8: Haunted Wasteland

- **Part 1:** Parse the instruction line and the `name = (left, right)` nodes
  into a map, then walk from `AAA`, repeating the instructions, until `ZZZ` is
  reached.
- **Part 2:** Walk from every node ending in `A` to the first node ending in
  `Z`. A ghost revisits its destination every walk length, so the answer is the
  least common multiple of those lengths.
- A walk costs `O(steps)` time and the network `O(nodes)` memory. Walking longer
  than `nodes * instructions` steps means the (node, instruction) state
  repeated, so the destination is unreachable and an error is returned instead
  of looping forever.

From `2023/rust`:

```sh
cargo run -q --bin day08 < examples/day08.txt  # 2
cargo run -q --bin day08 < examples/day08-2.txt # 6
cargo test day08::
```
