# Day 21: Monkey Math

Each line parses into a named job: a literal number, or an operator on two
other monkeys.

Part 1 evaluates `root` recursively. A budget of one recursion step per monkey
turns a self-referential riddle into an error instead of a hang.

Part 2 reads `root` as an equality with `humn` as the unknown. Every monkey is
mapped to the monkey listening to it, which gives the path from `humn` up to
`root`. Walking that path back down, the branch without `humn` always evaluates
to a known number, so each operation is inverted (`a + b = t` gives `a = t - b`,
`a / b = t` gives `a = t * b`, ...) until it pins down `humn`. Divisions have to
be exact, otherwise the riddle has no integer answer.

Both parts are `O(n)` in the number of monkeys for a tree-shaped riddle, and
space is `O(n)`: every subtree hanging off the path is evaluated once.

From `2022/go`:

```sh
go run ./day21 < day21/testdata/example.txt   # Part 1: 152, Part 2: 301
go test -race ./day21
go vet ./day21
```
