# Day 1: Sonar Sweep

## Problem

Given a list of sonar depth measurements, count increases.

**Part 1:** Count how many times a depth measurement increases from the previous measurement.

**Part 2:** Count how many times the sum of a 3-measurement sliding window increases from the previous window's sum. (Equivalently: count how many times `depths[i] < depths[i+3]`.)

## Approach

**Part 1:** Use `windows(2)` to compare consecutive pairs.

**Part 2:** Use `windows(4)` to compare elements 3 positions apart, which is equivalent to comparing sliding window sums without computing the sums explicitly.

## Complexity

- Time: O(n) for both parts
- Space: O(1) (excluding input storage)

## Run

```bash
cd 2021/rust
cargo run -q --bin day1 < input/day1-example.txt
```

## Test

```bash
cd 2021/rust
cargo test --bin day1
```
