# Day 6: Tuning Trouble

## Problem

Find the start-of-packet and start-of-message markers in a datastream.

- **Part 1:** Find the position where the first sequence of 4 distinct characters ends.
- **Part 2:** Find the position where the first sequence of 14 distinct characters ends.

## Approach

Both parts use a sliding window with frequency tracking:
1. Maintain a window of the required size (4 for part 1, 14 for part 2)
2. Track character frequencies in the current window using an array
3. Count unique characters in the window
4. Return the position when the count equals window size

A frequency array provides O(1) insertion, deletion and lookup for all 256 possible byte values.

## Complexity

- Time: O(n) where n is the input length
- Space: O(1) constant space for the frequency array

## Commands

```bash
cd 2022/zig
zig test src/day_6.zig
zig run src/day_6.zig < src/examples/day_6.txt
```
