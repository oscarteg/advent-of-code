# Day 15: Warehouse Woes

The input is a warehouse map (`#` wall, `O` box, `@` robot) and the robot's
moves, whose line breaks carry no meaning. Parsing rejects a ragged map, an
unknown tile, a map without exactly one robot, and one not enclosed by walls -
that wall border is what stops a push from running off the grid. Both parts
replay every move, then sum `100 * row + col` over the boxes.

## Part 1

Scan from the robot along the move to the first tile that is not a box: a wall
there blocks the push, otherwise the near box hops to the far end of the run,
which leaves the same map as shifting every box one tile. O(moves * width).

## Part 2

Tiles double in width, so a box becomes `[]` and can be half-aligned with its
neighbours. A sideways push shifts the run of halves tile by tile; a vertical
push collects the boxes breadth-first through both halves of each box, bails
out when any meets a wall, and moves them farthest first. O(moves * boxes).

## Run / test

```bash
go test -race -count=1 ./day15
go run ./day15 < day15/testdata/example.txt
```
