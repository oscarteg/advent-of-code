# Day 22: Crab Combat

## Part A
Simulate Combat: both players draw their top card each round; the higher
card wins and takes both cards (loser's, then winner's) onto the bottom of
the deck. Continue until one deck is empty, then score it by
`card * position_from_bottom`. O(rounds), bounded by the total card count.

## Part B
Recursive Combat: if both players hold at least as many cards as their
drawn card's value, the round is decided by a sub-game on a deck prefix
instead of a plain comparison. Every `(deck1, deck2)` state seen earlier in
a game makes player 1 win that game instantly, guaranteeing termination.
Exponential worst case, fast in practice for 25-50 card decks.

## Commands
```bash
go test -race ./challenge/day22
go run . 22 < ../rust/input/2020/day22.txt
```
