# Day 2: Gift Shop

## Problem

The input is one line of comma-separated inclusive ID ranges `a-b`.

**Part 1:** an ID qualifies if its decimal digits are some sequence repeated
exactly twice (`55`, `6464`, `123123`). Sum every qualifying ID in the ranges.

**Part 2:** an ID qualifies if its digits are a shorter sequence repeated two
or more times (so `123123123` and `1111111` also count). Sum them.

---

*Solution spoilers below.*

## Approach

Check each ID's digit string. Part 1: even length, and the first half equals
the second half. Part 2: for every candidate period length that evenly divides
the digit count, check whether that period tiles the whole string.

## Complexity

Let `N` be the number of IDs spanned by all ranges and `D` the digit count
of the largest ID. Checking one ID is `O(D^2)` in the worst case (one scan
per candidate period), so the whole solve is `O(N * D^2)`.

## Run / test

```bash
cd 2025
zig test src/day2.zig
zig run src/day2.zig < src/examples/day2.txt
```
