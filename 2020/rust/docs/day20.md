# Day 20: Jurassic Jigsaw

Reassemble square camera tiles that arrive shuffled, rotated and flipped.

## Part 1
Two tiles can only be joined where they share an edge, so index every edge by its
canonical form (the smaller of its two reading directions). An edge owned by a single
tile cannot be joined to anything and therefore lies on the border of the image, which
makes a tile with exactly two such edges a corner. Multiply the four corner ids.

## Part 2
Lay the tiles out row by row, starting from a corner turned so that its unjoinable
edges face up and left. Every later slot is pinned by an already placed neighbour: look
up the tiles owning that edge and keep the one orientation that joins it while leaving
the outer rows and columns unjoinable, so each placement is unique. Drop every tile's
border, stitch the interiors into one image and count sea monsters in each of its eight
orientations; the roughness is the lit pixels minus the 15 pixels of each monster found.

## Complexity
- **Time:** O(n) placements for n tiles, each resolved against 8 orientations of the few
  tiles owning the pinning edge, plus 8 scans of the p pixels of the stitched image.
- **Space:** O(p) for the stitched image.

## Run & Test
```bash
cd 2020/rust
cargo test --bin day20
cargo run -q --bin day20 < input/2020/day20-example.txt
```
