# Day 23: Crab Cups

## Part A
Parse the labels into a permutation of 1..n and hold the circle in a successor
array: `next[label]` is the label clockwise of it. A move splices the three cups
behind the current cup out of the chain, walks the destination label down
(wrapping at 1, skipping the picked-up labels) and splices them back in - four
writes, nothing shifts. After 100 moves the answer is the labels read clockwise
from cup 1. O(moves) time, O(n) memory.

## Part B
The same array sized to 1,000,000: input labels first, then 10..1,000,000 in
order, closed back onto the first cup. After ten million moves the answer is the
product of the two labels clockwise of cup 1. O(moves) time, O(total) memory;
every move stays O(1) because the destination walk skips at most three labels.

## Run and test
```bash
go run . 23 < ../rust/input/2020/day23.txt
go test -race ./challenge/day23
```
