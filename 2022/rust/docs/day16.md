# Day 16: Proboscidea Volcanium

The scan is a graph of valves; only the valves with a positive flow rate are
worth visiting. Floyd-Warshall turns the tunnels into walking times between
`AA` and those valves, so a route is a permutation of valves to open rather
than a step-by-step walk.

## Part 1
Depth-first search from `AA` over the valves still closed: walking to a valve
and opening it costs `distance + 1` minutes, and the valve then releases
`minutes_left * flow`. Recording the best total per set of opened valves (a
bitmask) gives part 1 as the maximum over all sets, within 30 minutes.

## Part 2
With 26 minutes the elephant and you open disjoint sets, so the same table
suffices. A subset maximum pass (`best[set]` becomes the best over every subset
of `set`) lets each set be paired with its complement; the answer is the best
`best[set] + best[!set]`. Complexity: O(2^v * v) for the pass plus the search
over at most v! routes, pruned hard by the 26/30-minute budget.

## Run / test
```bash
cargo test --bin day16
cargo run -q --bin day16 < input/day16-example.txt
```
