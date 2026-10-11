# Day 24: Arithmetic Logic Unit

Instructions are parsed into typed values and an interpreter runs them over the
registers `w`, `x`, `y` and `z`, rejecting the crashes the puzzle warns about:
`div` by zero and `mod` with a negative left or non-positive right operand.

MONAD repeats one 18-instruction block per digit, parameterised by `div` (1 or
26), `add x` and `add y`; each block is matched against that template, so other
shapes are rejected. `z` is a base-26 stack: a `div 1` block pushes
`digit + add_y` (its `add x` exceeds every digit, so its comparison always
fails), and a `div 26` block pops it, leaving `z` small only when
`digits[pop] == digits[push] + add_y[push] + add_x[pop]`. Part 1 takes the
largest digit pair satisfying each constraint, part 2 the smallest, and the
interpreter validates both: `O(n)` time and space, not 9^14 candidates.

`input/day24-example.txt` is a synthetic 14-digit MONAD: the puzzle ships only
small ALU examples (negation, triple check, binary split; all tested), no MONAD.

From `2021/rust`:

```sh
cargo run --bin day24 < input/day24-example.txt
# Part 1: 99999596891916
# Part 2: 49411151121911
cargo test --bin day24
```
