# Day 8: Seven Segment Search

## Part 1
Count digits 1, 4, 7, and 8 in the output values. These are the "easy" digits because they use unique numbers of segments:
- 1 uses 2 segments
- 4 uses 4 segments
- 7 uses 3 segments
- 8 uses 7 segments

**Approach:** Scan output digits and count those with lengths 2, 3, 4, or 7.

**Complexity:** O(n) where n is the total number of output digits.

## Part 2
Deduce the full mapping of signal patterns to digits by exploiting set-containment relationships.

**Approach:**
1. Identify easy digits (1, 4, 7, 8) by segment count
2. Use segment containment to deduce 6-segment digits (0, 6, 9):
   - 9 contains all segments of 4
   - 0 contains all segments of 1 but not 4
   - 6 lacks all segments of 1
3. Deduce 5-segment digits (2, 3, 5):
   - 3 contains all segments of 1
   - 5 is contained entirely in 6
   - 2 is the remaining one
4. Decode output value using the mapping

**Complexity:** O(n) where n is the number of entries (constant-factor segment comparisons per entry).

## Commands

```bash
go test ./day08
go run ./day08 < day08/testdata/example.txt
```
