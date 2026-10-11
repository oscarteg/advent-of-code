# Day 2: Password Philosophy

Each line is parsed with a regex into a `Policy { low, high, letter, password }`.

## Part 1

A password is valid when the number of occurrences of `letter` in `password`
falls within the inclusive `low..=high` range. O(n * m) over all lines and
characters; O(1) extra space per line.

## Part 2

A password is valid when exactly one of the 1-indexed positions `low`/`high`
holds `letter` (XOR of the two position checks). Same complexity as part 1.

## Run

```sh
cargo run -q --bin day2 < input/2020/day2-example.txt
```

## Test

```sh
cargo test --bin day2
```
