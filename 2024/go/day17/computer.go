package main

import (
	"errors"
	"fmt"
	"slices"
	"strconv"
	"strings"
)

var errInvalidInput = errors.New("invalid input")

type machine struct {
	a, b, c uint64
	program []uint8
}

type registers struct {
	a, b, c uint64
}

func (r registers) combo(operand uint8) (uint64, error) {
	switch operand {
	case 4:
		return r.a, nil
	case 5:
		return r.b, nil
	case 6:
		return r.c, nil
	case 7:
		return 0, fmt.Errorf("%w: combo operand 7 is reserved", errInvalidInput)
	default:
		return uint64(operand), nil
	}
}

func registerValue(line, label string) (uint64, error) {
	rest, found := strings.CutPrefix(line, label)
	if !found {
		return 0, fmt.Errorf("%w: want a line starting with %q, got %q", errInvalidInput, label, line)
	}
	value, err := strconv.ParseUint(strings.TrimSpace(rest), 10, 64)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", label, err)
	}
	return value, nil
}

func parse(input string) (machine, error) {
	var lines []string
	for _, line := range strings.Split(input, "\n") {
		if trimmed := strings.TrimSpace(line); trimmed != "" {
			lines = append(lines, trimmed)
		}
	}
	if len(lines) != 4 {
		return machine{}, fmt.Errorf("%w: want three register lines and one program line, got %d lines", errInvalidInput, len(lines))
	}

	var parsed machine
	var err error
	if parsed.a, err = registerValue(lines[0], "Register A:"); err != nil {
		return machine{}, err
	}
	if parsed.b, err = registerValue(lines[1], "Register B:"); err != nil {
		return machine{}, err
	}
	if parsed.c, err = registerValue(lines[2], "Register C:"); err != nil {
		return machine{}, err
	}

	fields, found := strings.CutPrefix(lines[3], "Program:")
	if !found {
		return machine{}, fmt.Errorf("%w: want a program line, got %q", errInvalidInput, lines[3])
	}
	for _, field := range strings.Split(fields, ",") {
		value, err := strconv.ParseUint(strings.TrimSpace(field), 10, 3)
		if err != nil {
			return machine{}, fmt.Errorf("program value: %w", err)
		}
		parsed.program = append(parsed.program, uint8(value))
	}
	return parsed, nil
}

// Opcodes 0, 2, 5, 6 and 7 read a combo operand; bxl and jnz take a literal one.
const comboOpcodes = 1<<0 | 1<<2 | 1<<5 | 1<<6 | 1<<7

func run(m machine, initialA uint64) ([]uint8, error) {
	r := registers{a: initialA, b: m.b, c: m.c}
	var output []uint8

	for ip := 0; ip+1 < len(m.program); {
		opcode, operand := m.program[ip], m.program[ip+1]
		next := ip + 2

		var combo uint64
		if comboOpcodes&(1<<opcode) != 0 {
			value, err := r.combo(operand)
			if err != nil {
				return nil, err
			}
			combo = value
		}

		switch opcode {
		case 0:
			r.a >>= combo
		case 1:
			r.b ^= uint64(operand)
		case 2:
			r.b = combo % 8
		case 3:
			if r.a != 0 {
				next = int(operand)
			}
		case 4:
			r.b ^= r.c
		case 5:
			output = append(output, uint8(combo%8))
		case 6:
			r.b = r.a >> combo
		case 7:
			r.c = r.a >> combo
		}
		ip = next
	}
	return output, nil
}

// A self-replicating program emits one value per loop and then drops the low
// three bits of A, so A is built three bits at a time from the last value
// backwards, keeping every candidate whose run reproduces the program's tail.
func findSelfReplicatingA(m machine) (uint64, error) {
	candidates := []uint64{0}
	for index := len(m.program) - 1; index >= 0; index-- {
		var next []uint64
		for _, prefix := range candidates {
			for digit := uint64(0); digit < 8; digit++ {
				candidate := prefix*8 + digit
				output, err := run(m, candidate)
				if err != nil {
					return 0, err
				}
				if slices.Equal(output, m.program[index:]) {
					next = append(next, candidate)
				}
			}
		}
		candidates = next
	}
	if len(candidates) == 0 {
		return 0, nil
	}
	return slices.Min(candidates), nil
}

func joinValues(values []uint8) string {
	digits := make([]string, len(values))
	for i, value := range values {
		digits[i] = strconv.FormatUint(uint64(value), 10)
	}
	return strings.Join(digits, ",")
}

// part2 is 0 when no initial A makes the program print a copy of itself.
func solve(input string) (part1 string, part2 uint64, err error) {
	m, err := parse(input)
	if err != nil {
		return "", 0, err
	}
	output, err := run(m, m.a)
	if err != nil {
		return "", 0, err
	}
	selfReplicatingA, err := findSelfReplicatingA(m)
	if err != nil {
		return "", 0, err
	}
	return joinValues(output), selfReplicatingA, nil
}
