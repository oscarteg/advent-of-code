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

Cards and jokers are enums with derives for comparison. `classify_hand` and `classify_hand_joker` compute hand types by counting card frequencies. Hands sort naturally via Ord trait.

## Complexity

- Time: O(n log n) for sort
- Space: O(n) for hands

## Commands

```bash
cargo test day07::
cargo clippy --all-targets -- -D warnings
cargo fmt -- --check
cargo run -q --bin day07 < examples/day07.txt
```
