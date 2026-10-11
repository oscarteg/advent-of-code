package main

import (
	"errors"
	"fmt"
	"slices"
	"strconv"
	"strings"
)

var errInvalidInput = errors.New("invalid input")

const (
	maxSwaps = 4
	maxBits  = 63
)

type operator uint8

const (
	opAnd operator = iota
	opOr
	opXor
)

type wireID int32

const (
	noWire wireID = -1
	noGate int    = -1
)

type gate struct {
	left, right wireID
	op          operator
}

type circuit struct {
	names   []string
	ids     map[string]wireID
	initial []int8   // per wire: -1 unknown, else 0 or 1
	driver  []int    // per wire: gate index, or noGate
	gates   []gate   // gate inputs, never swapped
	outputs []wireID // outputs[i] is the wire gate i drives
	bits    int      // number of x (and y) input wires
}

func (c *circuit) intern(name string) wireID {
	if id, ok := c.ids[name]; ok {
		return id
	}
	id := wireID(len(c.names))
	c.ids[name] = id
	c.names = append(c.names, name)
	c.initial = append(c.initial, -1)
	c.driver = append(c.driver, noGate)
	return id
}

func (c *circuit) wire(prefix byte, bit int) wireID {
	name := [3]byte{prefix, byte('0' + bit/10), byte('0' + bit%10)}
	if id, ok := c.ids[string(name[:])]; ok {
		return id
	}
	return noWire
}

func parseOperator(text string) (operator, error) {
	switch text {
	case "AND":
		return opAnd, nil
	case "OR":
		return opOr, nil
	case "XOR":
		return opXor, nil
	}
	return 0, fmt.Errorf("%w: unknown operator %q", errInvalidInput, text)
}

func (c *circuit) parseGate(line string) error {
	words := strings.Fields(line)
	if len(words) != 5 || words[3] != "->" {
		return fmt.Errorf("%w: malformed gate %q", errInvalidInput, line)
	}
	op, err := parseOperator(words[1])
	if err != nil {
		return fmt.Errorf("gate %q: %w", line, err)
	}

	left, right, output := c.intern(words[0]), c.intern(words[2]), c.intern(words[4])
	if c.driver[output] != noGate {
		return fmt.Errorf("%w: wire %q is driven twice", errInvalidInput, c.names[output])
	}
	c.driver[output] = len(c.gates)
	c.gates = append(c.gates, gate{left: left, right: right, op: op})
	c.outputs = append(c.outputs, output)
	return nil
}

func parse(input string) (*circuit, error) {
	c := &circuit{ids: make(map[string]wireID)}

	for _, raw := range strings.Split(input, "\n") {
		line := strings.TrimSpace(raw)
		switch {
		case line == "":
			continue
		case strings.Contains(line, "->"):
			if err := c.parseGate(line); err != nil {
				return nil, err
			}
		default:
			name, value, found := strings.Cut(line, ":")
			if !found {
				return nil, fmt.Errorf("%w: malformed line %q", errInvalidInput, line)
			}
			name, value = strings.TrimSpace(name), strings.TrimSpace(value)
			if name == "" || (value != "0" && value != "1") {
				return nil, fmt.Errorf("%w: malformed initial value %q", errInvalidInput, line)
			}
			if value == "1" {
				c.initial[c.intern(name)] = 1
			} else {
				c.initial[c.intern(name)] = 0
			}
		}
	}

	if !slices.ContainsFunc(c.outputs, func(wire wireID) bool {
		return strings.HasPrefix(c.names[wire], "z")
	}) {
		return nil, fmt.Errorf("%w: no z output wires", errInvalidInput)
	}
	for c.bits < maxBits && c.wire('x', c.bits) != noWire {
		c.bits++
	}
	return c, nil
}

type evaluator struct {
	circuit  *circuit
	values   []int8
	visiting []bool
}

func (e *evaluator) valueOf(wire wireID) (bool, error) {
	if e.values[wire] >= 0 {
		return e.values[wire] == 1, nil
	}
	driver := e.circuit.driver[wire]
	if driver == noGate {
		return false, fmt.Errorf("%w: wire %q never receives a value", errInvalidInput, e.circuit.names[wire])
	}
	if e.visiting[wire] {
		return false, fmt.Errorf("%w: circuit has a cycle at wire %q", errInvalidInput, e.circuit.names[wire])
	}

	e.visiting[wire] = true
	g := e.circuit.gates[driver]
	left, err := e.valueOf(g.left)
	if err != nil {
		return false, err
	}
	right, err := e.valueOf(g.right)
	if err != nil {
		return false, err
	}
	e.visiting[wire] = false

	var value bool
	switch g.op {
	case opAnd:
		value = left && right
	case opOr:
		value = left || right
	case opXor:
		value = left != right
	}
	e.values[wire] = 0
	if value {
		e.values[wire] = 1
	}
	return value, nil
}

