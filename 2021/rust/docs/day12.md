# Day 12: Passage Pathing

## Problem

Count paths through a cave system from `start` to `end`:
- Small caves (lowercase) can be visited at most once per path
- Large caves (uppercase) can be visited any number of times
- Part 1: Count all valid paths with standard rules
- Part 2: Count all valid paths where exactly one small cave can be visited twice

## Approach

**Parsing:** Parse edges as `cave1-cave2` and build an adjacency list (HashMap).

**Part 1:** Depth-first search with visited set tracking small caves. Backtrack when reaching `end` or a dead-end.

**Part 2:** DFS with a flag allowing one small cave to be visited twice. Track whether we've used this exception and prevent revisiting `start`.

## Complexity

- **Time:** O(2^E) worst-case where E is edges (exponential paths exist)
- **Space:** O(V + E) for graph storage, O(V) for recursion stack

## Commands

```bash
cd 2021/rust
cargo test --bin day12
cargo run -q --bin day12 < input/day12-example.txt
```
