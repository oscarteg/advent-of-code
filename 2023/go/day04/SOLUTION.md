# Day 4: Scratchcards

Each card has a winning-number set and a "your numbers" list; `matches` is the
count of your numbers also in the winning set.

- **Part 1:** A card with `n` matches scores `2^(n-1)` (0 for no matches); sum
  over all cards. `O(cards * numbers)`.
- **Part 2:** A card's matches win one extra copy of each of the following `n`
  cards. A single forward pass keeps a running copy count per card index
  (starting at 1) and propagates it; sum the counts. `O(cards * numbers)`.

From `2023/go`:

```sh
go run ./day04 < day04/testdata/example.txt # Part 1: 13, Part 2: 30
go test ./day04
```

Use the same commands with your puzzle input redirected into stdin.
