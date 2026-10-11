# Day 10: Factory

## Problem

One machine per line: an indicator-light diagram in `[]`, button wirings in
`()`, and joltage requirements in `{}`.

**Part 1:** each button press toggles the lights it is wired to. Per machine,
find the fewest presses that light up exactly the diagram; sum over machines.

**Part 2:** the lights are irrelevant: a press adds one to each counter its
button feeds. Per machine, find the fewest presses that make the counters
equal the joltage requirements; sum over machines.

(Part 1's per-machine presses (2, 3, 2) come from the official puzzle text;
part 2's (10, 12, 11) from the task brief, whose statement needs a login.)

---

*Solution spoilers below.*

## Approach

**Part 1:** toggling twice cancels out, so each button is pressed once or not
at all and the diagram is a linear system over GF(2). Elimination yields one
solution plus a null-space basis - button sets that toggle nothing - and the
answer is the lightest solution in that coset. `O(2^(b-r) * b)` for `b`
buttons of rank `r`.

**Part 2:** the targets are `A x = joltage` over the nonnegative integers.
Fraction-free Gaussian elimination makes the pivot presses a function of the
free ones, each free press is bounded by the smallest target it feeds, and a
pruned depth-first search over that box gives the cheapest total.
`O(l * b^2)` plus the search.

## Run / test

```sh
cd 2025
zig test src/day10.zig
zig run src/day10.zig < src/examples/day10.txt  # Part 1: 7, Part 2: 33
```
