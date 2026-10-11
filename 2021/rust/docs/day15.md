# Day 15: Chiton

Parse the grid of risk digits, then run Dijkstra's algorithm from the top-left
to the bottom-right cell, summing entered-cell risk along the way.

Part 1 runs Dijkstra on the grid as given. Part 2 tiles the grid 5x5: risk in
tile `(tx, ty)` is the base risk plus `tx + ty`, wrapping back to 1 after 9.
The tiled risk is computed lazily per cell instead of materializing the full
grid.

With `N` cells and 4-directional moves, each Dijkstra run is
`O(N log N)` using a binary heap, and space is `O(N)`.

From `2021/rust`:

```sh
cargo test --bin day15
cargo run -q --bin day15 < input/day15-example.txt
```

The example prints exactly:

```text
Part 1: 40
Part 2: 315
```
