# Day 21: Dirac Dice

Part 1 simulates turns with a deterministic die (1..100, cycling), moving each
player's pawn around the 10-space track and summing scores until one reaches
1000. Runtime is `O(rolls)`, a few hundred turns.

Part 2 explores every quantum universe. Each turn's three 3-sided die rolls
collapse to 7 possible sums (3..9) with known multiplicities, so the state
space is `(position, score)` for both players and whose turn it is. Memoizing
`countWins(current, other)` on that state resolves each state in O(1) after
its first visit; with 10 positions and 21 scores per player the memo has at
most `10 * 21 * 10 * 21` entries, so overall time and space are bounded by
the score limit and board size, independent of input.

Run and test from `2021/go`:

```sh
go run ./day21 < day21/testdata/example.txt
go test -race ./day21
go vet ./day21
```

The example prints:

```text
Part 1: 739785
Part 2: 444356092776315
```
