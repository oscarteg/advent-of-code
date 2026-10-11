# Day 20: Pulse Propagation

- **Part 1:** Simulate the button press as a BFS over pulses: flip-flops
  toggle and emit on a low pulse, conjunctions remember the last pulse per
  input and emit low only when every input is remembered high, broadcaster
  repeats its pulse. Push the button 1000 times, multiply total low pulses
  by total high pulses (including each press's own low pulse to
  `broadcaster`).
- **Part 2:** `rx` is never defined with its own line, only referenced as a
  destination, so it is found as the graph's one undefined sink. It has
  exactly one input, a conjunction fed by a handful of sub-networks that
  each raise a periodic high pulse. The answer is the LCM of each
  sub-network's first high pulse into that conjunction; exercised on a
  synthetic graph since the puzzle gives no part 2 example.
- Runtime is `O(presses * pulses per press)` for part 1 and bounded by the
  longest sub-network cycle for part 2; both are linear in the module graph.

From `2023/go`:

```sh
go run ./day20 < day20/testdata/example.txt
go test -race ./day20
go vet ./day20
```

The official example prints `Part 1: 32000000` (part 2 has no official
example; see the synthetic-cycle test in `solve_test.go`).
