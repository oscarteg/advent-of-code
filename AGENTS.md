# PROJECT KNOWLEDGE BASE

**Updated:** 2026-10-11
**Branch:** main

## OVERVIEW
Advent of Code 2020-2025. Every year keeps the language(s) it was started in and also has a Go
version. Each year/language directory is an independent project with its own toolchain; there
is no shared code and no root build.

## LAYOUT AND RUNNERS
Every runner reads the puzzle input on stdin. Day `N` is unpadded, `NN` is zero-padded.
The root `./run <year> <day> <lang> < input.txt` script dispatches to the right project
(interactive without a redirect: paste the input, end with Ctrl-D).

| Year | Lang | Project | Files for day N | Run | Test |
|------|------|---------|-----------------|-----|------|
| 2020 | Rust | `2020/rust` | `src/days/dayN.rs`, `src/bin/dayN.rs`, `docs/dayN.md`, `input/2020/dayN-example.txt` | `cargo run -q --bin dayN < in.txt` | `cargo test --bin dayN` |
| 2020 | Go | `2020/go` | `challenge/dayN/{a,b}.go`, `*_test.go`, `SOLUTION.md` | `go run . N < in.txt` (both parts; day 9: optional preamble arg) | `go test ./challenge/dayN` |
| 2021, 2022 | Rust | `<year>/rust` | as 2020, examples in `input/dayN-example.txt` | `cargo run -q --bin dayN < in.txt` | `cargo test --bin dayN` |
| 2021-2025 | Go | `<year>/go` | `dayNN/{main.go,<topic>.go,<topic>_test.go,testdata/example.txt,SOLUTION.md}` | `go run ./dayNN < in.txt` | `go test ./dayNN` |
| 2022 | Zig | `2022/zig` | `src/day_N.zig`, `src/examples/day_N.txt`, `docs/day_N.md` | `zig run src/day_N.zig < in.txt` | `zig test src/day_N.zig` |
| 2023 | Rust | `2023/rust` single crate | `src/dayNN/{mod,part1,part2,custom_error}.rs`, `src/bin/dayNN.rs`, `examples/dayNN.txt`, `docs/dayNN.md` | `cargo run -q --bin dayNN < in.txt` | `cargo test dayNN::` |
| 2024 | C++23 | `2024/cpp` (Makefile) | `dayN.hpp`, `dayN.cpp`, `dayN_test.cpp`, `examples/dayN.txt`, `docs/dayN.md` | `make -s run day=dayN < in.txt` | `make test day=dayN` |
| 2024 | Zig | `2024/zig` | `dayN.zig`, `examples/dayN.txt`, `docs/dayN.md` | `zig run dayN.zig < in.txt` | `zig test dayN.zig` |
| 2025 | Zig | `2025` | `src/dayN.zig`, `src/examples/dayN.txt`, `docs/dayN.md` | `zig run src/dayN.zig < in.txt` | `zig test src/dayN.zig` |

Output: `Part 1: X` and `Part 2: Y` lines; every day 25 has a single part.

## WHERE TO LOOK
| Task | Location | Notes |
|------|----------|-------|
| Reference Rust day (single crate) | `2021/rust/src/days/day4.rs` + `src/bin/day4.rs` | `pub fn solve(&str) -> anyhow::Result<(A, B)>`, bin includes the module with `#[path]` |
| Reference Go day | `2021/go/day04/` | stdlib only, `solve(input string)`, table tests |
| Reference 2023 day | `2023/rust/src/day05/` + `src/bin/day05.rs` | stdin bin, typed parser module, `AocError` variants |
| Reference Zig 0.16 I/O | `2025/src/utils.zig` `readStdin`, `2025/src/day1.zig` `main` | |
| C++ test helpers | `2024/cpp/check.hpp` | `check::equal`, `check::throws`, `check::exit_code` |
| 2020 Go dispatch | `2020/go/main.go` (`days` map), `2020/go/AGENTS.md` | |

## CONVENTIONS
- Parse input into typed values at the boundary; malformed input returns an error (Rust
  `anyhow`/`AocError`, Go `error`, Zig error union, C++ `std::invalid_argument`), never a panic.
- Tests: official examples for both parts, at least one edge case, and malformed input.
- One short Markdown file per solution: approach per part, complexity, run/test commands.
  Never copy puzzle prose into the repo; example inputs are fine.
- Real puzzle inputs are not added to git (Advent of Code policy). Existing committed ones stay.
- Only stdlib in new Go modules (go 1.23). 2020 Go keeps testify for its tests.

## TOOLCHAINS (mise / devbox, "latest")
- Zig 0.16: `pub fn main(init: std.process.Init)`, stdin via
  `std.Io.File.stdin().reader(init.io, &buf)` + `.interface.allocRemaining(arena, .unlimited)`,
  buffered `std.Io.File.stdout().writer(init.io, &buf)` + `flush()`. `std.ArrayList` is unmanaged
  (`.empty`, `append(allocator, x)`). `std.fs.File` and `std.io` no longer exist.
- Go toolchain 1.27, Rust stable, Apple clang (`-std=c++23`).
- `cargo-generate`, `just`, `cargo-nextest` and `zls` are not installed.

## ANTI-PATTERNS (THIS PROJECT)
- Don't trust `zig test` alone: it only compiles code reachable from tests, so always `zig run` too.
- Don't run crate- or module-wide formatters (`cargo fmt`, `go fmt ./...`) for a one-day change;
  format the files you touched.
- Don't add a root workspace or build that ties years together.
- `2024/.gitignore` ignores only `input/` (real inputs) and `zig/scratch/`; the Makefile pins
  `-std=c++23` with `override` because the shell exports `CXXFLAGS=-std=c++17`. `make clean`
  removes only `cpp/build/`.

## COMMANDS
```bash
cd 2021/rust && cargo test && cargo run -q --bin day4 < input/day4-example.txt
cd 2021/go && go test -race ./... && go run ./day04 < day04/testdata/example.txt
cd 2025 && zig test src/day1.zig && zig run src/day1.zig < input.txt
cd 2024/cpp && make test && make -s run day=day1 < examples/day1.txt
cd 2023/rust && cargo test day05:: && cargo run -q --bin day05 < examples/day05.txt
cd 2020/go && go test -race ./... && go run . 1 < ../rust/input/2020/day1.txt
```

## NOTES
- `2020/go` lived in `2022/go` by mistake (its module path always said 2020). `2022/go` is now the
  2022 Go module.
- Each year has a `README.md` with its run/test commands and a day table; the root `README.md`
  links them.
