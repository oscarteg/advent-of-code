# Day 8: Haunted Wasteland

Parse the instruction line and the `name = (left, right)` nodes into a map.

- **Part 1:** Walk from `AAA`, repeating the instructions, until `ZZZ` is
  reached, and count the steps.
- **Part 2:** Walk from every node ending in `A` to the first node ending in
  `Z`. Each ghost revisits its destination once per walk length, so the answer
  is the least common multiple of those lengths, combined with a checked
  multiplication.
- A walk costs `O(steps)` time and the network `O(nodes)` memory. Walking
  longer than `nodes * instructions` steps means the (node, instruction) state
  repeated, so the destination is unreachable and the solver returns an error
  instead of looping forever. Malformed instructions, node lines, unknown
  targets and duplicate nodes are rejected while parsing.

From `2023/go`:

```sh
go run ./day08 < day08/testdata/example.txt # Part 1: 2, Part 2: 2
go test -race ./day08
go vet ./day08
```
