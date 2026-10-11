# Day 4: Printing Department

## Problem

A grid of paper rolls (`@`) and empty space (`.`). A roll is accessible to a
forklift when fewer than 4 of its 8 neighbours are also rolls.

**Part 1:** count the accessible rolls in the original grid.

**Part 2:** repeatedly remove every currently accessible roll, simultaneously
per round (removals within a round do not affect each other), until no roll
is accessible. The total removed across all rounds is the answer.

---

*Solution spoilers below.*

## Approach

Parsing rejects rows of inconsistent width, characters other than `@`/`.`,
and empty input. Part 2 recomputes accessibility against a snapshot of the
round's grid and applies all removals at once.

With an `R`-row by `C`-column grid, each round is `O(R*C)` and the number of
rounds is bounded by the grid diameter, so part 2 is `O(R*C*rounds)`. Space
is `O(R*C)` for the working copy.

## Run / test

From `2025`:

```sh
zig test src/day4.zig
zig run src/day4.zig < src/examples/day4.txt
```

The example prints exactly:

```text
Part 1: 13
Part 2: 43
```
