# Day 4: Scratchcards

- Shared parsing in `src/scratchcards.rs`: each card is a winning-number set and a
  "your numbers" list; `num_matches` is their intersection size.
- **Part 1:** A card with `n` matches scores `2^(n-1)` (0 for no matches); sum over
  all cards. `O(cards * numbers)`.
- **Part 2:** A card's matches win one extra copy of each of the following `n`
  cards. Track a running copy count per card index (starting at 1) and propagate
  forward in a single pass; sum the counts. `O(cards * numbers)`.

From `2023/rust`:

```sh
cargo run -q --bin day04 < examples/day04.txt # 13
cargo run -q --bin day04 < examples/day04.txt # 30
cargo test day04::
```

Use the same commands with your puzzle input redirected into stdin.
