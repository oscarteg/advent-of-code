# Day 8: Seven Segment Search

Each signal pattern becomes a seven-bit mask. Digits `1`, `4`, `7`, and `8` are
identified by their unique segment counts. Among the six-segment patterns, `9`
contains `4`, `0` contains `1`, and the remaining pattern is `6`. Among the
five-segment patterns, `3` contains `1`, `5` is contained by `6`, and the
remaining pattern is `2`.

Part 1 counts output masks with a unique segment count. Part 2 uses the decoded
masks to turn each four-digit output into a number.

The work is O(n) in the number of entries with O(1) extra space per entry: each
entry has exactly ten signals and four outputs.

```sh
cd 2021/rust
cargo test --bin day8
cargo run -q --bin day8 < input/day8-example.txt
```
