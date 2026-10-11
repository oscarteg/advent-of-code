# Day 4: Printing Department

A grid of paper rolls (`@`) and empty space (`.`). A roll is accessible to a
forklift when fewer than 4 of its 8 neighbours are also rolls.

**Part 1** counts accessible rolls in the original grid.

**Part 2** repeatedly removes every currently accessible roll, simultaneously
per round, until none remain, and totals the removals across all rounds.

Parsing rejects rows of inconsistent width, characters other than `@`/`.`,
and empty input.

With an `R`-row by `C`-column grid, each round is `O(R*C)` and the number of
rounds is bounded by the grid diameter, so part 2 is `O(R*C*rounds)`. Space
is `O(R*C)` for the working copy.

From `2025/go`:

```sh
go run ./day04 < day04/testdata/example.txt
go test -race ./day04
go vet ./day04
```

The example prints exactly:

```text
Part 1: 13
Part 2: 43
```
