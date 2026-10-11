# Day 19: Beacon Scanner

Every report lists beacons in its own scanner's frame: rotated by one of the 24
axis-aligned rotations and shifted by the unknown scanner position.

Part 1 fixes scanner 0 at the origin and grows that frame breadth-first. For each
rotation of an unplaced report, count how often every `known - rotated` offset occurs;
an offset hit twelve times is the overlap, so it is the scanner position and the
rotated, shifted beacons join the global set. Part 1 is the size of that set.

Part 2 is the largest Manhattan distance over all pairs of recovered positions.

For `S` reports of `B` beacons, matching one pair costs `O(24 * B^2)`, so assembling
all of them is `O(S^2 * B^2)` in the worst case; space is `O(S * B)`.

From `2021/go`:

```sh
go run ./day19 < day19/testdata/example.txt
# Part 1: 79
# Part 2: 3621
go test -race ./day19
go vet ./day19
```
