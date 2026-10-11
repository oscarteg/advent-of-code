# Day 23: LAN Party

## Problem

Find sets of 3 inter-connected computers and the maximum fully-connected subset.

## Approach

- **Part 1**: Enumerate all triangles (3-cliques) in the graph; count those with at least one 't' node.
- **Part 2**: Search for the maximum clique by checking all 4-node and 3-node subsets.

## Complexity

- **Time**: O(n³) for triangle enumeration; O(n⁴) for clique search on small inputs.
- **Space**: O(n + m) for graph storage.

## Run & Test

`go run ./day23 < testdata/example.txt`
`go test ./day23`
