# Day 12: Garden Groups

Orthogonally touching plots of the same plant letter form a region.

## Part 1

Flood-fill every region from the first plot not yet visited. Each step that
leaves the region (different letter or grid edge) adds one to its perimeter;
the price is area * perimeter. Complexity: O(width * height).

## Part 2

A region has as many sides as corners, so the same fill counts corners per
plot: for each of the four diagonals, both adjacent plots outside the region
is a convex corner, and both inside while the diagonal plot is outside is a
concave corner. The price is area * sides. Complexity: O(width * height).

## Run / test

```bash
go test -race -count=1 ./day12
go run ./day12 < day12/testdata/example.txt
```
