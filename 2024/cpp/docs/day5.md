# Day 5: Print Queue (C++)

## Approach

**Part 1:** Parse ordering rules (X|Y means X must come before Y) and page updates. For each update, validate that it respects all ordering rules by checking that whenever both pages from a rule appear, they're in the correct order. Sum the middle page numbers of valid updates.

**Part 2:** For each invalid update, topologically sort the pages using Kahn's algorithm with the rules as a directed graph. Sum the middle page numbers of corrected updates.

The key insight is that the ordering rules form a partial order, and Kahn's algorithm efficiently finds a valid total order respecting those constraints.

## Complexity

- Parsing: O(R + U) where R is total rules, U is total updates
- Validation per update: O(R) 
- Topological sort per invalid update: O(P + R) where P is pages in update
- Overall: O(U × (R + P))

## Run / test

```bash
cd 2024/cpp
make test day=day5
make -s run day=day5 < examples/day5.txt
```
