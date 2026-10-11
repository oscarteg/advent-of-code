# Day 4: Camp Cleanup

## Problem

Each Elf is assigned a range of section IDs to clean. Input contains pairs of assignments.
We need to find overlapping assignments.

## Part 1

Count assignment pairs where one range fully contains the other.

**Approach**: Parse each pair as two ranges with start and end. Check if one range's
boundaries fully enclose the other's.

**Complexity**: O(n) where n is the number of pairs.

## Part 2

Count assignment pairs where ranges overlap at all.

**Approach**: Two ranges overlap if `max(a.start, b.start) <= min(a.end, b.end)`.

**Complexity**: O(n).

## Run and Test

```bash
cd 2022/go
go test -race ./day04
go run ./day04 < day04/testdata/example.txt
```
