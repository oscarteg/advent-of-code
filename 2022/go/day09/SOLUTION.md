# Day 9: Rope Bridge

## Part 1
Simulate a 2-knot rope: move the head by one unit per step, then if the tail
is no longer adjacent to the head, move it one step toward the head (sign of
the delta per axis). Record every position the tail visits in a set; the
answer is the set's size. O(n) in total steps.

## Part 2
Generalize to a 10-knot rope: each knot follows the knot in front of it with
the same adjacency rule, applied head-to-tail after every head move. Track
visited positions of the last knot only. O(n * k) for n steps and k knots.

## Run
```bash
go run ./day09 < day09/testdata/example.txt
```

## Test
```bash
go test -race -count=1 ./day09
```
