# Day 22: Sand Slabs

Each line parses into an axis-aligned brick with inclusive bounds per axis.
Sorting by lowest `z` and dropping the bricks in that order settles the stack: a
height map keyed by `(x, y)` column holds the highest occupied `z` plus the
brick owning it, so a falling brick lands one cube above the tallest column
under its footprint and the bricks owning those top cubes become its supporters.
That yields `supports` and `supportedBy` for the settled stack.

Part 1 counts bricks where every brick resting on them has a second supporter.
Part 2 walks up the support graph per brick with a queue, marking a brick as
falling once all of its supporters fall, and sums those chain reactions.

Settling costs `O(n log n)` plus the total footprint area; part 2 traverses the
graph once per brick, `O(n * edges)`.

From `2023/go`:

```sh
go run ./day22 < day22/testdata/example.txt
go test -race ./day22
go vet ./day22
```

The example prints exactly:

```text
Part 1: 5
Part 2: 7
```
