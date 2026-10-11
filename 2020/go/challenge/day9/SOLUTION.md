# Day 9: Encoding Error

## Approach

**Part A:** Parse all numbers and check each one after the preamble (25 numbers) to see if it's a sum of two different numbers from the preceding preamble window. Return the first invalid number.

**Part B:** Given the invalid number from part A, find a contiguous sequence of at least 2 numbers that sum to it. Return the sum of the minimum and maximum values in that sequence.

## Complexity

- **Part A:** O(n × p²) where n is total numbers and p is preamble size. For each number after the preamble, check all pairs in the window.
- **Part B:** O(n²) worst case with early termination. For each start position, sum forward until target is found or exceeded.

## Commands

```bash
cd 2020/go
go test ./challenge/day9
go run . 9 < ../rust/input/2020/day9.txt
go run . 9 5 < ../rust/input/2020/day9-example.txt   # the example uses a preamble of 5
```
