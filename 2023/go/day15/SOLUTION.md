# Day 15: Lens Library

Part 1: run the HASH algorithm (`value = (value + byte) * 17 % 256` per byte)
over each comma-separated step and sum the results. `O(n)` in input length.

Part 2: parse each step as `label-` (remove) or `label=focal` (insert), hash
the label to pick one of 256 boxes, then insert/update/remove the lens in
that box's ordered slice. Sum `(box + 1) * (slot + 1) * focalLength` over all
boxes. `O(n * k)` where `k` is the average lenses per box, `O(n)` space.

From `2023/go`:

```sh
go run ./day15 < day15/testdata/example.txt
go test -race ./day15
go vet ./day15
```

The example prints exactly:

```text
Part 1: 1320
Part 2: 145
```
