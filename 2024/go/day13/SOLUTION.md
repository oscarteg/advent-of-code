# Day 13: Claw Contraption

## Approach

**Part 1:** Brute-force all combinations of A and B button presses (0-100 each) to find solutions where `a*ax + b*bx == px` and `a*ay + b*by == py`. Sum costs (3 tokens per A, 1 per B) for all winning machines.

**Part 2:** Add 10000000000000 to both prize coordinates, then use Cramer's rule to solve the 2x2 linear system directly (no bounds). Check if solutions are integers; if yes, add cost to total.

## Complexity

- Part 1: O(n * 100^2) = O(n) per machine, 10,000 iterations per machine
- Part 2: O(n) per machine, constant-time Cramer's rule solution

## Commands

```bash
go run ./day13 < day13/testdata/example.txt
go test ./day13
```
