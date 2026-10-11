# Advent of Code 2020 - Go

All 25 days. One command runs both parts, reading the puzzle input on stdin:

```bash
go run . 12 < ../rust/input/2020/day12.txt
go run . 9 5 < ../rust/input/2020/day9-example.txt   # day 9 takes an optional preamble length (the example uses 5)
```

Output is `Part 1: X` and `Part 2: Y` (day 25 has no part 2).

Each day lives in `challenge/dayN/` with its solvers (`A`, `B`), tests and a short `SOLUTION.md`.

```bash
go test -race ./...             # everything
go test -race ./challenge/day12 # one day
```
