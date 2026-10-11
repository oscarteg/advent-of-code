# Day 17: Clumsy Crucible

- **Part 1 & 2:** Dijkstra over the state `(x, y, last direction, run length)` instead of plain
  `(x, y)`, since the heat loss to reach a cell depends on how many consecutive blocks were
  already travelled in the current direction. From a state, the crucible may continue the same
  direction while `run < maxRun`, or turn once `run >= minRun`; it never reverses. The goal is
  reached only once `run >= minRun`. Part 1 uses `minRun = 0, maxRun = 3`; part 2 (the "ultra
  crucible") uses `minRun = 4, maxRun = 10`.
- **Complexity:** `O(width * height * 4 * maxRun * log(...))` for the heap-based search.

From `2023/go`:

```sh
go run ./day17 < day17/testdata/example.txt # Part 1: 102, Part 2: 94
go test -race -count=1 ./day17
```

Use the same command with your puzzle input redirected into stdin.
