# Day 23: Crab Cups

## Approach

**Part 1:** Parse the labels into a permutation of `1..=n` and hold the circle in a successor
array: `next[label]` is the label clockwise of it. A move splices the three cups behind the
current cup out of the chain, walks the destination label down (wrapping at 1, skipping the
picked-up labels) and splices them back in - four writes, nothing shifts. After 100 moves the
answer is the labels read clockwise from cup 1.

**Part 2:** The same array sized to 1,000,000: input labels first, then `10..=1,000,000` in
order, closed back onto the first cup. After ten million moves the answer is the product of
the two labels clockwise of cup 1.

## Complexity

O(moves) time and O(total cups) memory; the destination walk skips at most three labels, so
every move is O(1).

## Commands

```bash
cd 2020/rust
cargo test --bin day23
cargo run -q --bin day23 < input/2020/day23-example.txt
```
