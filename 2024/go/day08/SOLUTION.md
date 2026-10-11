# Day 8: Resonant Collinearity

## Approach

**Part 1:** Group antennas by frequency. For each pair of same-frequency antennas at positions p1 and p2, compute antinodes at p1 + (p1 - p2) and p2 + (p2 - p1). Keep only those within bounds.

**Part 2:** Extend to all collinear points. For each antenna pair, normalize the direction vector using GCD to get unit steps, then walk the line in both directions from both antennas, adding all points within bounds.

## Complexity

- Time: O(n² × m) where n is antennas per frequency, m is grid dimension (for line walking)
- Space: O(m²) for storing antinodes

## Commands

```bash
test -z "$(gofmt -l day08)" && go vet ./day08 && go test -race -count=1 ./day08 && go run ./day08 < day08/testdata/example.txt
```

## Example Answers

Part 1: 14
Part 2: 34
