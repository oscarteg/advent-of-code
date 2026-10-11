# Day 19: Linen Layout

The first line lists the available towel patterns; every later line is a design that must
be built by concatenating patterns.

- **Part 1:** count the designs that can be built at all.
- **Part 2:** sum, over all designs, the number of different ways to build them.

`countWays` tries every pattern that prefixes the remaining design and recurses on the rest,
memoised per suffix, so each suffix is solved once.

## Complexity

O(designs x length x patterns) time with memoisation; O(length) memo entries per design.

## Run / test

```bash
cd 2024/go
go test -race ./day19
go run ./day19 < day19/testdata/example.txt
```
