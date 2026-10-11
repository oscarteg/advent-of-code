# Day 6: Tuning Trouble

## Problem

Find the start-of-packet and start-of-message markers in a datastream.

- **Part 1:** Find the position where the first sequence of 4 distinct characters ends.
- **Part 2:** Find the position where the first sequence of 14 distinct characters ends.

## Approach

Both parts use a sliding window algorithm:
1. Maintain a window of the required size (4 for part 1, 14 for part 2)
2. Check if all characters in the window are distinct using a map
3. Return the position when all characters are unique

Maps provide O(1) average insertion, deletion and lookup.

## Complexity

- Time: O(n * window_size) where n is the input length
- Space: O(window_size) for the map

## Commands

```bash
cd 2022/go
go test -race -count=1 ./day06
go run ./day06 < day06/testdata/example.txt
```
