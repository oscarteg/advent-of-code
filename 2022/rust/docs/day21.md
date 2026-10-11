# Day 21: Monkey Math

## Part 1

Parse each line into a named job: a literal number, or an operator on two other
monkeys. Evaluate `root` recursively, O(n) over a tree-shaped riddle. A budget
of one recursion step per monkey turns a self-referential riddle into an error
instead of a hang.

## Part 2

`root` becomes an equality and `humn` the unknown. Map every monkey to the
monkey listening to it, walk from `humn` up to `root`, then back down: at each
step the branch without `humn` evaluates to a known number, so the operation is
inverted (`a + b = t` gives `a = t - b`, `a / b = t` gives `a = t * b`, ...)
until it pins down `humn`. Divisions must be exact, otherwise the riddle has no
integer answer. Also O(n): every subtree hanging off the path is evaluated once.

## Run and test

```bash
cargo run -q --bin day21 < input/day21-example.txt
cargo test --bin day21
```
