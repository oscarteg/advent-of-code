# Day 1: Sonar Sweep

## Problem

Given a list of sonar depth measurements, count increases.

**Part 1:** Count how many times a depth measurement increases from the previous measurement.

**Part 2:** Count how many times the sum of a 3-measurement sliding window increases from the previous window's sum. (Equivalently: count how many times `depths[i] > depths[i-3]`.)

## Approach

**Part 1:** Iterate through depths and count whenever the current depth exceeds the previous one.

**Part 2:** Iterate through depths starting from index 3 and count whenever the current depth exceeds the depth 3 positions back (which is equivalent to comparing sliding window sums).

## Complexity

- Time: O(n) for both parts
- Space: O(n) for storing depths

## Run

```bash
cd 2021/go
go run ./day01 < day01/testdata/example.txt
```

## Test

```bash
cd 2021/go
go test ./day01
```
