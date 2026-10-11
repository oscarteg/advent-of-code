# Day 24: Never Tell Me The Odds

- **Part 1:** Cramer's rule on `pa + s*va == pb + t*vb` in the xy plane for every
  pair. A zero determinant means parallel paths, a non-positive time means they
  crossed in the past; the crossing itself is compared against the test area.
  `O(n^2)` pairs with `O(1)` work each.
- **Part 2:** Seen from the first hailstone the rock's path runs through the
  origin, so every other path spans a plane with it (normal `offset x velocity`).
  The rock's direction lies in all of those planes, which makes it the cross
  product of two normals. Crossing that line with two paths gives two hit times
  and two points: their difference over the time span is the rock's velocity, and
  walking back to time zero gives its start position. `O(n)`.
- Every step runs on `math/big.Rat`, because the real test area needs more than
  64 bits and a `float64` would misjudge crossings near the boundary.

The runner uses the real test area, so part 1 prints `0` for the example; its
`7..27` area is covered by the tests.

From `2023/go`:

```sh
go vet ./day24 && go test -race ./day24
go run ./day24 < day24/testdata/example.txt   # Part 1: 0, Part 2: 47
go run ./day24 < input.txt
```
