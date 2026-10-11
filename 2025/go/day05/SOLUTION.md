# Day 5: Cafeteria

## Approach

**Part 1:** Count how many available ingredient IDs fall within at least one fresh range. For each ID, check if it's contained in any of the given ranges.

**Part 2:** Count the total number of distinct integer IDs covered by all ranges after merging overlaps. Sort ranges by start, then merge adjacent or overlapping ranges by extending the end point if needed, finally sum the coverage of each merged range.

## Complexity

- Part 1: O(n·m) where n = number of IDs, m = number of ranges
- Part 2: O(m log m) for sorting + O(m) for merging

## Commands

```bash
cd 2025/go
go test -race -count=1 ./day05
go run ./day05 < day05/testdata/example.txt
```
