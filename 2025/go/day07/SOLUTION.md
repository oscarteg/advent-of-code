# Day 7: Laboratories

A beam starts at `S` and travels straight down. A splitter (`^`) stops the
beam and emits two beams, immediately left and right, continuing down from
the next row.

Track active beams as a column-to-multiplicity map, walked one row at a
time. Landing on a splitter counts as one split (part 1) and forwards its
multiplicity to the neighbouring columns; landing elsewhere just carries the
count down unchanged. A split at the grid edge drops the offshoot that
falls off the grid. Part 2 sums the multiplicities once the walk reaches
the bottom row, since each split also forks a timeline.

With `R` rows and `W` columns, each row touches at most `W` active columns,
so runtime is `O(R * W)` and space is `O(W)`.

From `2025/go`:

```sh
go run ./day07 < day07/testdata/example.txt
go test -race ./day07
go vet ./day07
```

The example prints exactly:

```text
Part 1: 21
Part 2: 40
```
