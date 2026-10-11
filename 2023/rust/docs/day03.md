# Day 3: Gear Ratios

## Part 1
Parse the engine schematic grid and find all numbers adjacent (including diagonally) to any non-period symbol. Sum these part numbers.

**Approach:** Build a map of (x, y) positions to values (number, symbol, empty). For each symbol, check its 8 adjacent cells for numbers. Track starting positions to avoid counting the same multi-digit number twice.

**Complexity:** O(n×m) where n, m are grid dimensions.

## Part 2
Find all `*` symbols (gears) that have exactly 2 adjacent numbers. Multiply those numbers and sum all gear ratios.

**Approach:** Same grid construction, but only process `*` symbols and collect adjacent numbers into a vector. If exactly 2 numbers exist, multiply them and add to sum.

**Complexity:** O(n×m).

## Run
```bash
cd 2023/rust
cargo run -q --bin day03 < examples/day03.txt
```

## Test
```bash
cd 2023/rust
cargo test day03::
```
