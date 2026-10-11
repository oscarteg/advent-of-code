# Day 7: Camel Cards

## Part 1: Hand Rankings

Parse poker-like hands and rank by strength. Hand types (strongest to weakest):
- Five of a kind
- Four of a kind
- Full house
- Three of a kind
- Two pair
- One pair
- High card

Within same type, compare cards left-to-right (A > K > Q > J > T > 9..2).

Calculate winnings: sum of (bid × rank) for each hand.

## Part 2: Joker Wild

J becomes a joker (wildcard, lowest value). Adjust hand classification to maximize type by treating jokers optimally. J now ranks lowest (before 2).

## Implementation

Card and Joker types use iota constants. Hand classification counts card frequencies and maps to types. Custom Hands and HandsJoker types implement sort.Interface for natural sorting by type then cards.

## Complexity

- Time: O(n log n) for sort
- Space: O(n) for hands

## Commands

```bash
go test -race ./day07
go vet ./day07
go run ./day07 < day07/testdata/example.txt
```
