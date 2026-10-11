# Day 12: Rain Risk

## Part A

Track ship position `(x, y)` and a heading vector, starting facing east.
`N`/`S`/`E`/`W` move the ship directly; `L`/`R` rotate the heading vector by
a multiple of 90 degrees via `(x, y) -> (-y, x)` applied `degrees/90` times;
`F` moves the ship by the heading vector. Answer is the Manhattan distance
from the origin. O(n) time, O(1) extra space.

## Part B

Same instructions, but `N`/`S`/`E`/`W`/`L`/`R` move/rotate a waypoint
relative to the ship; `F` moves the ship by the waypoint offset, repeated by
the given value. Same rotation helper, applied to the waypoint vector.
O(n) time, O(1) extra space.

## Commands

```bash
go test -race -count=1 ./challenge/day12
go run . 12 < ../rust/input/2020/day12.txt
```
