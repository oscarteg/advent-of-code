# Day 12: Passage Pathing

## Problem

Count paths through a cave system from `start` to `end`:
- Small caves (lowercase) can be visited at most once per path
- Large caves (uppercase) can be visited any number of times
- Part 1: Count all valid paths with standard rules
- Part 2: Count all valid paths where exactly one small cave can be visited twice

## Approach

**Parsing:** Parse edges as `cave1-cave2` and build an adjacency list (map).

**Part 1:** Depth-first search with visited map tracking small caves. Recurse when reaching neighbors, backtrack naturally.

**Part 2:** DFS with a flag allowing one small cave to be visited twice. Track whether we've used this exception and prevent revisiting `start`.

## Complexity

- **Time:** O(2^E) worst-case where E is edges (exponential paths exist)
- **Space:** O(V + E) for graph storage, O(V) for recursion stack and visited copies

## Commands

```bash
cd 2021/go
go test ./day12
go run ./day12 < day12/testdata/example.txt
```
