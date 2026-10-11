# Day 9: Rope Bridge

## Part 1
Simulate a 2-knot rope: move the head by one unit per step, then if the tail
is no longer adjacent to the head, move it one step toward the head (sign of
the delta per axis). Track visited tail positions in an `AutoHashMap`; the
answer is its count. O(n) in total steps.

## Part 2
Generalize to a 10-knot rope: each knot follows the knot in front of it with
the same adjacency rule, applied head-to-tail after every head move. Track
visited positions of the last knot only. O(n * k) for n steps and k knots.

## Run
```bash
zig run src/day_9.zig < src/examples/day_9.txt
```

## Test
```bash
zig test src/day_9.zig
```
