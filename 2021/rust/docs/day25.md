# Day 25: Sea Cucumber

## Part 1: Find the first step where sea cucumbers stop moving

Parse a grid containing sea cucumbers (`>` facing east, `v` facing south).
Simulate movement in two phases each step:

1. All east-facing cucumbers move right one cell if empty (wrapping around)
2. All south-facing cucumbers move down one cell if empty (wrapping around)

Continue until no cucumbers move. Return the step number.

## Complexity

**Time:** O(steps × h × w) where steps ≈ 400 (typical), h is height, w is width  
**Space:** O(h × w) for the grid copy each step

## Commands

```bash
cd 2021/rust
cargo test --bin day25
cargo run -q --bin day25 < input/day25-example.txt
```
