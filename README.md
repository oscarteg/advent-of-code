# Advent of Code

My [Advent of Code](https://adventofcode.com) solutions for 2020-2025. Each year is solved in the
language(s) I started it in, plus Go for every year.

| Year | Languages | Guide |
|------|-----------|-------|
| 2020 | Rust, Go | [2020/README.md](2020/README.md) |
| 2021 | Rust, Go | [2021/README.md](2021/README.md) |
| 2022 | Rust, Zig, Go | [2022/README.md](2022/README.md) |
| 2023 | Rust, Go | [2023/README.md](2023/README.md) |
| 2024 | C++, Zig, Go | [2024/README.md](2024/README.md) |
| 2025 | Zig, Go | [2025/README.md](2025/README.md) |

Every day has a solution, tests (the official examples plus edge and malformed-input cases), a
runner and a short write-up of the approach. Each year/language directory is its own project;
there is no root build.

## Running a day

Every runner reads the puzzle input on stdin and prints `Part 1: ...` / `Part 2: ...`.
The `./run` script picks the right project for you:

```bash
./run 2025 day12 zig < input.txt
pbpaste | ./run 2023 day5 rust   # input straight from the clipboard
./run 2020 day9 go               # no redirect: paste the input, end with Ctrl-D
```

Or run a project directly:

```bash
cd 2020/go   && go run . 1 < ../rust/input/2020/day1.txt
cd 2021/rust && cargo run -q --bin day4 < input/day4-example.txt
cd 2021/go   && go run ./day04 < day04/testdata/example.txt
cd 2022/zig  && zig run src/day_7.zig < src/examples/day_7.txt
cd 2023/rust && cargo run -q --bin day05 < examples/day05.txt
cd 2024/cpp  && make -s run day=day7 < examples/day7.txt
cd 2025      && zig run src/day3.zig < src/examples/day3.txt
```

Each year's guide lists the exact run and test commands per language.

## Toolchains

Rust (stable), Go 1.23+, Zig 0.16 and clang with C++23, installed with `mise` or `devbox shell`.

## Inputs

The example inputs from the puzzle pages sit next to the tests. Real puzzle inputs are
personal ([Advent of Code asks](https://adventofcode.com/about) not to share them): inputs that
were committed before stay, new ones are not added, and 2024 keeps its inputs in the git-ignored
`2024/input/`.
