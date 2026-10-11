package main

import (
	"errors"
	"fmt"
	"sort"
	"strconv"
	"strings"
)

var errInvalidMonkey = errors.New("invalid monkey business input")

const (
	roundsPart1  = 20
	roundsPart2  = 10000
	reliefFactor = 3
)

type operator int

const (
	opAdd operator = iota
	opMultiply
)

type operation struct {
	op      operator
	useOld  bool
	operand uint64
}

func (o operation) apply(old uint64) uint64 {
	rhs := o.operand
	if o.useOld {
		rhs = old
	}
	if o.op == opAdd {
		return old + rhs
	}
	return old * rhs
}

type monkey struct {
	items       []uint64
	op          operation
	divisor     uint64
	ifTrue      int
	ifFalse     int
	inspections uint64
}

func parseSuffixNumbers(line, prefix string) ([]uint64, error) {
	rest := strings.TrimSpace(line)
	if !strings.HasPrefix(rest, prefix) {
		return nil, fmt.Errorf("%w: expected prefix %q in %q", errInvalidMonkey, prefix, line)
	}
	rest = strings.TrimPrefix(rest, prefix)
	fields := strings.Split(rest, ",")
	numbers := make([]uint64, 0, len(fields))
	for _, field := range fields {
		n, err := strconv.ParseUint(strings.TrimSpace(field), 10, 64)
		if err != nil {
			return nil, fmt.Errorf("%w: invalid number in %q: %v", errInvalidMonkey, line, err)
		}
		numbers = append(numbers, n)
	}
	return numbers, nil
}

func parseSuffixNumber(line, prefix string) (uint64, error) {
	rest := strings.TrimSpace(line)
	if !strings.HasPrefix(rest, prefix) {
		return 0, fmt.Errorf("%w: expected prefix %q in %q", errInvalidMonkey, prefix, line)
	}
	n, err := strconv.ParseUint(strings.TrimSpace(strings.TrimPrefix(rest, prefix)), 10, 64)
	if err != nil {
		return 0, fmt.Errorf("%w: invalid number in %q: %v", errInvalidMonkey, line, err)
	}
	return n, nil
}

func parseOperation(line string) (operation, error) {
	rest := strings.TrimSpace(line)
	const prefix = "Operation: new = old "
	if !strings.HasPrefix(rest, prefix) {
		return operation{}, fmt.Errorf("%w: invalid operation line %q", errInvalidMonkey, line)
	}
	fields := strings.Fields(strings.TrimPrefix(rest, prefix))
	if len(fields) != 2 {
		return operation{}, fmt.Errorf("%w: invalid operation line %q", errInvalidMonkey, line)
	}
	var op operation
	switch fields[0] {
	case "+":
		op.op = opAdd
	case "*":
		op.op = opMultiply
	default:
		return operation{}, fmt.Errorf("%w: unsupported operator %q", errInvalidMonkey, fields[0])
	}
	if fields[1] == "old" {
		op.useOld = true
		return op, nil
	}
	value, err := strconv.ParseUint(fields[1], 10, 64)
	if err != nil {
		return operation{}, fmt.Errorf("%w: invalid operand %q: %v", errInvalidMonkey, fields[1], err)
	}
	op.operand = value
	return op, nil
}

func parseMonkey(block string) (monkey, error) {
	var m monkey
	lines := make([]string, 0, 6)
	for _, line := range strings.Split(strings.TrimSpace(block), "\n") {
		trimmed := strings.TrimSpace(line)
		if trimmed != "" {
			lines = append(lines, trimmed)
		}
	}
	if len(lines) != 6 {
		return m, fmt.Errorf("%w: monkey block needs six lines", errInvalidMonkey)
	}

	items, err := parseSuffixNumbers(lines[1], "Starting items:")
	if err != nil {
		return m, err
	}
	m.items = items

	op, err := parseOperation(lines[2])
	if err != nil {
		return m, err
	}
	m.op = op

	divisor, err := parseSuffixNumber(lines[3], "Test: divisible by")
	if err != nil {
		return m, err
	}
	m.divisor = divisor

	ifTrue, err := parseSuffixNumber(lines[4], "If true: throw to monkey")
	if err != nil {
		return m, err
	}
	m.ifTrue = int(ifTrue)

	ifFalse, err := parseSuffixNumber(lines[5], "If false: throw to monkey")
	if err != nil {
		return m, err
	}
	m.ifFalse = int(ifFalse)

	return m, nil
}

func parse(input string) ([]monkey, error) {
	blocks := strings.Split(strings.TrimSpace(input), "\n\n")
	monkeys := make([]monkey, 0, len(blocks))
	for _, block := range blocks {
		if strings.TrimSpace(block) == "" {
			continue
		}
		m, err := parseMonkey(block)
		if err != nil {
			return nil, err
		}
		monkeys = append(monkeys, m)
	}
	if len(monkeys) == 0 {
		return nil, fmt.Errorf("%w: no monkeys found", errInvalidMonkey)
	}
	for _, m := range monkeys {
		if m.divisor == 0 {
			return nil, fmt.Errorf("%w: divisor must not be zero", errInvalidMonkey)
		}
		if m.ifTrue >= len(monkeys) || m.ifFalse >= len(monkeys) {
			return nil, fmt.Errorf("%w: throw target out of range", errInvalidMonkey)
		}
	}
	return monkeys, nil
}

func run(monkeys []monkey, rounds int, relief func(uint64) uint64) uint64 {
	for round := 0; round < rounds; round++ {
		for i := range monkeys {
			for len(monkeys[i].items) > 0 {
				item := monkeys[i].items[0]
				monkeys[i].items = monkeys[i].items[1:]
				monkeys[i].inspections++
				worry := relief(monkeys[i].op.apply(item))
				target := monkeys[i].ifFalse
				if worry%monkeys[i].divisor == 0 {
					target = monkeys[i].ifTrue
				}
				monkeys[target].items = append(monkeys[target].items, worry)
			}
		}
	}
	counts := make([]uint64, len(monkeys))
	for i, m := range monkeys {
		counts[i] = m.inspections
	}
	sort.Slice(counts, func(i, j int) bool { return counts[i] > counts[j] })
	return counts[0] * counts[1]
}

func solve(input string) (uint64, uint64, error) {
	part1Monkeys, err := parse(input)
	if err != nil {
		return 0, 0, err
	}
	part1 := run(part1Monkeys, roundsPart1, func(worry uint64) uint64 { return worry / reliefFactor })

	part2Monkeys, err := parse(input)
	if err != nil {
		return 0, 0, err
	}
	var modulus uint64 = 1
	for _, m := range part2Monkeys {
		modulus *= m.divisor
	}
	part2 := run(part2Monkeys, roundsPart2, func(worry uint64) uint64 { return worry % modulus })

	return part1, part2, nil
}
