# Day 7: Bridge Repair (C++)

Each line gives a target value and a list of operands. Operators (`+`, `*`,
and for part 2 `||` concatenation) are evaluated strictly left to right, with
no precedence; find which equations can reach their target with some
assignment of operators.

## Approach

Parse each line into a target and its operands. Recursively try every
operator at each position, pruning as soon as the running total exceeds the
target (all operators only grow or hold a positive value, so partial sums
never need to decrease). Part 1 allows `+`/`*`; part 2 adds `||`, which
concatenates the decimal digits of its two operands. Sum the targets of every
solvable equation.

## Complexity

With `n` operands per equation there are up to `2^(n-1)` (or `3^(n-1)` with
concatenation) operator assignments; the target-exceeded pruning keeps this
fast in practice for the given input sizes. Time: O(equations x branches).
Space: O(n) recursion depth per equation.

## Run / test

```bash
cd 2024/cpp
make test day=day7
make -s run day=day7 < examples/day7.txt
```
