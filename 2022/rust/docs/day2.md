# Day 2: Rock Paper Scissors

## Part 1
Interpret the strategy guide literally: the second column is our move (X=Rock, Y=Paper, Z=Scissors). Score each round as the move value (1-3) plus outcome (0=loss, 3=draw, 6=win), then sum.

**Approach:** Parse pairs, determine outcome via move comparison, accumulate scores.
**Complexity:** O(n) where n is the number of lines.

## Part 2
Reinterpret the strategy guide: the second column is the desired outcome (X=lose, Y=draw, Z=win). Determine the move that achieves it, then score as before.

**Approach:** For each opponent move and desired outcome, select the winning move via a lookup table, then score.
**Complexity:** O(n) where n is the number of lines.

## Run
```bash
cargo run -q --bin day2 < input/day2-example.txt
```

## Test
```bash
cargo test --bin day2
```
