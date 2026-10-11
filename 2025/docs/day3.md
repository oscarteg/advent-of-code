# Day 3: Lobby

## Problem

Each input line is a bank of digit batteries. From every bank, pick K digits
while preserving their order in the string, forming the largest possible
K-digit number; the answer is the sum over all banks.

**Part 1:** K = 2. **Part 2:** K = 12.

---

*Solution spoilers below.*

## Approach

Greedy per position. For each digit position in the result (left to right):
- Search the remaining input for the largest digit
- Ensure enough digits remain after it for the remaining positions
- Use that digit and continue from the position after it

## Complexity

- Time: O(n * k²) where n is total digits and k is target length (k=2 or k=12)
- Space: O(k) for the result

## Run

```bash
cd 2025
zig run src/day3.zig < src/examples/day3.txt
zig test src/day3.zig
```
