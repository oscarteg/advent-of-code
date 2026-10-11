# Day 13: Claw Contraption (Zig)

## Approach

**Part 1:** Each claw machine has buttons A (costs 3 tokens) and B (costs 1 token). Brute-force all combinations of A and B presses from 0–100 to reach the prize location. The goal is to find non-negative integers `a` and `b` such that:
- `a*ax + b*bx == px`
- `a*ay + b*by == py`

Sum the token costs for all winnable machines.

**Part 2:** Add 10,000,000,000,000 to both prize coordinates, then solve the same 2x2 linear system using Cramer's rule. No press bounds apply—solve directly via:
- `det = ax*by - ay*bx`
- `a = (px*by - py*bx) / det`
- `b = (ax*py - ay*px) / det`

Check if `a` and `b` are non-negative integers; if yes, add the cost to the total.

## Complexity

- Part 1: O(n * 100²) per puzzle input (10,000 iterations per machine)
- Part 2: O(n) per puzzle input (constant-time Cramer's rule)

## Run / test

```bash
cd 2024/zig
zig test day13.zig
zig run day13.zig < examples/day13.txt
```
