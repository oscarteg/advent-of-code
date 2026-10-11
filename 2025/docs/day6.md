# Day 6: Trash Compactor

## Problem

The input is a worksheet of arithmetic problems laid out side by side: rows of
numbers with an operator row underneath, and columns that are blank on every
line separating one problem from the next.

**Part 1:** solve each problem reading each row's slice as one operand, fold the
operands with the problem's operator, and sum the results.

**Part 2:** solve each problem reading its columns right to left instead: each
column's digits (top to bottom) form one operand. Sum the results.

---

*Solution spoilers below.*

## Approach

The worksheet is parsed without trimming lines: lines are padded to the
widest line with spaces so every line can be indexed by column. A column
that is a space on every line (including the operator row) is a separator;
contiguous runs of non-separator columns are the problems.

Both parts share `findProblems`/`findOperator`; only the operand source
differs (`solveRows` vs `solveColumns`, the latter skipping blank rows).

## Complexity

O(rows * width) time and space: every cell is visited once per part.

## Run / Test

```bash
cd 2025
zig test src/day6.zig
zig run src/day6.zig < src/examples/day6.txt
```
