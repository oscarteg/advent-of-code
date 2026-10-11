# Day 20: Trench Map

Parse the 512-entry enhancement algorithm and the starting image as a sparse
set of lit pixel coordinates. Each enhancement step looks at every pixel in a
1-cell-wider bounding box, builds a 9-bit index from its 3x3 neighborhood
(out-of-bounds pixels use the current infinite background value), and looks
that index up in the algorithm to decide whether the pixel is lit next.

Because `algorithm[0]` can be `#`, the infinite background itself flips each
step (dark then light then dark...). The background for the next step is
`algorithm[511]` when currently lit, otherwise `algorithm[0]`; the returned
lit-pixel count only covers the finite image, which is correct after an even
number of steps.

Part 1 runs 2 enhancements, part 2 runs 50. Each step is `O(N)` in the current
image area, which grows by a constant border per step, so `k` steps cost
`O(k^2)` overall; space is `O(k^2)` for the lit-pixel set.

From `2021/go`:

```sh
go test -race ./day20
go vet ./day20
go run ./day20 < day20/testdata/example.txt
```

The example prints exactly:

```text
Part 1: 35
Part 2: 3351
```
