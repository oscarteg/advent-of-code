# Day 3: Binary Diagnostic

## Part 1: Power Consumption

Submarine diagnostic systems report a list of binary numbers. For each bit position, find the most common bit (0 or 1); this becomes the corresponding bit of the **gamma rate**. The **epsilon rate** is the opposite (least common bits). Power consumption = gamma × epsilon.

**Approach:** Count 1-bits at each position. If the count exceeds half the lines, that position's gamma bit is 1.

**Complexity:** O(n·b) where n = number of readings, b = bits per reading.

## Part 2: Life Support Rating

Filter the readings using two criteria:
- **Oxygen generator rating:** Keep readings with the most common bit at each position (ties favor 1), repeat until one remains.
- **CO₂ scrubber rating:** Keep readings with the least common bit at each position (ties favor 0), repeat until one remains.

Life support rating = oxygen × CO₂.

**Approach:** Iterative filtering: for each bit position, count 1-bits among candidates and filter based on the criterion.

**Complexity:** O(n·b) worst case, typically faster as candidates shrink.

## Run & Test

```bash
cd 2021/go
go test ./day03 -v
go run ./day03 < day03/testdata/example.txt
```

Expected output:
```
Part 1: 198
Part 2: 230
```
