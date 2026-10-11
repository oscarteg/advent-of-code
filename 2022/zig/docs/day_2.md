# Day 2: Rock Paper Scissors

## Part 1
Interpret the strategy guide literally: the second column is our move (X=Rock, Y=Paper, Z=Scissors). Score each round as the move value (1-3) plus outcome (0=loss, 3=draw, 6=win), then sum.

**Approach:** Parse character pairs, determine outcome by checking move relationships, accumulate scores.
**Complexity:** O(n) where n is the number of lines.

## Part 2
Reinterpret the strategy guide: the second column is the desired outcome (X=lose, Y=draw, Z=win). Determine the move that achieves it, then score as before.

**Approach:** For each opponent move and desired outcome, select the correct move via a lookup pattern, then score.
**Complexity:** O(n) where n is the number of lines.

## Run
```bash
zig run src/day_2.zig < src/examples/day_2.txt
```

## Test
```bash
zig test src/day_2.zig
zig run src/day_2.zig < src/examples/day_2.txt
```
