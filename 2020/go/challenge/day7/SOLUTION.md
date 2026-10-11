# Day 7: Handy Haversacks

## Part A: Bag Containers

Parse bag containment rules and count how many bag colors can eventually contain at least one "shiny gold" bag.

Build a reverse index mapping each bag type to the set of bags that can contain it, then traverse upward from "shiny gold" using BFS to find all reachable containers.

**Complexity:** O(n + c) where n is the number of rules and c is the number of containment relationships.

## Part B: Bag Contents

Count the total number of individual bags required inside a single "shiny gold" bag, accounting for the multiplicative nature of nested bags.

Recursively count the bags in each container with memoization to handle repeated subproblems, accumulating counts from nested bags multiplied by their required quantity.

**Complexity:** O(n) with memoization.

## Run

```bash
go run . 7 < ../rust/input/2020/day7.txt
```

## Test

```bash
go test -race ./challenge/day7
```
