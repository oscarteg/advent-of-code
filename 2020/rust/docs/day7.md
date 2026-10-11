# Day 7: Handy Haversacks

## Part 1: Bag Containers

Parse bag containment rules and count how many bag colors can eventually contain at least one "shiny gold" bag.

Build a reverse index mapping each bag type to the bags that can contain it, then traverse upward from "shiny gold" using BFS to find all reachable containers.

**Complexity:** O(n + c) where n is the number of rules and c is the number of containment relationships.

## Part 2: Bag Contents

Count the total number of individual bags required inside a single "shiny gold" bag, accounting for the multiplicative nature of nested bags.

Recursively count the bags in each container, accumulating counts from nested bags multiplied by their required quantity.

**Complexity:** O(n) with memoization (worst case exponential without it on deeply nested structures).

## Run

```bash
cargo run -q --bin day7 < input/2020/day7.txt
```

## Test

```bash
cargo test --bin day7
```
