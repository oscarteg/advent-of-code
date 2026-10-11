# Day 7: The Treachery of Whales

Part 1 aligns at a median position and sums absolute distances, minimizing the
linear fuel cost. Part 2 evaluates the convex triangular-distance fuel function
with binary search for its minimum. Parsing and sorting take `O(n log n)` time;
Part 2 takes `O(n log R)`, where `R` is the position range. Memory use is `O(n)`.

Run and test from this directory:

```sh
cargo run -q --bin day7 < input/day7-example.txt
cargo test --bin day7
cargo clippy --bin day7 -- -D warnings
```
