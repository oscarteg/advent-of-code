# Day 12: Hot Springs

- **Part 1:** Parse each row into its spring states (`.`, `#`, `?`) and group
  sizes, then walk it with a memoised recursion over the pair *(spring
  position, next group)*. An unknown or operational spring can be skipped; a
  group may start at the current spring when the next `size` springs are all
  non-operational and the spring behind them is not damaged. A row counts when
  springs and groups run out together.
- **Part 2:** The same count after unfolding every row five times, joining the
  copies with `?` and repeating the group sizes. Memoisation turns the
  otherwise exponential search into a small table.
- Each row costs `O(springs * groups)` time and space, so the input as a whole
  is `O(rows * springs * groups)`.

From `2023/go`:

```sh
go run ./day12 < day12/testdata/example.txt
go test -race ./day12
go vet ./day12
```

The example prints exactly:

```text
Part 1: 21
Part 2: 525152
```
