# Day 9: Rope Bridge

## Part 1
Simulate a 2-knot rope: move the head by one unit per step, then if the tail
is no longer adjacent to the head, move it one step diagonally/straight
toward the head (`signum` of the delta per axis). Track every position the
tail visits in a `HashSet`; the answer is its size. O(n) in total steps.

## Part 2
Generalize to a rope of 10 knots: each knot follows the knot in front of it
using the same adjacency rule, applied in order from head to tail after every
head move. Track visited positions of the last knot only. O(n * k) for n
steps and k knots.

## Run
```bash
cargo run -q --bin day9 < input/day9-example.txt
```

## Test
```bash
cargo test --bin day9
```
