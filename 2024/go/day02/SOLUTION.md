# Day 2: Red-Nosed Reports

A report is safe when every adjacent pair of levels differs by 1-3 and the
sequence is consistently increasing or decreasing.

- **Part 1:** parse each line into ints, count safe reports.
- **Part 2:** a report also counts as safe if removing exactly one level
  makes it safe (the "problem dampener"), checked by brute force.

With `R` reports of length `L`, part 1 is `O(R * L)`, part 2 is `O(R * L^2)`.

From `2024/go`:

```sh
go run ./day02 < day02/testdata/example.txt
go test -race ./day02
go vet ./day02
```

The example prints exactly:

```text
Part 1: 2
Part 2: 4
```
