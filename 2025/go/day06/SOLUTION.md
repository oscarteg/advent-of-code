# Day 6: Trash Compactor

## Approach

Lines are split without trimming (whitespace is significant) and padded to
the widest line. A column that is a space on every row, including the
operator row, is a separator; contiguous non-separator column runs are the
problems.

**Part 1**: each number row's slice of a problem's columns (trimmed) is one
operand, folded left to right with the problem's operator.

**Part 2**: operands come from columns instead of rows, scanned right to
left; each column's digits (top to bottom, skipping blanks) form one
operand, folded the same way.

`solveByRows` and `solveByColumns` share `problems()`/`operator()`; only the
operand source differs.

## Complexity

O(rows * width) time and space: every cell is visited once per part.

## Run / Test

```bash
cd 2025/go
go test ./day06
go run ./day06 < day06/testdata/example.txt
```
