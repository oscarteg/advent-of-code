# Day 8: Playground

## Problem

Junction boxes hang in 3D space, one `x,y,z` per line. Pairs of boxes get
connected by strings of lights, closest straight-line distance first, and
every connection merges the circuits of the two boxes it joins. A pair whose
boxes already share a circuit still counts as a connection, it just changes
nothing.

**Part 1:** after the prescribed number of connections (1000 for the real
input, 10 for the 20-box example), multiply the sizes of the three largest
circuits.

**Part 2:** keep connecting until every box is in one circuit; multiply the
X coordinates of the pair whose connection merges the final two circuits.

---

*Solution spoilers below.*

## Approach

Order all `N*(N-1)/2` pairs by squared distance - the same order as the
straight-line distance, without a square root - and feed the closest ones
into a union-find forest. Ignoring the pairs that change nothing makes
part 2 Kruskal's algorithm, so the merging connection is the last pair joined.

## Complexity

`O(N^2 log N)` time, dominated by sorting the pairs, and `O(N^2)` space.

## Run / test

```sh
cd 2025
zig test src/day8.zig
zig run src/day8.zig < src/examples/day8.txt
```

The example prints exactly:

```text
Part 1: 40
Part 2: 25272
```
