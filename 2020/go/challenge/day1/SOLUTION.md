# Day 1: Report Repair

## Problem Summary

Part A: Find two entries in the expense report that sum to 2020 and multiply them together.

Part B: Find three entries that sum to 2020 and multiply them together.

## Approach

Both parts use a map for O(1) average-case lookups.

**Part A:** Build a map of all entries. For each entry x, check if (2020 - x) exists in the map. Time: O(n).

**Part B:** Build a map of all entries and a list. For each pair (x, y), check if (2020 - x - y) exists in the map. Time: O(n²).

## Complexity

- Time: O(n) for part a, O(n²) for part b
- Space: O(n) for the map and list

## Commands

```bash
cd 2020/go
go test -race ./challenge/day1
go run . 1 < ../rust/input/2020/day1-example.txt
go run . 1 < ../rust/input/2020/day1.txt
```
