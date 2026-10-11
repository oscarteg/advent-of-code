package main

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

var errInvalidEquation = errors.New("invalid equation")

type equation struct {
	target   int64
	operands []int64
}

func parseEquations(input string) ([]equation, error) {
	var equations []equation
	for _, line := range strings.Split(strings.TrimSpace(input), "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		eq, err := parseEquation(line)
		if err != nil {
			return nil, err
		}
		equations = append(equations, eq)
	}
	if len(equations) == 0 {
		return nil, fmt.Errorf("%w: no equations found", errInvalidEquation)
	}
	return equations, nil
}

func parseEquation(line string) (equation, error) {
	parts := strings.SplitN(line, ":", 2)
	if len(parts) != 2 {
		return equation{}, fmt.Errorf("%w: missing ':' in %q", errInvalidEquation, line)
	}
	target, err := strconv.ParseInt(strings.TrimSpace(parts[0]), 10, 64)
	if err != nil {
		return equation{}, fmt.Errorf("target: %w", err)
	}

	fields := strings.Fields(parts[1])
	if len(fields) == 0 {
		return equation{}, fmt.Errorf("%w: no operands in %q", errInvalidEquation, line)
	}
	operands := make([]int64, len(fields))
	for i, field := range fields {
		operand, err := strconv.ParseInt(field, 10, 64)
		if err != nil {
			return equation{}, fmt.Errorf("operand: %w", err)
		}
		operands[i] = operand
	}
	return equation{target: target, operands: operands}, nil
}

func concatenate(left, right int64) int64 {
	multiplier := int64(1)
	for multiplier <= right {
		multiplier *= 10
	}
	return left*multiplier + right
}

func canReach(target int64, operands []int64, index int, accumulated int64, allowConcat bool) bool {
	if accumulated > target {
		return false
	}
	if index == len(operands) {
		return accumulated == target
	}

	next := operands[index]
	if canReach(target, operands, index+1, accumulated+next, allowConcat) {
		return true
	}
	if canReach(target, operands, index+1, accumulated*next, allowConcat) {
		return true
	}
	if allowConcat && canReach(target, operands, index+1, concatenate(accumulated, next), allowConcat) {
		return true
	}
	return false
}

func (e equation) isSolvable(allowConcat bool) bool {
	return canReach(e.target, e.operands, 1, e.operands[0], allowConcat)
}

func solve(input string) (part1, part2 int64, err error) {
	equations, err := parseEquations(input)
	if err != nil {
		return 0, 0, err
	}

	for _, eq := range equations {
		switch {
		case eq.isSolvable(false):
			part1 += eq.target
			part2 += eq.target
		case eq.isSolvable(true):
			part2 += eq.target
		}
	}
	return part1, part2, nil
}
