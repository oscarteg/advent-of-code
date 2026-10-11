# Day 1: Secret Entrance

## Problem

A dial with numbers 0-99 starts at position 50. Each input line specifies a rotation:
- `L<n>` rotates left by n steps
- `R<n>` rotates right by n steps

The dial wraps around (e.g., rotating left from 1 goes to 99, 98, etc.).

## Part 1

Count how many rotations end exactly on position 0.

**Example:** Starting at 50, the sequence L68 L30 R48 L5 R60 L55 L1 L99 R14 L82 results in 3 rotations ending on 0.

## Part 2

Count every click that lands on 0, including intermediate clicks during a rotation.

When rotating R1000 from position 50, the dial passes through 0 multiple times (every 100 clicks), so this counts 10 crossings even though it only makes 1 rotation.

**Example:** The same sequence results in 6 total landings on 0.

---

*Solution spoilers below.*

## Approach

Use modular arithmetic to track position on the dial. When rotating:
- Right: new_pos = (pos + distance) mod 100
- Left: new_pos = (pos - distance) mod 100

For crossings, use integer division to count how many complete wraps occur during each rotation. Special-case left rotation from 0 and account for partial wraps carefully.

## Complexity

- **Time:** O(n) where n is the number of rotations
- **Space:** O(1) (excluding input storage)

## Commands

```bash
cd 2025
zig fmt --check src/day1.zig
zig test src/day1.zig
zig run src/day1.zig < src/examples/day1.txt
```
