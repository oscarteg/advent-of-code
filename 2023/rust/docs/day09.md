# Day 9: Mirage Maintenance

- **Part 1:** For each history, recursively take successive differences until
  all values are zero, then sum the last value of every row on the way back up
  to extrapolate the next reading. Sum that over all histories.
- **Part 2:** Same recursion, but extrapolate backward: each row's leading
  value minus the extrapolated previous value of its differences.
- Each history of length `n` produces `O(n)` recursive calls with `O(n)` work
  per call (building the differences), so part 1 and part 2 are both
  `O(n^2)` per history.

From `2023/rust`:

```sh
cargo run -q --bin day09 < examples/day09.txt # 114
cargo run -q --bin day09 < examples/day09.txt # 2
cargo test day09::
```

Use the same commands with your puzzle input redirected into stdin.
