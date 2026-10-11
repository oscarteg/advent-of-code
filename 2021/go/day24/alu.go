package main

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

var errInvalidProgram = errors.New("invalid ALU program")

type register int

const (
	registerW register = iota
	registerX
	registerY
	registerZ
)

func parseRegister(token string) (register, error) {
	switch token {
	case "w":
		return registerW, nil
	case "x":
		return registerX, nil
	case "y":
		return registerY, nil
	case "z":
		return registerZ, nil
	}
	return 0, fmt.Errorf("%w: unknown register %q", errInvalidProgram, token)
}

type operand struct {
	register  register
	literal   int64
	isLiteral bool
}

func parseOperand(token string) (operand, error) {
	switch token {
	case "w", "x", "y", "z":
		source, err := parseRegister(token)
		if err != nil {
			return operand{}, err
		}
		return operand{register: source}, nil
	}
	literal, err := strconv.ParseInt(token, 10, 64)
	if err != nil {
		return operand{}, fmt.Errorf("%w: operand %q: %w", errInvalidProgram, token, err)
	}
	return operand{literal: literal, isLiteral: true}, nil
}

// The five binary operations; inp is opInput, which takes no second operand.
type operation int

const (
	opInput operation = iota
	opAdd
	opMul
	opDiv
	opMod
	opEql
)

type instruction struct {
	op      operation
	target  register
	operand operand
}

func parseInstruction(line string) (instruction, error) {
	fields := strings.Fields(line)
	if len(fields) < 2 {
		return instruction{}, fmt.Errorf("%w: incomplete instruction %q", errInvalidProgram, line)
	}
	target, err := parseRegister(fields[1])
	if err != nil {
		return instruction{}, err
	}
	if fields[0] == "inp" {
		if len(fields) != 2 {
			return instruction{}, fmt.Errorf("%w: trailing tokens in %q", errInvalidProgram, line)
		}
		return instruction{op: opInput, target: target}, nil
	}
	if len(fields) != 3 {
		return instruction{}, fmt.Errorf("%w: %q takes a target and one operand", errInvalidProgram, line)
	}
	var op operation
	switch fields[0] {
	case "add":
		op = opAdd
	case "mul":
		op = opMul
	case "div":
		op = opDiv
	case "mod":
		op = opMod
	case "eql":
		op = opEql
	default:
		return instruction{}, fmt.Errorf("%w: unknown instruction %q", errInvalidProgram, fields[0])
	}
	argument, err := parseOperand(fields[2])
	if err != nil {
		return instruction{}, err
	}
	return instruction{op: op, target: target, operand: argument}, nil
}

func parseProgram(input string) ([]instruction, error) {
	lines := strings.Split(strings.ReplaceAll(input, "\r\n", "\n"), "\n")
	program := make([]instruction, 0, len(lines))
	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}
		parsed, err := parseInstruction(line)
		if err != nil {
			return nil, err
		}
		program = append(program, parsed)
	}
	return program, nil
}

// apply runs one binary operation and rejects the cases that crash the ALU.
func apply(op operation, left, right int64) (int64, error) {
	switch op {
	case opAdd:
		return left + right, nil
	case opMul:
		return left * right, nil
	case opDiv:
		if right == 0 {
			return 0, fmt.Errorf("%w: div %d 0 crashes the ALU", errInvalidProgram, left)
		}
		return left / right, nil
	case opMod:
		if left < 0 || right <= 0 {
			return 0, fmt.Errorf("%w: mod %d %d crashes the ALU", errInvalidProgram, left, right)
		}
		return left % right, nil
	case opEql:
		if left == right {
			return 1, nil
		}
		return 0, nil
	}
	return 0, fmt.Errorf("%w: unsupported operation", errInvalidProgram)
}

// run feeds inputs to the inp instructions in order and returns [w, x, y, z].
func run(program []instruction, inputs []int64) ([4]int64, error) {
	var registers [4]int64
	consumed := 0
	for _, current := range program {
		if current.op == opInput {
			if consumed == len(inputs) {
				return registers, fmt.Errorf("%w: the program reads more inputs than supplied", errInvalidProgram)
			}
			registers[current.target] = inputs[consumed]
			consumed++
			continue
		}
		right := current.operand.literal
		if !current.operand.isLiteral {
			right = registers[current.operand.register]
		}
		result, err := apply(current.op, registers[current.target], right)
		if err != nil {
			return registers, err
		}
		registers[current.target] = result
	}
	return registers, nil
}

const blockLength = 18

// block holds the parameters of one MONAD digit block: it divides z by div,
// compares z%26+addX with the digit and pushes digit+addY when they differ.
type block struct {
	div  int64
	addX int64
	addY int64
}

// instructions spells the block out as MONAD writes it, so that programs of any
// other shape can be rejected.
func (b block) instructions() [blockLength]instruction {
	fromRegister := func(op operation, target, source register) instruction {
		return instruction{op: op, target: target, operand: operand{register: source}}
	}
	fromLiteral := func(op operation, target register, value int64) instruction {
		return instruction{op: op, target: target, operand: operand{literal: value, isLiteral: true}}
	}
	return [blockLength]instruction{
		{op: opInput, target: registerW},
		fromLiteral(opMul, registerX, 0),
		fromRegister(opAdd, registerX, registerZ),
		fromLiteral(opMod, registerX, 26),
		fromLiteral(opDiv, registerZ, b.div),
		fromLiteral(opAdd, registerX, b.addX),
		fromRegister(opEql, registerX, registerW),
		fromLiteral(opEql, registerX, 0),
		fromLiteral(opMul, registerY, 0),
		fromLiteral(opAdd, registerY, 25),
		fromRegister(opMul, registerY, registerX),
		fromLiteral(opAdd, registerY, 1),
		fromRegister(opMul, registerZ, registerY),
		fromLiteral(opMul, registerY, 0),
		fromRegister(opAdd, registerY, registerW),
		fromLiteral(opAdd, registerY, b.addY),
		fromRegister(opMul, registerY, registerX),
		fromRegister(opAdd, registerZ, registerY),
	}
}

