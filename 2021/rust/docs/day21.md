# Day 21: Dirac Dice

Part 1 simulates turns with a deterministic die (1..100, cycling), moving each
player's pawn around the 10-space track and summing scores until one reaches
1000. Runtime is `O(rolls)`, a few hundred turns.

Part 2 explores every quantum universe. Each turn's three 3-sided die rolls
collapse to 7 possible sums (3..9) with known multiplicities, so the state
space is `(position, score)` for both players and whose turn it is. Memoizing
`count_wins(current, other)` on that state makes each state resolve in O(1)
after its first visit; with 10 positions and 21 scores per player the memo
has at most `10 * 21 * 10 * 21` entries, so overall time and space are O(1)
relative to score limit and board size.

Run and test from `2021/rust`:

```sh
cargo run -q --bin day21 < input/day21-example.txt
cargo test --bin day21
```

The example prints:

```text
Part 1: 739785
Part 2: 444356092776315
```
