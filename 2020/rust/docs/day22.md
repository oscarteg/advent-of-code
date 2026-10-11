# Day 22: Crab Combat

## Part 1
Simulate Combat: each round both players draw their top card, the higher
card wins and takes both cards (loser's first, winner's second) onto the
bottom of their deck. Play continues until one deck is empty; score the
winning deck by summing `card * position_from_bottom`.
Complexity: O(rounds), bounded by the finite total card count.

## Part 2
Recursive Combat adds a rule: if both players have enough cards left to
match their drawn card's value, the round winner is decided by a sub-game
played on a prefix of each deck instead of a plain comparison. To
guarantee termination, every `(deck1, deck2)` state seen earlier in a game
causes player 1 to win that game instantly.
Complexity: exponential in the worst case (bounded by the state space of
reachable deck configurations), fast in practice for 25-50 card decks.

## Commands
```bash
cargo test --bin day22
cargo run -q --bin day22 < input/2020/day22-example.txt
```
