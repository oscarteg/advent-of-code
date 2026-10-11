# Day 5: Print Queue

## Approach

**Part 1:** Parse ordering rules and page updates. For each update, check if it's in valid order by verifying that for every rule (a before b), a appears before b in the update. Sum the middle page numbers of valid updates.

**Part 2:** For invalid updates, sort them topologically using Kahn's algorithm with the ordering rules as a directed graph. Sum the middle page numbers of corrected updates.

## Complexity

- Parsing: O(R + U) where R is rules and U is updates
- Validation: O(U × R) per update
- Topological sort: O(P + R) per invalid update (Kahn's algorithm)
- Overall: O(U × (R + P))

## Commands

```bash
cd 2024/go && go test -race ./day05 && go run ./day05 < day05/testdata/example.txt
```