func simulate(c *circuit) (int64, error) {
	e := &evaluator{circuit: c, values: slices.Clone(c.initial), visiting: make([]bool, len(c.names))}

	var result int64
	for _, wire := range c.outputs {
		name := c.names[wire]
		if !strings.HasPrefix(name, "z") {
			continue
		}
		bit, err := strconv.Atoi(name[1:])
		if err != nil {
			return 0, fmt.Errorf("bit index of wire %q: %w", name, err)
		}
		if bit > maxBits {
			return 0, fmt.Errorf("%w: wire %q is above bit %d", errInvalidInput, name, maxBits)
		}
		value, err := e.valueOf(wire)
		if err != nil {
			return 0, err
		}
		if value {
			result |= 1 << bit
		}
	}
	return result, nil
}

type gateKey struct {
	low, high wireID
	op        operator
}

func newGateKey(left, right wireID, op operator) gateKey {
	if left > right {
		left, right = right, left
	}
	return gateKey{low: left, high: right, op: op}
}

type adder struct {
	circuit *circuit
	lookup  map[gateKey]int
	outputs []wireID
	swapped []wireID
}

func newAdder(c *circuit) *adder {
	lookup := make(map[gateKey]int, len(c.gates))
	for index, g := range c.gates {
		key := newGateKey(g.left, g.right, g.op)
		if _, exists := lookup[key]; !exists {
			lookup[key] = index
		}
	}
	return &adder{circuit: c, lookup: lookup, outputs: slices.Clone(c.outputs)}
}

func (a *adder) driven(left, right wireID, op operator) wireID {
	if left == noWire || right == noWire {
		return noWire
	}
	index, ok := a.lookup[newGateKey(left, right, op)]
	if !ok {
		return noWire
	}
	return a.outputs[index]
}

// A correct ripple-carry adder computes z0 = x0 XOR y0 and carry0 = x0 AND y0, then for
// every higher bit z_i = (x_i XOR y_i) XOR carry_(i-1) and
// carry_i = (x_i AND y_i) OR ((x_i XOR y_i) AND carry_(i-1)), with the last carry on the
// top z wire. Returns the lowest bit that does not match, or bits + 1 when all of them do.
func (a *adder) firstBrokenBit() int {
	carry := noWire
	for bit := 0; bit < a.circuit.bits; bit++ {
		output := a.circuit.wire('z', bit)
		x, y := a.circuit.wire('x', bit), a.circuit.wire('y', bit)
		sum, generated := a.driven(x, y, opXor), a.driven(x, y, opAnd)
		if output == noWire || sum == noWire || generated == noWire {
			return bit
		}

		if bit == 0 {
			if sum != output {
				return 0
			}
			carry = generated
			continue
		}
		if a.driven(sum, carry, opXor) != output {
			return bit
		}
		propagated := a.driven(sum, carry, opAnd)
		if propagated == noWire {
			return bit
		}
		nextCarry := a.driven(propagated, generated, opOr)
		if nextCarry == noWire {
			return bit
		}
		carry = nextCarry
	}

	if top := a.circuit.wire('z', a.circuit.bits); top == noWire || carry != top {
		return a.circuit.bits
	}
	return a.circuit.bits + 1
}

// Depth-first search over output swaps, keeping only swaps that push the first broken bit
// higher; the puzzle input needs exactly four of them.
func (a *adder) repair(remaining int) bool {
	broken := a.firstBrokenBit()
	if broken > a.circuit.bits {
		return true
	}
	if remaining == 0 {
		return false
	}

	for first := range a.outputs {
		for second := first + 1; second < len(a.outputs); second++ {
			a.outputs[first], a.outputs[second] = a.outputs[second], a.outputs[first]
			if a.firstBrokenBit() > broken {
				a.swapped = append(a.swapped, a.outputs[first], a.outputs[second])
				if a.repair(remaining - 1) {
					return true
				}
				a.swapped = a.swapped[:len(a.swapped)-2]
			}
			a.outputs[first], a.outputs[second] = a.outputs[second], a.outputs[first]
		}
	}
	return false
}

// Empty when the circuit is not an adder that four swaps can repair, which includes both
// part 1 examples.
func findSwappedWires(c *circuit) string {
	a := newAdder(c)
	if !a.repair(maxSwaps) {
		return ""
	}
	names := make([]string, 0, len(a.swapped))
	for _, wire := range a.swapped {
		names = append(names, c.names[wire])
	}
	slices.Sort(names)
	return strings.Join(names, ",")
}

func solve(input string) (part1 int64, part2 string, err error) {
	c, err := parse(input)
	if err != nil {
		return 0, "", err
	}
	part1, err = simulate(c)
	if err != nil {
		return 0, "", err
	}
	return part1, findSwappedWires(c), nil
}
