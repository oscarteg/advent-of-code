# Day 19: Linen Layout (C++)

## Approach

**Part 1:** Count how many target designs can be constructed from available towel patterns.

**Part 2:** Sum the total number of ways to construct all designs.

Both parts use recursive dynamic programming with memoization:
- For each target design, try to match each available pattern as a prefix
- If a pattern matches, recursively solve the remaining substring
- Memoize results to avoid recomputation
- Part 1 counts designs with at least one way; Part 2 sums all ways

## Complexity

- Time: O(D × |R| × P × M) where D = designs, R = remaining string lengths, P = patterns, M = memoization lookups
- Space: O(D × |R|) for memoization across all designs

## Run / test

```bash
cd 2024/cpp
make test day=day19
make -s run day=day19 < examples/day19.txt
```
