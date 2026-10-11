# Day 24: Crossed Wires (C++)

Part 1 evaluates every wire named `z..` with a memoised depth-first walk over the gate graph
(a revisited wire means a cycle, which is rejected) and reads those bits as a number, `z00`
least significant.

## Approach

Part 2: the real input is a ripple-carry adder with four swapped gate output wires. Gates are
indexed by their inputs plus operator, which a swap never changes, so the expected adder can be
walked from bit 0 up: `z0 = x0 XOR y0`, `carry0 = x0 AND y0`, then
`z_i = (x_i XOR y_i) XOR carry_(i-1)` and
`carry_i = (x_i AND y_i) OR ((x_i XOR y_i) AND carry_(i-1))`, the last carry on the top `z` wire.
The walk returns the lowest bit that does not match; a depth-first search swaps pairs of gate
outputs, keeps only swaps that push that bit higher and recurses until the chain verifies. The
eight names are reported sorted and comma-joined. Part 2 is empty when no four swaps repair the
circuit: both part 1 examples, and the official part 2 example, a different (AND) circuit.

## Complexity

Part 1 is O(gates); one walk is O(bits) and the search tries O(gates^2) swaps per level for at
most four levels: 0.1 s for a 45-bit adder.

## Run / test

```bash
cd 2024/cpp
make test day=day24
make -s run day=day24 < examples/day24.txt
```
