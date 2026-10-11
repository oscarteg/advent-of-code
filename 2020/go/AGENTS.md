# 2020/go

## OVERVIEW
Go solutions for the 2020 puzzles: module `github.com/oscarteg/advent-of-code-2020` (go 1.23).
A single `main.go` dispatches `go run . <day>` to each day's exported `A`/`B` solvers and
prints both parts from one stdin read. testify stays for tests; cobra/viper are gone.

## WHERE TO LOOK
| Task | Location | Notes |
|------|----------|-------|
| Dispatch | `main.go` | `days` map: `"N": both(dayN.A, dayN.B)`; day 9 takes an optional preamble arg; day 25 is part 1 only |
| Input abstraction | `challenge/input.go` | `FromLiteral(s)`, `Lines()`, `LinesSlice()`, `ToGrid()` (stdin is read once in `main`) |
| A day | `challenge/dayN/` | `a.go`/`b.go` with `func A(*challenge.Input) (T, error)` (and `B`), `a_test.go`/`b_test.go`, `SOLUTION.md` |
| End-to-end runner test | `main_test.go` | builds the binary with `-race`, one example input per day, asserts both part lines |
| Real inputs | `../rust/input/2020/dayN.txt` | every day; days 1-3 also have `challenge/dayN/input.txt` |

## CONVENTIONS
- Output is `Part 1: X` / `Part 2: Y`; solver errors reach `main` and exit 1.
- Tests build input with `challenge.FromLiteral(...)` and assert with `testify/require`.

## COMMANDS
```bash
go run . 1 < ../rust/input/2020/day1.txt
go run . 9 5 < ../rust/input/2020/day9-example.txt
go test -race ./...
```
