# Day 5: Hydrothermal Venture

## Algorithm

The solution tracks overlapping lines of hydrothermal vents by counting how many lines pass through each point.

1. **Parse input**: Convert each line from `x1,y1 -> x2,y2` format into start and end points.
2. **Generate points**: For each line, compute all integer points it covers by stepping from start to end with unit increments in x and/or y (horizontal, vertical, or diagonal).
3. **Count overlaps (Part 1)**: Filter to horizontal/vertical lines only (where x₁=x₂ or y₁=y₂), track point frequencies in a map, count points with frequency ≥ 2.
4. **Count overlaps (Part 2)**: Include all lines (adding 45° diagonals), use the same counting logic.

## Complexity

- **Time**: O(n·L) where n is the number of lines and L is the average line length. Line length is bounded by max coordinate value (~1000 in typical inputs).
- **Space**: O(P) where P is the total number of unique points covered across all lines.

## Run & Test

```bash
cd /Users/oscar/Developer/Personal/advent-of-code/2021/go

# Run tests with race detection
go test -race ./day05

# Run diagnostics
go vet ./day05

# Run with sample input
go run ./day05 < day05/testdata/example.txt
```

Expected output:
```
Part 1: 5
Part 2: 12
```

## Implementation Notes

- Lines are represented as `struct { start, end point }`; points as `struct { x, y int }`
- Point iteration uses a direction vector (dx, dy ∈ {-1, 0, 1}) to handle all 8 directions
- Overlap counting uses a `map[point]int` for O(1) lookups per point
- Input parsing validates format and handles whitespace robustly
- Error cases include invalid line format, missing coordinates, and empty input
