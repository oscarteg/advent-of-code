# Day 2: Dive!

## Problem

Given a list of `forward`/`down`/`up` commands, compute horizontal position times depth.

**Part 1:** `forward` moves horizontally, `down`/`up` change depth directly.

**Part 2:** `down`/`up` change an "aim" value; `forward` moves horizontally and also
changes depth by `aim * distance`.

## Approach

Parse each line into a `command{dir, distance}`, then iterate once tracking
`(horizontal, depth)` for part 1 and `(horizontal, depth, aim)` for part 2.
`up` clamps depth/aim at zero instead of underflowing the unsigned counters.

## Complexity

- Time: O(n) for both parts
- Space: O(n) to hold the parsed commands

## Run

```bash
cd 2021/go
go run ./day02 < day02/testdata/example.txt
```

## Test

```bash
cd 2021/go
go test ./day02
```
