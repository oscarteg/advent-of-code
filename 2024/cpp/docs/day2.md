# Day 2: Red-Nosed Reports (C++)

Each report is a list of levels. It is safe when every adjacent pair differs
by 1-3 and the whole sequence is consistently increasing or decreasing.

- **Part 1:** parse each line into integers and count safe reports.
- **Part 2:** a report also counts as safe if removing exactly one level
  makes it safe (the "problem dampener"). Checked by brute force: try
  removing each index in turn.

With `R` reports of length `L`, part 1 is `O(R * L)` and part 2 is
`O(R * L^2)` since each of the `L` removals re-checks safety.

Run and test commands (from `2024`):

The example prints exactly:

```text
Part 1: 2
Part 2: 4
```

## Run / test

```bash
cd 2024/cpp
make test day=day2
make -s run day=day2 < examples/day2.txt
```
