# Day 25: Snowverload

## Part 1

The wires form an undirected graph whose global minimum cut is the three wires
to pull. Every wire becomes two unit-capacity arcs (`arc` and `arc^1`), then
component 0 is fixed as the source and Edmonds-Karp runs towards every other
component. The cheapest of those maximum flows is the global minimum cut,
because the source lies on one of the cut's two sides. Whatever stays reachable
from the source in the residual graph is one group, the rest is the other;
multiplying both sizes gives the answer.

A candidate sink is abandoned as soon as its flow can no longer beat the
cheapest cut so far, so it costs `O(cut * wires)` with `cut == 3` for puzzle
input: `O(components * cut * wires)` in total, `O(components + wires)` memory.

Single part puzzle: part 2 is just pushing the button.

## Run & test

```bash
go test -race -count=1 ./day25
go run ./day25 < day25/testdata/example.txt
```

The example prints exactly:

```text
Part 1: 54
```
