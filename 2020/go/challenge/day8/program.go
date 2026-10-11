package day8

import (
	"fmt"
	"strconv"
	"strings"
)

type opcode string

const (
	opAcc opcode = "acc"
	opJmp opcode = "jmp"
	opNop opcode = "nop"
)

type instruction struct {
	op  opcode
	arg int
}

func parseProgram(lines []string) ([]instruction, error) {
	var program []instruction
	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}
		op, argText, found := strings.Cut(line, " ")
		if !found {
			return nil, fmt.Errorf("malformed instruction: %q", line)
		}
		arg, err := strconv.Atoi(argText)
		if err != nil {
			return nil, fmt.Errorf("malformed argument %q: %w", argText, err)
		}
		switch opcode(op) {
		case opAcc, opJmp, opNop:
			program = append(program, instruction{op: opcode(op), arg: arg})
		default:
			return nil, fmt.Errorf("unknown operation: %q", op)
		}
	}
	if len(program) == 0 {
		return nil, fmt.Errorf("program must not be empty")
	}
	return program, nil
}

// run executes the program and reports whether it terminated by running past
// the last instruction, along with the accumulator at the point it stopped
// (either on termination or on the first repeated instruction).
func run(program []instruction) (terminated bool, acc int) {
	visited := make([]bool, len(program))
	pc := 0
	for pc >= 0 && pc < len(program) {
		if visited[pc] {
			return false, acc
		}
		visited[pc] = true
		switch program[pc].op {
		case opAcc:
			acc += program[pc].arg
			pc++
		case opJmp:
			pc += program[pc].arg
		case opNop:
			pc++
		}
	}
	return true, acc
}

func swapped(op opcode) (opcode, bool) {
	switch op {
	case opJmp:
		return opNop, true
	case opNop:
		return opJmp, true
	default:
		return "", false
	}
}

func fixAndRun(program []instruction) (int, error) {
	for i, instr := range program {
		newOp, ok := swapped(instr.op)
		if !ok {
			continue
		}
		patched := make([]instruction, len(program))
		copy(patched, program)
		patched[i] = instruction{op: newOp, arg: instr.arg}
		if terminated, acc := run(patched); terminated {
			return acc, nil
		}
	}
	return 0, fmt.Errorf("no single jmp/nop swap makes the program terminate")
}
