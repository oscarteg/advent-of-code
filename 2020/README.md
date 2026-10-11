# Advent of Code 2020

All 25 days, solved in Rust and Go.

## Languages

| Lang | Project | Run | Test |
|------|---------|-----|------|
| Rust | `rust/` | `cargo run -q --bin dayN < input/2020/dayN.txt` | `cargo test --bin dayN` |
| Go | `go/` | `go run . N < ../rust/input/2020/dayN.txt` (prints both parts) | `go test -race -shuffle=on -count=1 ./...` |

Both read the puzzle input on stdin and print `Part 1: X` / `Part 2: Y`. Real inputs live under `rust/input/2020/`; days 1-3 additionally ship a copy
under `go/challenge/dayN/input.txt`. `rust/input/2020/dayN-example.txt` holds the official example
for each Rust day.

## Days

| Day | Title | Rust | Go |
|-----|-------|------|-----|
| 1 | Report Repair | [docs](rust/docs/day1.md) | [SOLUTION](go/challenge/day1/SOLUTION.md) |
| 2 | Password Philosophy | [docs](rust/docs/day2.md) | [SOLUTION](go/challenge/day2/SOLUTION.md) |
| 3 | Toboggan Trajectory | [docs](rust/docs/day3.md) | [SOLUTION](go/challenge/day3/SOLUTION.md) |
| 4 | Passport Processing | [docs](rust/docs/day4.md) | [SOLUTION](go/challenge/day4/SOLUTION.md) |
| 5 | Binary Boarding | [docs](rust/docs/day5.md) | [SOLUTION](go/challenge/day5/SOLUTION.md) |
| 6 | Custom Customs | [docs](rust/docs/day6.md) | [SOLUTION](go/challenge/day6/SOLUTION.md) |
| 7 | Handy Haversacks | [docs](rust/docs/day7.md) | [SOLUTION](go/challenge/day7/SOLUTION.md) |
| 8 | Handheld Halting | [docs](rust/docs/day8.md) | [SOLUTION](go/challenge/day8/SOLUTION.md) |
| 9 | Encoding Error | [docs](rust/docs/day9.md) | [SOLUTION](go/challenge/day9/SOLUTION.md) |
| 10 | Adapter Array | [docs](rust/docs/day10.md) | [SOLUTION](go/challenge/day10/SOLUTION.md) |
| 11 | Seating System | [docs](rust/docs/day11.md) | [SOLUTION](go/challenge/day11/SOLUTION.md) |
| 12 | Rain Risk | [docs](rust/docs/day12.md) | [SOLUTION](go/challenge/day12/SOLUTION.md) |
| 13 | Shuttle Search | [docs](rust/docs/day13.md) | [SOLUTION](go/challenge/day13/SOLUTION.md) |
| 14 | Docking Data | [docs](rust/docs/day14.md) | [SOLUTION](go/challenge/day14/SOLUTION.md) |
| 15 | Rambunctious Recitation | [docs](rust/docs/day15.md) | [SOLUTION](go/challenge/day15/SOLUTION.md) |
| 16 | Ticket Translation | [docs](rust/docs/day16.md) | [SOLUTION](go/challenge/day16/SOLUTION.md) |
| 17 | Conway Cubes | [docs](rust/docs/day17.md) | [SOLUTION](go/challenge/day17/SOLUTION.md) |
| 18 | Operation Order | [docs](rust/docs/day18.md) | [SOLUTION](go/challenge/day18/SOLUTION.md) |
| 19 | Monster Messages | [docs](rust/docs/day19.md) | [SOLUTION](go/challenge/day19/SOLUTION.md) |
| 20 | Jurassic Jigsaw | [docs](rust/docs/day20.md) | [SOLUTION](go/challenge/day20/SOLUTION.md) |
| 21 | Allergen Assessment | [docs](rust/docs/day21.md) | [SOLUTION](go/challenge/day21/SOLUTION.md) |
| 22 | Crab Combat | [docs](rust/docs/day22.md) | [SOLUTION](go/challenge/day22/SOLUTION.md) |
| 23 | Crab Cups | [docs](rust/docs/day23.md) | [SOLUTION](go/challenge/day23/SOLUTION.md) |
| 24 | Lobby Layout | [docs](rust/docs/day24.md) | [SOLUTION](go/challenge/day24/SOLUTION.md) |
| 25 | Combo Breaker | [docs](rust/docs/day25.md) | [SOLUTION](go/challenge/day25/SOLUTION.md) |

Day 9's Go CLI takes an optional preamble argument (`go run . 9 5 < in.txt`; it defaults to the
real-input preamble of 25, the official example uses 5). Day 25 has no part B.
