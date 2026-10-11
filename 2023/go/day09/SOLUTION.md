# Day 9: Mirage Maintenance

- **Part 1:** For each history, recursively take successive differences until
  all values are zero, summing the last value of every row on the way back up
  to extrapolate the next reading. Sum over all histories.
- **Part 2:** Same recursion, extrapolating backward: each row's leading value
  minus the extrapolated previous value of its differences.
- Each history of length `n` makes `O(n)` recursive calls doing `O(n)` work
  each, so both parts are `O(n^2)` per history.

From `2023/go`:

```sh
go run ./day09 < day09/testdata/example.txt
go test -race ./day09
go vet ./day09
```

The example prints exactly:

```text
Part 1: 114
Part 2: 2
```
