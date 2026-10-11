# Day 12: Hot Springs

- **Part 1:** Parse every row into spring states (`.`, `#`, `?`) and group
  sizes. Walk the row with a memoised recursion over the pair
  *(spring position, next group)*: an unknown or operational spring may be
  skipped, and a group can start here when the next `size` springs are all
  non-operational and the spring behind them is not damaged. The row is counted
  when both the springs and the groups run out together.
- **Part 2:** The same count after unfolding each row five times, joining the
  copies with `?` and repeating the group list five times. Memoisation keeps it
  linear instead of exponential.
- Each row costs `O(springs * groups)` time and space, so the whole input is
  `O(rows * springs * groups)`.

From `2023/rust`:

```sh
cargo run -q --bin day12 < examples/day12.txt # 21
cargo run -q --bin day12 < examples/day12.txt # 525152
cargo test day12::
```

Use the same commands with your puzzle input redirected into stdin.
