# Day 13: Point of Incidence

- **Part 1:** For each pattern, test every possible horizontal and vertical
  mirror line. A line is valid when the rows (or columns) reflected across it
  match exactly on the overlapping span. Sum `100 * rows above` for horizontal
  mirrors and `columns to the left` for vertical ones.
- **Part 2:** Same search, but a line is valid when the reflected cells differ
  by exactly one ("smudge") across the whole pattern instead of zero.
- Runtime is `O(patterns * cells)`: each candidate line compares at most the
  full grid once, and there are at most `rows + columns` candidate lines.

From `2023/go`:

```sh
go run ./day13 < day13/testdata/example.txt
go test -race ./day13
go vet ./day13
```

The example prints exactly:

```text
Part 1: 405
Part 2: 400
```
