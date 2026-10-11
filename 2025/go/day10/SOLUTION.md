# Day 10: Factory

One machine per line: indicator diagram in `[]`, button wirings in `()`,
joltage requirements in `{}`.

Part 1: toggling twice cancels out, so each button is pressed once or not at
all and the diagram is a linear system over GF(2). Elimination yields one
solution plus a null-space basis - button sets that toggle nothing - and the
answer is the lightest solution in that coset. `O(2^(b-r) * b)` for `b`
buttons of rank `r`.

Part 2: the lights are irrelevant: a press adds one to each counter its button
feeds, so the targets are `A x = joltage` over the nonnegative integers.
Fraction-free Gaussian elimination makes the pivot presses a function of the
free ones, each free press is bounded by the smallest target it feeds, and a
pruned depth-first search over that box gives the cheapest total.
`O(l * b^2)` elimination plus the search.

From `2025/go`:

```sh
go run ./day10 < day10/testdata/example.txt  # Part 1: 7, Part 2: 33
go test -race ./day10
go vet ./day10
```

Part 1's per-machine presses (2, 3, 2) come from the official puzzle text;
part 2's (10, 12, 11) from the task brief, whose statement needs a login.
