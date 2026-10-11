# Day 8: Playground

Junction boxes hang in 3D space, one `x,y,z` per line. Pairs of boxes get
connected by strings of lights, closest straight-line distance first, and
every connection merges the circuits of the two boxes it joins.

Part 1 orders all `N*(N-1)/2` pairs by squared distance (the same order as
the straight-line distance, without a square root) and feeds the closest
ones into a union-find forest. A pair whose boxes already share a circuit
still counts as one of the connections, it just changes nothing. After the
prescribed number of connections - 1000 for the real input, 10 for the
20-box example, overridable through `solveWith` - the sizes of the three
largest circuits are multiplied.

Part 2 keeps walking the same ordered list. Ignoring the pairs that change
nothing makes this Kruskal's algorithm, so the connection that merges the
final two circuits is the last pair joined; part 2 multiplies its X
coordinates.

`O(N^2 log N)` time, dominated by sorting the pairs, and `O(N^2)` space.

From `2025/go`:

```sh
go run ./day08 < day08/testdata/example.txt
go test -race ./day08
go vet ./day08
```

The example prints exactly:

```text
Part 1: 40
Part 2: 25272
```
