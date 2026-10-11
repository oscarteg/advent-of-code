# Day 6: Lanternfish

Each fish has a timer counting down to spawning. Rather than simulating every
fish (the population grows exponentially), track how many fish share each
timer value (0-8) in a 9-slot bucket array. Each day, fish at timer 0 spawn: the
buckets rotate down by one, timer-0 fish reset to 6, and an equal number of new
fish start at timer 8.

Part 1 runs the simulation for 80 days, part 2 for 256 days; both sum all
buckets.

With `D` days, runtime is `O(D)` and space is `O(1)`, independent of the fish
population size.

From `2021/go`:

```sh
go test -race ./day06
go run ./day06 < day06/testdata/example.txt
```

The example prints exactly:

```text
Part 1: 5934
Part 2: 26984457539
```