func literalOperand(current instruction) (int64, error) {
	if current.op == opInput || !current.operand.isLiteral {
		return 0, fmt.Errorf("%w: expected a literal operand", errInvalidProgram)
	}
	return current.operand.literal, nil
}

func parseBlock(instructions []instruction) (block, error) {
	if len(instructions) != blockLength {
		return block{}, fmt.Errorf("%w: a MONAD block has %d instructions, found %d", errInvalidProgram, blockLength, len(instructions))
	}
	div, err := literalOperand(instructions[4])
	if err != nil {
		return block{}, err
	}
	addX, err := literalOperand(instructions[5])
	if err != nil {
		return block{}, err
	}
	addY, err := literalOperand(instructions[15])
	if err != nil {
		return block{}, err
	}
	if div != 1 && div != 26 {
		return block{}, fmt.Errorf("%w: a MONAD block divides z by 1 or 26, found %d", errInvalidProgram, div)
	}
	parsed := block{div: div, addX: addX, addY: addY}
	for index, want := range parsed.instructions() {
		if instructions[index] != want {
			return block{}, fmt.Errorf("%w: instruction %d does not follow the MONAD template", errInvalidProgram, index)
		}
	}
	return parsed, nil
}

// constraint records that a div-26 block pops what an earlier div-1 block
// pushed, which forces digits[pop] == digits[push] + offset.
type constraint struct {
	push   int
	pop    int
	offset int64
}

type bound int

const (
	largest bound = iota
	smallest
)

type monad struct {
	program []instruction
	blocks  []block
}

func parseMonad(input string) (monad, error) {
	program, err := parseProgram(input)
	if err != nil {
		return monad{}, err
	}
	if len(program) == 0 {
		return monad{}, fmt.Errorf("%w: the program is empty", errInvalidProgram)
	}
	if len(program)%blockLength != 0 {
		return monad{}, fmt.Errorf("%w: MONAD repeats blocks of %d instructions, found %d in total", errInvalidProgram, blockLength, len(program))
	}
	blocks := make([]block, 0, len(program)/blockLength)
	for start := 0; start < len(program); start += blockLength {
		parsed, err := parseBlock(program[start : start+blockLength])
		if err != nil {
			return monad{}, err
		}
		blocks = append(blocks, parsed)
	}
	return monad{program: program, blocks: blocks}, nil
}

func (m monad) constraints() ([]constraint, error) {
	type pushedValue struct {
		index int
		addY  int64
	}
	stack := make([]pushedValue, 0, len(m.blocks))
	constraints := make([]constraint, 0, len(m.blocks)/2)
	for index, current := range m.blocks {
		if current.div == 1 {
			if current.addX <= 9 {
				return nil, fmt.Errorf("%w: block %d pushes, so add x %d must exceed every digit", errInvalidProgram, index, current.addX)
			}
			stack = append(stack, pushedValue{index: index, addY: current.addY})
			continue
		}
		if len(stack) == 0 {
			return nil, fmt.Errorf("%w: block %d pops while z is empty", errInvalidProgram, index)
		}
		top := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		constraints = append(constraints, constraint{push: top.index, pop: index, offset: top.addY + current.addX})
	}
	if len(stack) != 0 {
		return nil, fmt.Errorf("%w: %d pushed values are never popped, so z cannot end at 0", errInvalidProgram, len(stack))
	}
	return constraints, nil
}

// modelNumber picks the extreme digit for every constrained pair and validates
// the resulting model number with the interpreter.
func (m monad) modelNumber(want bound) (uint64, error) {
	constraints, err := m.constraints()
	if err != nil {
		return 0, err
	}
	digits := make([]int64, len(m.blocks))
	for _, current := range constraints {
		if current.offset < -8 || current.offset > 8 {
			return 0, fmt.Errorf("%w: digits %d and %d must differ by %d, which digits 1..9 cannot do", errInvalidProgram, current.push, current.pop, current.offset)
		}
		if want == largest {
			pop := min(int64(9), 9+current.offset)
			digits[current.push], digits[current.pop] = pop-current.offset, pop
		} else {
			push := max(int64(1), 1-current.offset)
			digits[current.push], digits[current.pop] = push, push+current.offset
		}
	}
	registers, err := run(m.program, digits)
	if err != nil {
		return 0, err
	}
	if registers[registerZ] != 0 {
		return 0, fmt.Errorf("%w: the chosen model number leaves z = %d", errInvalidProgram, registers[registerZ])
	}
	var number uint64
	for _, digit := range digits {
		if digit < 1 || digit > 9 {
			return 0, fmt.Errorf("%w: digit %d is not between 1 and 9", errInvalidProgram, digit)
		}
		number = number*10 + uint64(digit)
	}
	return number, nil
}

func solve(input string) (uint64, uint64, error) {
	parsed, err := parseMonad(input)
	if err != nil {
		return 0, 0, err
	}
	part1, err := parsed.modelNumber(largest)
	if err != nil {
		return 0, 0, err
	}
	part2, err := parsed.modelNumber(smallest)
	if err != nil {
		return 0, 0, err
	}
	return part1, part2, nil
}
