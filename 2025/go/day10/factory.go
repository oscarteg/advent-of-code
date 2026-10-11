// Package main solves Advent of Code 2025 day 10, Factory: every line
// describes one machine with an indicator light diagram in [], button wiring
// schematics in () and joltage requirements in {}.
package main

import (
	"errors"
	"fmt"
	"math/bits"
	"strconv"
	"strings"
)

const (
	// maxLights is the width of the bitmask holding the diagram.
	maxLights = 64
	// maxButtons bounds the button subsets enumerated as uint64 bitmasks.
	maxButtons = 63
)

var (
	errInvalidMachine = errors.New("invalid machine")
	errUnreachable    = errors.New("machine cannot be configured")
)

// machine is one line of the manual.
type machine struct {
	lights int
	// target has bit j set when light j has to end up on.
	target uint64
	// buttons lists, per button, the lights it toggles and the counters it
	// increments.
	buttons [][]int
	joltage []int
}

// unwrap strips the brackets off a field like [.##.], (0,2) or {3,5}.
func unwrap(field string, open, close byte) (string, error) {
	if len(field) < 3 || field[0] != open || field[len(field)-1] != close {
		return "", fmt.Errorf("%w: %q is not %c...%c", errInvalidMachine, field, open, close)
	}
	return field[1 : len(field)-1], nil
}

func parseIndices(wiring string, lights int) ([]int, error) {
	parts := strings.Split(wiring, ",")
	indices := make([]int, 0, len(parts))
	for _, part := range parts {
		index, err := strconv.Atoi(part)
		if err != nil {
			return nil, fmt.Errorf("button wiring: %w", err)
		}
		if index < 0 || index >= lights {
			return nil, fmt.Errorf("%w: button wired to light %d of %d", errInvalidMachine, index, lights)
		}
		indices = append(indices, index)
	}
	return indices, nil
}

func parseMachine(line string) (machine, error) {
	fields := strings.Fields(line)
	if len(fields) < 3 {
		return machine{}, fmt.Errorf("%w: %q needs a diagram, a button and joltage", errInvalidMachine, line)
	}

	diagram, err := unwrap(fields[0], '[', ']')
	if err != nil {
		return machine{}, err
	}
	if len(diagram) > maxLights {
		return machine{}, fmt.Errorf("%w: %d lights exceed the supported %d", errInvalidMachine, len(diagram), maxLights)
	}
	var target uint64
	for light := range len(diagram) {
		switch diagram[light] {
		case '.':
		case '#':
			target |= 1 << light
		default:
			return machine{}, fmt.Errorf("%w: unknown light state %q", errInvalidMachine, diagram[light])
		}
	}

	buttonFields := fields[1 : len(fields)-1]
	if len(buttonFields) > maxButtons {
		return machine{}, fmt.Errorf("%w: %d buttons exceed the supported %d", errInvalidMachine, len(buttonFields), maxButtons)
	}
	buttons := make([][]int, 0, len(buttonFields))
	for _, field := range buttonFields {
		wiring, err := unwrap(field, '(', ')')
		if err != nil {
			return machine{}, err
		}
		indices, err := parseIndices(wiring, len(diagram))
		if err != nil {
			return machine{}, err
		}
		buttons = append(buttons, indices)
	}

	requirements, err := unwrap(fields[len(fields)-1], '{', '}')
	if err != nil {
		return machine{}, err
	}
	joltage := make([]int, 0, len(diagram))
	for _, requirement := range strings.Split(requirements, ",") {
		value, err := strconv.Atoi(requirement)
		if err != nil {
			return machine{}, fmt.Errorf("joltage requirement: %w", err)
		}
		if value < 0 {
			return machine{}, fmt.Errorf("%w: negative joltage %d", errInvalidMachine, value)
		}
		joltage = append(joltage, value)
	}
	if len(joltage) != len(diagram) {
		return machine{}, fmt.Errorf("%w: %d joltage requirements for %d lights", errInvalidMachine, len(joltage), len(diagram))
	}

	return machine{lights: len(diagram), target: target, buttons: buttons, joltage: joltage}, nil
}

// fewestToggles returns the fewest presses that toggle the all-off lights
// into the diagram. Pressing a button twice cancels out, so every button is
// pressed once or not at all and the lights form a system over GF(2).
// Elimination yields one solution plus a basis of the null space - button
// sets that toggle nothing - and the answer is the lightest solution in that
// coset.
func fewestToggles(m machine) (int, error) {
	type row struct {
		lights  uint64
		buttons uint64
	}
	rows := make([]row, len(m.buttons))
	for index, button := range m.buttons {
		var lights uint64
		for _, light := range button {
			lights |= 1 << light
		}
		rows[index] = row{lights: lights, buttons: 1 << index}
	}

	pivotLights := make([]int, 0, len(rows))
	rank := 0
	for light := 0; light < m.lights && rank < len(rows); light++ {
		bit := uint64(1) << light
		pivot := -1
		for index := rank; index < len(rows); index++ {
			if rows[index].lights&bit != 0 {
				pivot = index
				break
			}
		}
		if pivot < 0 {
			continue
		}
		rows[rank], rows[pivot] = rows[pivot], rows[rank]
		for index := range rows {
			if index == rank || rows[index].lights&bit == 0 {
				continue
			}
			rows[index].lights ^= rows[rank].lights
			rows[index].buttons ^= rows[rank].buttons
		}
		pivotLights = append(pivotLights, light)
		rank++
	}

	residual, pressed := m.target, uint64(0)
	for index, light := range pivotLights {
		if residual&(1<<light) == 0 {
			continue
		}
		residual ^= rows[index].lights
		pressed ^= rows[index].buttons
	}
	if residual != 0 {
		return 0, fmt.Errorf("%w: no button set toggles the lights into the diagram", errUnreachable)
	}

	basis := rows[rank:]
	best := bits.OnesCount64(pressed)
	for combination := uint64(1); combination < uint64(1)<<len(basis); combination++ {
		buttons := pressed
		for index := range basis {
			if combination&(1<<index) != 0 {
				buttons ^= basis[index].buttons
			}
		}
		if count := bits.OnesCount64(buttons); count < best {
			best = count
		}
	}
	return best, nil
}

