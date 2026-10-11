# Day 18: Snailfish

A snailfish number is stored flat: every regular number with its nesting depth.
Neighbours in that list are exactly the numbers an explode feeds, so no tree
walking is needed.

Part 1 adds the numbers in order. An addition concatenates both lists with every
depth raised by one, then reduces: explode the leftmost neighbours sharing a
depth of five or more (replaced by a zero one level up), otherwise split the
leftmost number of ten or more. Magnitude folds equal depths with a stack.

Part 2 adds every ordered pair of two different numbers and keeps the largest
magnitude, because snailfish addition is not commutative.

For `n` numbers of `m` regular numbers each, a reduce is `O(m^2)`, part 1 is
`O(n * m^2)` and part 2 is `O(n^2 * m^2)`. Space is `O(n * m)`.

From `2021/rust`:

```sh
cargo run -q --bin day18 < input/day18-example.txt
# Part 1: 4140
# Part 2: 3993
cargo test --bin day18
```
