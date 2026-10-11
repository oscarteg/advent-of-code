# Day 24: Never Tell Me The Odds

- **Part 1:** Cramer's rule on `pa + s * va == pb + t * vb` in the xy plane for
  every pair. A zero determinant means parallel paths, a non-positive time means
  the paths crossed in the past. Both times and the crossing coordinates stay
  multiplied by the determinant, so the test-area comparison needs no division
  and stays exact in `i128`. `O(n^2)` pairs with `O(1)` work each.
- **Part 2:** Seen from the first hailstone the rock's path runs through the
  origin, so every other path spans a plane with it, with normal
  `offset x velocity`. The rock's direction lies in all of those planes, which
  makes it the cross product of two normals. Crossing that line with two paths
  gives two hit times and two points: their difference over the time span is the
  rock's velocity, and walking back to time zero gives its start position.
  Exact rationals over `i128` (`src/ratio.rs`) keep every step exact. `O(n)`.

The runners use the real test area, so part 1 prints `0` for `example.txt`; the
example's `7..=27` area is covered by the tests.

From `2023/rust`:

```sh
cargo test day24::                            # example: 2 crossings in 7..=27, rock sum 47
cargo run -q --bin day24 < input.txt
```
