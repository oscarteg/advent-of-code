# Day 21: Step Counter

- **Part 1:** Breadth-first flood from `S` over the map as handed out. The elf
  can step off a plot and back onto it, so every plot whose distance has the
  same parity as the step budget is reachable. `O(width * height)`.
- **Part 2:** The same flood on the infinitely tiled map reproduces the
  example's counts, but 26501365 steps are far too many to flood. The real map
  is square with an open border and an open row and column through `S`, so the
  flood swallows one more map copy in every direction per `width` steps: the
  counts after `offset + tiles * width` steps are a quadratic in `tiles`. Flood
  the first three tile-aligned counts, fit the quadratic, evaluate it.
  `O(width^3)` for those three floods, then constant time.
- Both runners use the real parameters (64 and 26501365 steps), so on the
  example they print those applied to the example. Its own counts (6, 10, 50,
  100 and 500 steps give 16, 50, 1594, 6536 and 167004 plots) are asserted in
  the tests, which also check the extrapolation against the flood on a map with
  the real structure; the example lacks the open row and column.

From `2023/rust`:

```sh
cargo run -q --bin day21 < examples/day21.txt # 42: the example at 64 steps
cargo run -q --bin day21 < examples/day21.txt # the fit for 26501365 steps
cargo test day21::
```
