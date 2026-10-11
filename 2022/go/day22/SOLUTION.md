# Day 22: Monkey Map

Parse the notes into empty/open/wall tiles plus a path of forward counts and
turns, then walk it; only the wrapping rule differs between the parts.

**Part 1** stays flat: keep stepping with floor-mod over the bounding box until a
non-empty tile shows up, keeping the facing.

**Part 2** folds the net generically. Face size is `sqrt(tiles / 6)`; a BFS over
the net blocks rotates a 3D basis (`right`, `down`, `normal`) a quarter turn
across each shared edge. Cell `(i, j)` of a face is the doubled lattice point
`P = size*normal + (2i-size+1)*right + (2j-size+1)*down`; leaving it in 3D
direction `t` lands at `P + t - normal` on the face whose normal is `t`, now
heading along `-normal`. Any of the 11 cube nets folds this way.

Time `O(tiles + path length)`, space `O(tiles)`.

From `2022/go`:

```sh
go test -race ./day22
go run ./day22 < day22/testdata/example.txt   # Part 1: 6032, Part 2: 5031
```
