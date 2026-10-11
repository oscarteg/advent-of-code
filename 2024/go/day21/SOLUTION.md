# Day 21: Keypad Conundrum

Every code is typed on the door's numeric keypad by a robot that is itself commanded from a
directional keypad by another robot, and so on; the human types the last keypad. The answer
sums each code's complexity: the presses the human needs times the code's numeric part.

## Part 1

Two directional-keypad robots. An arm only walks a shortest path between two keys, and just two
orders can be shortest: every horizontal step first, or every vertical step first - one of them
possibly blocked because its turning point is the keypad's gap. Pressing a key therefore costs
the cheaper of typing `route + "A"` one keypad up, and with no robot left the human just presses
the route's buttons. Complexity: `O(len * keys^2)` per code.

## Part 2

The same recursion with 25 robots, which is only affordable memoised. A move's cost depends on
nothing but the keypad, the key pair and the robots above it, so `solver.cache` keyed on those
three collapses the exponential expansion. Complexity: `O(robots * keys^2)` cached moves.

## Run / test

```bash
go test -race -count=1 ./day21
go run ./day21 < day21/testdata/example.txt
```

The example prints `Part 1: 126384` and `Part 2: 154115708116294`.
