# Day 22: Reactor Reboot

Each reboot step switches every cube in a cuboid on or off. Rather than tracking
single cubes, the solver keeps signed cuboids: an `on` step adds its cuboid with
sign +1, and every step first adds its overlap with each already counted cuboid
under the opposite sign, so no cube is counted twice. The cubes left on are the
sum of sign * volume. Part 1 clips each step to the initialization region
(`-50..50` per axis) and drops steps outside it; part 2 repeats the same pass
unclipped. With `n` steps the list grows per step: `O(n^2)` intersections and
`O(n^2)` space worst case, milliseconds for both parts.

`testdata/example.txt` is the puzzle's part 2 example (official answers 474140
and 2758514936282235); the part 1 examples (39, 590784) live in `reactor_test.go`.

From `2021/go`:

```sh
go vet ./day22
go test -race ./day22
go run ./day22 < day22/testdata/example.txt
```

The example prints exactly:

```text
Part 1: 474140
Part 2: 2758514936282235
```
