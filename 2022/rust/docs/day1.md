# Day 1: Calorie Counting

Elves are separated by blank lines; each line within an elf's block is a
calorie count. Part 1 finds the elf carrying the most calories. Part 2 sums
the top three elves' totals.

Input is parsed into one total per elf by splitting on blank lines and
summing each block's numbers, then sorted descending. Runtime is
`O(n log n)` for the sort; space is `O(n)` for the totals, where `n` is the
number of elves. Malformed calorie lines return an error via `anyhow`.

From `2022/rust`:

```sh
cargo test --bin day1
cargo run -q --bin day1 < input/day1-example.txt
```

The example prints exactly:

```text
Part 1: 24000
Part 2: 45000
```
