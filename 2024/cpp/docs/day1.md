# Day 1: Historian Hysteria (C++)

## Part 1: Distance Calculation

Sort both lists and sum the absolute differences between corresponding elements.

## Part 2: Similarity Score

Count occurrences of each left-list element in the right list, then sum `element * count`.

## Complexity

- Time: O(n log n) for sorting, O(n) for similarity calculation
- Space: O(n) for storing the lists and hash map

## Run / test

```bash
cd 2024/cpp
make test day=day1
make -s run day=day1 < examples/day1.txt
```
