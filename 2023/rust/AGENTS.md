# 2023/rust

## OVERVIEW
One crate `aoc2023`, no workspace. `src/dayNN/` holds each day's modules (`mod.rs`,
`part1.rs`, `part2.rs`, `custom_error.rs`, optional parser module); `src/bin/dayNN.rs`
reads stdin and prints both parts. Official examples live in `examples/dayNN.txt`
(`dayNN-2.txt` etc. for extra examples), committed real inputs in `input/`, and the
write-ups in `docs/dayNN.md`.

## NEW DAY
1. `mkdir src/dayNN` with `mod.rs` (`pub mod custom_error; pub mod part1; pub mod part2;`).
2. Add `pub mod dayNN;` to `src/lib.rs` and copy `src/bin/day05.rs` to `src/bin/dayNN.rs`.
3. `process(&str) -> miette::Result<String, AocError>`; the official example is tested via
   `include_str!("../../examples/dayNN.txt")` inside `partN.rs`.

## CONVENTIONS
- `process` returns the answer as a `String`; bad input is an `AocError`, never a panic.
- Day 25 has no part 2; its bin prints only `Part 1`.

## COMMANDS
```bash
cargo test day05::
cargo clippy --all-targets -- -D warnings
cargo run -q --bin day05 < examples/day05.txt
```
