# Day 1: Calorie Counting

Elves are separated by blank lines; each line within an elf's block is a
calorie count. Part 1 finds the elf carrying the most calories. Part 2 sums
the top three elves' totals.

Each elf's block is summed into a total, the totals are sorted descending,
and the top one (part 1) or three (part 2) are summed. Runtime is
`O(n log n)` for the sort; space is `O(n)` for the totals, where `n` is the
number of elves. Malformed calorie lines and empty input return an error.

From `2022/zig`:

```sh
zig test src/day_1.zig
zig run src/day_1.zig < src/examples/day_1.txt
```

The example prints exactly:

```text
Part 1: 24000
Part 2: 45000
```
