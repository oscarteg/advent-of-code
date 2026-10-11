# Day 1: Secret Entrance

## Problem

A dial with numbers 0-99 starts at position 50. Each input line specifies a rotation:
- `L<n>` rotates left by n steps
- `R<n>` rotates right by n steps

The dial wraps around (e.g., rotating left from 1 goes to 99, 98, etc.).

## Part 1

Count how many rotations end exactly on position 0.

## Part 2

Count every click that lands on 0, including intermediate clicks during a rotation. When rotating across the 0 boundary multiple times, count each crossing.

## Approach

Use modular arithmetic to track position on the dial:
- Right: new_pos = (pos + distance) % 100, crossings = distance / 100
- Left: new_pos = (pos - distance) mod 100, with special handling for boundaries

For left rotation, carefully count crossings based on whether the distance passes through 0:
- If distance < current position: no crossing
- If distance >= current position: at least one crossing, plus additional full wraps

## Complexity

- **Time:** O(n) where n is the number of rotations
- **Space:** O(1) (excluding input storage)

## Commands

```bash
cd 2025/go
gofmt -l day01
go vet ./day01
go test -race -count=1 ./day01
go run ./day01 < day01/testdata/example.txt
```
