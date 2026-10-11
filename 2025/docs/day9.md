# Day 9: Movie Theater

## Problem

Red tiles are listed as `x,y` in the order they are walked around the room.

**Part 1:** take any two red tiles as opposite corners of a rectangle whose
tiles count inclusively (`(|dx| + 1) * (|dy| + 1)`); find the largest such
rectangle.

**Part 2:** consecutive red tiles are joined by straight green lines and the
enclosed area is filled green. A rectangle now counts only when every tile it
covers is red or green; find the largest.

(Part 2's statement needs a login, so its rule - same two red corners, only
red and green tiles - comes from the puzzle notes: it is the only reading
that produces the documented example answer of 24.)

---

*Solution spoilers below.*

## Approach

Part 1 scans all pairs. For part 2, all boundary lines sit on coordinates of
the loop, so x and y compress to the distinct coordinates plus one
representative per gap, which makes every strip of that grid a single colour.
A sweep per row marks the strips a line runs through and counts the vertical
lines right of each strip (an odd count means enclosed); a prefix sum over
the uncovered strips then checks a rectangle in four lookups.

Both parts run in `O(n^2)` time and memory.

## Run / test

```sh
cd 2025
zig test src/day9.zig
zig run src/day9.zig < src/examples/day9.txt   # Part 1: 50, Part 2: 24
```
