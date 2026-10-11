# Day 10: Pipe Maze

- **Parsing:** the grid is kept as raw lines; each byte maps to the two
  directions a pipe connects (or none for `.`/`S`).
- **Finding the loop:** check which neighbors of `S` have a pipe connecting
  back toward it (exactly two must); walk the loop one step at a time,
  always exiting through the connection that isn't where we came from, until
  we return to `S`. This infers `S`'s real shape without ever storing it.
- **Part 1:** the farthest tile is half the loop's length.
- **Part 2:** the shoelace formula gives the loop's area, then Pick's theorem
  (`area = interior + boundary/2 - 1`) recovers the enclosed tile count
  directly, without scanning the whole grid.
- Both parts run in `O(rows * cols)`.

From `2023/go`:

```sh
go run ./day10 < day10/testdata/example.txt
go test -race ./day10
go vet ./day10
```

The example prints:

```text
Part 1: 8
Part 2: 1
```
