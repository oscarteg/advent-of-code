# Day 7: Laboratories

## Problem

A grid where a beam starts at `S` and travels straight down. Hitting a
splitter (`^`) stops that beam and emits two new beams, one column left and
one column right of the splitter, continuing down from the next row. An
offshoot that would leave the grid is dropped.

**Part 1:** count the splits.

**Part 2:** every split forks a timeline (left and right each continue in
their own copy of events). Count the timelines that reach the bottom row.

---

*Solution spoilers below.*

## Approach

Track the active beams as a column-to-multiplicity map, processed one row
at a time. Whenever a column lands on a splitter, that counts as one split
(regardless of how many merged beams hit it together), and the beam count
moves to the neighbouring columns. For part 2 the same walk tracks how many
timelines occupy each column, merging counts when two timelines land on the
same column; the answer is the sum of those counts at the bottom row.

## Complexity

With `R` rows and `W` columns, each row touches at most `W` active columns,
so runtime is `O(R * W)` and space is `O(W)` for the per-row map.

## Run / test

```sh
cd 2025
zig test src/day7.zig
zig run src/day7.zig < src/examples/day7.txt
```

The example prints exactly:

```text
Part 1: 21
Part 2: 40
```
