# Day 2: Dive!

## Problem

Given a list of `forward`/`down`/`up` commands, compute horizontal position times depth.

**Part 1:** `forward` moves horizontally, `down`/`up` change depth directly.

**Part 2:** `down`/`up` change an "aim" value; `forward` moves horizontally and also
changes depth by `aim * distance`.

## Approach

Parse each line into a `Command` enum, then fold over the commands accumulating
`(horizontal, depth)` for part 1 and `(horizontal, depth, aim)` for part 2.
Depth and aim use `saturating_sub` to avoid underflow on malformed `up` commands.

## Complexity

- Time: O(n) for both parts
- Space: O(n) to hold the parsed commands

## Run

```bash
cd 2021/rust
cargo run -q --bin day2 < input/day2-example.txt
```

## Test

```bash
cd 2021/rust
cargo test --bin day2
```
