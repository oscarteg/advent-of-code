# Advent of Code 2023 - Rust

One crate, one binary per day. Every binary reads the puzzle input on stdin and
prints both parts:

```bash
cargo run -q --bin day05 < examples/day05.txt
cargo test             # everything
cargo test day05::     # one day
```

Layout: `src/dayNN/` holds the day's modules (`part1.rs`, `part2.rs`, `custom_error.rs`,
optional parser module), `src/bin/dayNN.rs` is its runner, `examples/` the official
examples, `input/` the committed real inputs (days 01-05), and `docs/dayNN.md` the
write-ups.
