# Day 23: LAN Party (C++)

## Problem

Given a network of computer connections, find:
1. **Part 1**: Count all sets of 3 inter-connected computers where at least one has a name starting with 't'.
2. **Part 2**: Find the largest fully-connected subset of computers (maximum clique).

## Approach

- **Graph representation**: Adjacency lists using hash maps keyed by computer names.
- **Part 1**: Find all triangles (3-node cliques) by checking each node and its neighbor pairs. Filter for triangles containing a 't' node.
- **Part 2**: Use the Bron-Kerbosch algorithm to find all maximal cliques, then select the largest. For smaller inputs, enumerate subsets directly.

## Complexity

- **Time**: O(n³) for triangle enumeration; Bron-Kerbosch is exponential in worst case but efficient on sparse graphs.
- **Space**: O(n + m) for graph storage.

## Run / test

```bash
cd 2024/cpp
make test day=day23
make -s run day=day23 < examples/day23.txt
```
