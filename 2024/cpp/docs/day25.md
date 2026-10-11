# Day 25: Code Chronicle (C++)

Single part: count lock/key pairs whose pins cannot overlap.

## Approach

Each 7x5 schematic is a lock (top row full of `#`) or a key (bottom row full).
For each column, count the `#` characters across all 7 rows minus 1 (the full
row) to get its pin height, 0-5. A lock and key fit if, for every column, the
sum of their heights does not exceed 5 (the available pin space between the
two full rows). Count all fitting lock/key pairs by brute force.

## Complexity

Parsing is O(n) in input size. Checking fit is O(L * K * 5) for L locks and
K keys; with small inputs this is effectively O(L * K).

## Run / test

```bash
cd 2024/cpp
make test day=day25
make -s run day=day25 < examples/day25.txt
```
