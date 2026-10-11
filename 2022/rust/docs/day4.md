# Day 4: Camp Cleanup

## Problem

Each Elf is assigned a range of section IDs to clean. Input contains pairs of assignments.
We need to find overlapping assignments.

## Part 1

Count assignment pairs where one range fully contains the other.

**Approach**: Parse each pair as two ranges `(start, end)`. For ranges to be contained,
one must have start ≤ the other's start AND end ≥ the other's end.

**Complexity**: O(n) where n is the number of pairs.

## Part 2

Count assignment pairs where ranges overlap at all.

**Approach**: Two ranges overlap if `max(a.start, b.start) <= min(a.end, b.end)`.

**Complexity**: O(n).

## Run and Test

```bash
cd 2022/rust
cargo test --bin day4
cargo run -q --bin day4 < input/day4-example.txt
```
