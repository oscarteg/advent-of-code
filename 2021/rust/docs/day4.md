# Day 4: Giant Squid

Mark each drawn number on every board. A complete row or column wins; diagonals
do not. Record the first winning score for part 1 and the last for part 2.
Skip boards that already won so repeated draws cannot count them twice.

A score is the sum of unmarked numbers multiplied by the winning draw.
For `D` draws and `B` fixed-size boards, time is `O(D * B)` and space is `O(B)`.

From `2021/rust`:

```sh
cargo run --bin day4 < input/day4-example.txt
# Part 1: 4512
# Part 2: 1924
cargo test --bin day4
```
