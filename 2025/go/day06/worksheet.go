package main

import (
	"fmt"
	"strconv"
	"strings"
)

// problem is the half-open column range [start, end) of one problem on the
// worksheet, not including its separating all-space columns.
type problem struct {
	start, end int
}

// worksheet holds the padded number rows and operator row of a parsed
// trash-compactor worksheet.
type worksheet struct {
	numRows []string
	opRow   string
	width   int
}

// parseWorksheet splits input into lines without trimming surrounding
// whitespace (it is significant), pads every line to the widest line, and
// separates the trailing operator row from the number rows above it.
func parseWorksheet(input string) (worksheet, error) {
	lines := strings.Split(input, "\n")
	for i, line := range lines {
		lines[i] = strings.TrimRight(line, "\r")
	}
	for len(lines) > 0 && lines[len(lines)-1] == "" {
		lines = lines[:len(lines)-1]
	}
	if len(lines) < 2 {
		return worksheet{}, fmt.Errorf("parse worksheet: need at least one number row and an operator row, got %d lines", len(lines))
	}

	width := 0
	for _, line := range lines {
		if len(line) > width {
			width = len(line)
		}
	}
	if width == 0 {
		return worksheet{}, fmt.Errorf("parse worksheet: empty worksheet")
	}
	for i, line := range lines {
		lines[i] = line + strings.Repeat(" ", width-len(line))
	}

	return worksheet{numRows: lines[:len(lines)-1], opRow: lines[len(lines)-1], width: width}, nil
}

// problems finds the contiguous column ranges separated by all-space
// columns (checked across the number rows and the operator row).
func (w worksheet) problems() []problem {
	isSeparator := func(col int) bool {
		if w.opRow[col] != ' ' {
			return false
		}
		for _, row := range w.numRows {
			if row[col] != ' ' {
				return false
			}
		}
		return true
	}

	var problems []problem
	col := 0
	for col < w.width {
		if isSeparator(col) {
			col++
			continue
		}
		start := col
		for col < w.width && !isSeparator(col) {
			col++
		}
		problems = append(problems, problem{start: start, end: col})
	}
	return problems
}

func (w worksheet) operator(p problem) (byte, error) {
	for _, ch := range w.opRow[p.start:p.end] {
		if ch == '+' || ch == '*' {
			return byte(ch), nil
		}
	}
	return 0, fmt.Errorf("find operator: no + or * in columns %d-%d", p.start, p.end)
}

func applyOperator(op byte, acc, n int64) int64 {
	if op == '+' {
		return acc + n
	}
	return acc * n
}

// solveByRows treats each number row's column slice within a problem as one
// operand (part 1).
func (w worksheet) solveByRows(p problem, op byte) (int64, error) {
	var acc int64
	hasAcc := false
	for _, row := range w.numRows {
		seg := strings.TrimSpace(row[p.start:p.end])
		if seg == "" {
			continue
		}
		n, err := strconv.ParseInt(seg, 10, 64)
		if err != nil {
			return 0, fmt.Errorf("solve by rows: %w", err)
		}
		if hasAcc {
			acc = applyOperator(op, acc, n)
		} else {
			acc, hasAcc = n, true
		}
	}
	if !hasAcc {
		return 0, fmt.Errorf("solve by rows: problem %d-%d has no operands", p.start, p.end)
	}
	return acc, nil
}

// solveByColumns treats each column, right to left, as one operand formed by
// its digits read top to bottom (part 2).
func (w worksheet) solveByColumns(p problem, op byte) (int64, error) {
	var acc int64
	hasAcc := false
	for col := p.end - 1; col >= p.start; col-- {
		var digits []byte
		for _, row := range w.numRows {
			ch := row[col]
			if ch == ' ' {
				continue
			}
			if ch < '0' || ch > '9' {
				return 0, fmt.Errorf("solve by columns: invalid digit %q at column %d", ch, col)
			}
			digits = append(digits, ch)
		}
		if len(digits) == 0 {
			continue
		}
		n, err := strconv.ParseInt(string(digits), 10, 64)
		if err != nil {
			return 0, fmt.Errorf("solve by columns: %w", err)
		}
		if hasAcc {
			acc = applyOperator(op, acc, n)
		} else {
			acc, hasAcc = n, true
		}
	}
	if !hasAcc {
		return 0, fmt.Errorf("solve by columns: problem %d-%d has no operands", p.start, p.end)
	}
	return acc, nil
}

func solve(input string) (part1, part2 int64, err error) {
	w, err := parseWorksheet(input)
	if err != nil {
		return 0, 0, err
	}

	problems := w.problems()
	if len(problems) == 0 {
		return 0, 0, fmt.Errorf("solve: no problems found on worksheet")
	}

	for _, p := range problems {
		op, err := w.operator(p)
		if err != nil {
			return 0, 0, err
		}
		rowsTotal, err := w.solveByRows(p, op)
		if err != nil {
			return 0, 0, err
		}
		colsTotal, err := w.solveByColumns(p, op)
		if err != nil {
			return 0, 0, err
		}
		part1 += rowsTotal
		part2 += colsTotal
	}

	return part1, part2, nil
}
