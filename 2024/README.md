# Advent of Code 2024

All 25 days in C++23, Zig 0.16 and Go, one directory per language. Every runner reads the
puzzle input on stdin and prints `Part 1: X` / `Part 2: Y`.

```
2024/
├── cpp/     # dayN.hpp (solution), dayN.cpp (runner), dayN_test.cpp, Makefile, examples/, docs/
├── zig/     # dayN.zig (solution, runner and tests), examples/, docs/
├── go/      # dayNN/ packages with tests, testdata/ and SOLUTION.md
└── input/   # your real puzzle inputs (git-ignored)
```

## Run and test

| Language | Directory | Run one day | Test one day | Test everything |
|----------|-----------|-------------|--------------|-----------------|
| C++23 | `cpp/` | `make -s run day=day3 < ../input/day3_input.txt` | `make test day=day3` | `make test` |
| Zig 0.16 | `zig/` | `zig run day3.zig < ../input/day3_input.txt` | `zig test day3.zig` | `for f in day*.zig; do zig test "$f" \|\| break; done` |
| Go | `go/` | `go run ./day03 < ../input/day3_input.txt` | `go test ./day03` | `go test -race ./...` |

Example inputs live in `cpp/examples/`, `zig/examples/` and `go/dayNN/testdata/`.

## Days

