# Day 9: Smoke Basin

Parse the grid of single-digit heights. A low point is lower than every
(up to four) orthogonal neighbor; part 1 sums `height + 1` over all low
points.

A basin is a maximal region of non-9 cells; it contains exactly one low
point, so flood-filling from each low point over cells `< 9` finds every
basin. Part 2 multiplies the sizes of the three largest basins.

For a grid with `N` cells, the low-point scan is `O(N)` and each
flood fill visits at most `N` cells, so total time is `O(N)`; space is
`O(N)` for the visited set.

From `2021/rust`:

```sh
cargo run --bin day9 < input/day9-example.txt
# Part 1: 15
# Part 2: 1134
cargo test --bin day9
```
