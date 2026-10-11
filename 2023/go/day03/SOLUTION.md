# Day 3: Gear Ratios

## Part 1
Parse the engine schematic grid and find all numbers adjacent (including diagonally) to any non-period symbol. Sum these part numbers.

**Approach:** Build a map of positions to values. For each symbol, check all 8 adjacent cells for numbers. Track processed positions to avoid counting the same multi-digit number twice.

**Complexity:** O(n×m) where n, m are grid dimensions.

## Part 2
Find all `*` symbols (gears) that have exactly 2 adjacent numbers. Multiply those numbers and sum all gear ratios.

**Approach:** Same grid construction, but only process `*` symbols and collect adjacent numbers. If exactly 2 numbers exist, multiply them and add to sum.

**Complexity:** O(n×m).

## Run
```bash
cd 2023/go
go run ./day03 < day03/testdata/example.txt
```

## Test
```bash
cd 2023/go
go test ./day03
```