func gcd(a, b int) int {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

func rowGCD(row []int) int {
	divisor := 0
	for _, cell := range row {
		if cell < 0 {
			cell = -cell
		}
		divisor = gcd(divisor, cell)
	}
	return divisor
}

func allZero(row []int) bool {
	for _, cell := range row {
		if cell != 0 {
			return false
		}
	}
	return true
}

// fewestPresses returns the fewest presses that make every counter reach its
// joltage target exactly. A press adds one to each counter its button feeds,
// so the targets are a linear system A x = joltage to be solved over the
// nonnegative integers. Integer-preserving Gaussian elimination turns the
// pivot presses into a function of the free ones; each free press is bounded
// by the smallest target its button feeds, so a depth-first search over that
// box - pruned by the best total so far - finds the cheapest solution.
func fewestPresses(m machine) (int, error) {
	buttons := len(m.buttons)
	// One row per counter; the last column holds that counter's target.
	rows := make([][]int, m.lights)
	for light := range rows {
		rows[light] = make([]int, buttons+1)
		rows[light][buttons] = m.joltage[light]
	}
	for column, button := range m.buttons {
		for _, light := range button {
			rows[light][column] = 1
		}
	}

	pivotRow := make([]int, buttons)
	for column := range pivotRow {
		pivotRow[column] = -1
	}
	rank := 0
	for column := 0; column < buttons && rank < m.lights; column++ {
		pivot := -1
		for row := rank; row < m.lights; row++ {
			if rows[row][column] != 0 {
				pivot = row
				break
			}
		}
		if pivot < 0 {
			continue
		}
		rows[rank], rows[pivot] = rows[pivot], rows[rank]
		pivotValue := rows[rank][column]
		for row := range rows {
			if row == rank || rows[row][column] == 0 {
				continue
			}
			factor := rows[row][column]
			for index, pivotCell := range rows[rank] {
				rows[row][index] = rows[row][index]*pivotValue - pivotCell*factor
			}
			// Keep the entries small; the row stays an equivalent equation.
			if divisor := rowGCD(rows[row]); divisor > 1 {
				for index := range rows[row] {
					rows[row][index] /= divisor
				}
			}
		}
		pivotRow[column] = rank
		rank++
	}

	for _, row := range rows {
		if allZero(row[:buttons]) && row[buttons] != 0 {
			return 0, fmt.Errorf("%w: joltage %v is unreachable", errUnreachable, m.joltage)
		}
	}

	var pivotColumns, freeColumns, bounds []int
	for column, row := range pivotRow {
		if row >= 0 {
			pivotColumns = append(pivotColumns, column)
			continue
		}
		// A button cannot be pressed more often than the smallest target it feeds.
		bound := m.joltage[m.buttons[column][0]]
		for _, light := range m.buttons[column][1:] {
			if m.joltage[light] < bound {
				bound = m.joltage[light]
			}
		}
		freeColumns = append(freeColumns, column)
		bounds = append(bounds, bound)
	}

	presses := make([]int, buttons)
	best := -1
	var search func(index, assigned int)
	search = func(index, assigned int) {
		if best >= 0 && assigned >= best {
			return
		}
		if index == len(freeColumns) {
			total := assigned
			for _, column := range pivotColumns {
				row := rows[pivotRow[column]]
				remainder := row[buttons]
				for _, free := range freeColumns {
					remainder -= row[free] * presses[free]
				}
				if remainder%row[column] != 0 {
					return
				}
				count := remainder / row[column]
				if count < 0 {
					return
				}
				total += count
			}
			if best < 0 || total < best {
				best = total
			}
			return
		}
		column := freeColumns[index]
		for count := 0; count <= bounds[index]; count++ {
			presses[column] = count
			search(index+1, assigned+count)
		}
		presses[column] = 0
	}
	search(0, 0)
	if best < 0 {
		return 0, fmt.Errorf("%w: joltage %v is unreachable", errUnreachable, m.joltage)
	}
	return best, nil
}

func solve(input string) (part1, part2 int, err error) {
	machines := 0
	for number, line := range strings.Split(strings.ReplaceAll(input, "\r\n", "\n"), "\n") {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			continue
		}
		parsed, err := parseMachine(trimmed)
		if err != nil {
			return 0, 0, fmt.Errorf("line %d: %w", number+1, err)
		}
		toggles, err := fewestToggles(parsed)
		if err != nil {
			return 0, 0, fmt.Errorf("line %d: %w", number+1, err)
		}
		presses, err := fewestPresses(parsed)
		if err != nil {
			return 0, 0, fmt.Errorf("line %d: %w", number+1, err)
		}
		part1 += toggles
		part2 += presses
		machines++
	}
	if machines == 0 {
		return 0, 0, fmt.Errorf("%w: input has no machines", errInvalidMachine)
	}
	return part1, part2, nil
}