| Day | Title | C++ | Zig | Go |
|-----|-------|-----|-----|----|
| 1 | Historian Hysteria | [day1.hpp](cpp/day1.hpp) · [doc](cpp/docs/day1.md) | [day1.zig](zig/day1.zig) · [doc](zig/docs/day1.md) | [day01](go/day01/) · [doc](go/day01/SOLUTION.md) |
| 2 | Red-Nosed Reports | [day2.hpp](cpp/day2.hpp) · [doc](cpp/docs/day2.md) | [day2.zig](zig/day2.zig) · [doc](zig/docs/day2.md) | [day02](go/day02/) · [doc](go/day02/SOLUTION.md) |
| 3 | Mull It Over | [day3.hpp](cpp/day3.hpp) · [doc](cpp/docs/day3.md) | [day3.zig](zig/day3.zig) · [doc](zig/docs/day3.md) | [day03](go/day03/) · [doc](go/day03/SOLUTION.md) |
| 4 | Ceres Search | [day4.hpp](cpp/day4.hpp) · [doc](cpp/docs/day4.md) | [day4.zig](zig/day4.zig) · [doc](zig/docs/day4.md) | [day04](go/day04/) · [doc](go/day04/SOLUTION.md) |
| 5 | Print Queue | [day5.hpp](cpp/day5.hpp) · [doc](cpp/docs/day5.md) | [day5.zig](zig/day5.zig) · [doc](zig/docs/day5.md) | [day05](go/day05/) · [doc](go/day05/SOLUTION.md) |
| 6 | Guard Gallivant | [day6.hpp](cpp/day6.hpp) · [doc](cpp/docs/day6.md) | [day6.zig](zig/day6.zig) · [doc](zig/docs/day6.md) | [day06](go/day06/) · [doc](go/day06/SOLUTION.md) |
| 7 | Bridge Repair | [day7.hpp](cpp/day7.hpp) · [doc](cpp/docs/day7.md) | [day7.zig](zig/day7.zig) · [doc](zig/docs/day7.md) | [day07](go/day07/) · [doc](go/day07/SOLUTION.md) |
| 8 | Resonant Collinearity | [day8.hpp](cpp/day8.hpp) · [doc](cpp/docs/day8.md) | [day8.zig](zig/day8.zig) · [doc](zig/docs/day8.md) | [day08](go/day08/) · [doc](go/day08/SOLUTION.md) |
| 9 | Disk Fragmenter | [day9.hpp](cpp/day9.hpp) · [doc](cpp/docs/day9.md) | [day9.zig](zig/day9.zig) · [doc](zig/docs/day9.md) | [day09](go/day09/) · [doc](go/day09/SOLUTION.md) |
| 10 | Hoof It | [day10.hpp](cpp/day10.hpp) · [doc](cpp/docs/day10.md) | [day10.zig](zig/day10.zig) · [doc](zig/docs/day10.md) | [day10](go/day10/) · [doc](go/day10/SOLUTION.md) |
| 11 | Plutonian Pebbles | [day11.hpp](cpp/day11.hpp) · [doc](cpp/docs/day11.md) | [day11.zig](zig/day11.zig) · [doc](zig/docs/day11.md) | [day11](go/day11/) · [doc](go/day11/SOLUTION.md) |
| 12 | Garden Groups | [day12.hpp](cpp/day12.hpp) · [doc](cpp/docs/day12.md) | [day12.zig](zig/day12.zig) · [doc](zig/docs/day12.md) | [day12](go/day12/) · [doc](go/day12/SOLUTION.md) |
| 13 | Claw Contraption | [day13.hpp](cpp/day13.hpp) · [doc](cpp/docs/day13.md) | [day13.zig](zig/day13.zig) · [doc](zig/docs/day13.md) | [day13](go/day13/) · [doc](go/day13/SOLUTION.md) |
| 14 | Restroom Redoubt | [day14.hpp](cpp/day14.hpp) · [doc](cpp/docs/day14.md) | [day14.zig](zig/day14.zig) · [doc](zig/docs/day14.md) | [day14](go/day14/) · [doc](go/day14/SOLUTION.md) |
| 15 | Warehouse Woes | [day15.hpp](cpp/day15.hpp) · [doc](cpp/docs/day15.md) | [day15.zig](zig/day15.zig) · [doc](zig/docs/day15.md) | [day15](go/day15/) · [doc](go/day15/SOLUTION.md) |
| 16 | Reindeer Maze | [day16.hpp](cpp/day16.hpp) · [doc](cpp/docs/day16.md) | [day16.zig](zig/day16.zig) · [doc](zig/docs/day16.md) | [day16](go/day16/) · [doc](go/day16/SOLUTION.md) |
| 17 | Chronospatial Computer | [day17.hpp](cpp/day17.hpp) · [doc](cpp/docs/day17.md) | [day17.zig](zig/day17.zig) · [doc](zig/docs/day17.md) | [day17](go/day17/) · [doc](go/day17/SOLUTION.md) |
| 18 | RAM Run | [day18.hpp](cpp/day18.hpp) · [doc](cpp/docs/day18.md) | [day18.zig](zig/day18.zig) · [doc](zig/docs/day18.md) | [day18](go/day18/) · [doc](go/day18/SOLUTION.md) |
| 19 | Linen Layout | [day19.hpp](cpp/day19.hpp) · [doc](cpp/docs/day19.md) | [day19.zig](zig/day19.zig) · [doc](zig/docs/day19.md) | [day19](go/day19/) · [doc](go/day19/SOLUTION.md) |
| 20 | Race Condition | [day20.hpp](cpp/day20.hpp) · [doc](cpp/docs/day20.md) | [day20.zig](zig/day20.zig) · [doc](zig/docs/day20.md) | [day20](go/day20/) · [doc](go/day20/SOLUTION.md) |
| 21 | Keypad Conundrum | [day21.hpp](cpp/day21.hpp) · [doc](cpp/docs/day21.md) | [day21.zig](zig/day21.zig) · [doc](zig/docs/day21.md) | [day21](go/day21/) · [doc](go/day21/SOLUTION.md) |
| 22 | Monkey Market | [day22.hpp](cpp/day22.hpp) · [doc](cpp/docs/day22.md) | [day22.zig](zig/day22.zig) · [doc](zig/docs/day22.md) | [day22](go/day22/) · [doc](go/day22/SOLUTION.md) |
| 23 | LAN Party | [day23.hpp](cpp/day23.hpp) · [doc](cpp/docs/day23.md) | [day23.zig](zig/day23.zig) · [doc](zig/docs/day23.md) | [day23](go/day23/) · [doc](go/day23/SOLUTION.md) |
| 24 | Crossed Wires | [day24.hpp](cpp/day24.hpp) · [doc](cpp/docs/day24.md) | [day24.zig](zig/day24.zig) · [doc](zig/docs/day24.md) | [day24](go/day24/) · [doc](go/day24/SOLUTION.md) |
| 25 | Code Chronicle | [day25.hpp](cpp/day25.hpp) · [doc](cpp/docs/day25.md) | [day25.zig](zig/day25.zig) · [doc](zig/docs/day25.md) | [day25](go/day25/) · [doc](go/day25/SOLUTION.md) |
