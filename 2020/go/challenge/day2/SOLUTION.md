# Day 2: Password Philosophy

Each line is parsed with a regex into a `policy{low, high, letter, password}`.

## Part A

Valid when the count of `letter` in `password` is between `low` and `high`
inclusive. O(n*m) over all lines and characters.

## Part B

Valid when exactly one of the 1-indexed positions `low`/`high` in `password`
holds `letter` (XOR of the two checks). Same complexity as part A.

## Run

```bash
go run . 2 < ../rust/input/2020/day2.txt
```

## Test

```bash
go test -race -count=1 ./challenge/day2
```
