package main

import (
	"errors"
	"fmt"
	"slices"
	"strconv"
	"strings"
)

const (
	rootName  = "root"
	humanName = "humn"
)

var errMonkeyMath = errors.New("invalid monkey riddle")

type operator byte

const (
	noOperator operator = 0
	add        operator = '+'
	subtract   operator = '-'
	multiply   operator = '*'
	divide     operator = '/'
)

func parseOperator(token string) (operator, error) {
	if len(token) == 1 {
		switch parsed := operator(token[0]); parsed {
		case add, subtract, multiply, divide:
			return parsed, nil
		}
	}
	return noOperator, fmt.Errorf("%w: unknown operator %q", errMonkeyMath, token)
}

func (op operator) apply(left, right int64) (int64, error) {
	switch op {
	case add:
		return left + right, nil
	case subtract:
		return left - right, nil
	case multiply:
		return left * right, nil
	case divide:
		if right == 0 {
			return 0, fmt.Errorf("%w: division by zero", errMonkeyMath)
		}
		return left / right, nil
	}
	return 0, fmt.Errorf("%w: unknown operator %q", errMonkeyMath, string(op))
}

// job is either a literal number or an operation on two other monkeys.
type job struct {
	number      int64
	left, right string
	operator    operator
}

func (j job) isNumber() bool { return j.operator == noOperator }

type monkeys map[string]job

func parseMonkeys(input string) (monkeys, error) {
	parsed := make(monkeys)
	for _, rawLine := range strings.Split(input, "\n") {
		line := strings.TrimSpace(rawLine)
		if line == "" {
			continue
		}
		name, rawJob, found := strings.Cut(line, ":")
		name = strings.TrimSpace(name)
		if !found || name == "" {
			return nil, fmt.Errorf("%w: malformed line %q", errMonkeyMath, line)
		}

		var parsedJob job
		fields := strings.Fields(rawJob)
		switch len(fields) {
		case 1:
			number, err := strconv.ParseInt(fields[0], 10, 64)
			if err != nil {
				return nil, fmt.Errorf("monkey %q number: %w", name, err)
			}
			parsedJob = job{number: number}
		case 3:
			parsedOperator, err := parseOperator(fields[1])
			if err != nil {
				return nil, fmt.Errorf("monkey %q: %w", name, err)
			}
			parsedJob = job{left: fields[0], operator: parsedOperator, right: fields[2]}
		default:
			return nil, fmt.Errorf("%w: malformed job %q", errMonkeyMath, line)
		}

		if _, duplicate := parsed[name]; duplicate {
			return nil, fmt.Errorf("%w: duplicate monkey %q", errMonkeyMath, name)
		}
		parsed[name] = parsedJob
	}
	if _, ok := parsed[rootName]; !ok {
		return nil, fmt.Errorf("%w: missing monkey %q", errMonkeyMath, rootName)
	}
	return parsed, nil
}

// evaluate yells name's number. budget bounds the depth by one step per monkey,
// so a riddle that refers back to itself fails instead of looping forever.
func (m monkeys) evaluate(name string, budget int) (int64, error) {
	if budget == 0 {
		return 0, fmt.Errorf("%w: cycle at monkey %q", errMonkeyMath, name)
	}
	current, ok := m[name]
	if !ok {
		return 0, fmt.Errorf("%w: unknown monkey %q", errMonkeyMath, name)
	}
	if current.isNumber() {
		return current.number, nil
	}
	left, err := m.evaluate(current.left, budget-1)
	if err != nil {
		return 0, err
	}
	right, err := m.evaluate(current.right, budget-1)
	if err != nil {
		return 0, err
	}
	return current.operator.apply(left, right)
}

// listeners maps every monkey to the monkey waiting for its number. Part 2
// inverts a single path, so a monkey heard by two monkeys is rejected.
func (m monkeys) listeners() (map[string]string, error) {
	heardBy := make(map[string]string, len(m))
	for name, current := range m {
		if current.isNumber() {
			continue
		}
		for _, operand := range []string{current.left, current.right} {
			if _, duplicate := heardBy[operand]; duplicate {
				return nil, fmt.Errorf("%w: monkey %q is heard by two monkeys", errMonkeyMath, operand)
			}
			heardBy[operand] = name
		}
	}
	return heardBy, nil
}

func (m monkeys) pathToHuman() ([]string, error) {
	heardBy, err := m.listeners()
	if err != nil {
		return nil, err
	}
	if _, ok := m[humanName]; !ok {
		return nil, fmt.Errorf("%w: missing monkey %q", errMonkeyMath, humanName)
	}

	path := []string{humanName}
	for current := humanName; current != rootName; {
		listener, ok := heardBy[current]
		if !ok {
			return nil, fmt.Errorf("%w: monkey %q is not heard by %q", errMonkeyMath, current, rootName)
		}
		path = append(path, listener)
		if len(path) > len(m) {
			return nil, fmt.Errorf("%w: cycle between %q and %q", errMonkeyMath, humanName, rootName)
		}
		current = listener
	}
	slices.Reverse(path)
	return path, nil
}

// invert solves target = left <operator> right for the operand on humn's side.
func invert(op operator, target, known int64, humanOnLeft bool) (int64, error) {
	switch op {
	case add:
		return target - known, nil
	case subtract:
		if humanOnLeft {
			return target + known, nil
		}
		return known - target, nil
	case multiply:
		return exactDivision(target, known)
	case divide:
		if humanOnLeft {
			return target * known, nil
		}
		return exactDivision(known, target)
	}
	return 0, fmt.Errorf("%w: unknown operator %q", errMonkeyMath, string(op))
}

// exactDivision refuses a remainder, which would leave the riddle without an
// integer answer.
func exactDivision(dividend, divisor int64) (int64, error) {
	if divisor == 0 || dividend%divisor != 0 {
		return 0, fmt.Errorf("%w: %d is not exactly divisible by %d", errMonkeyMath, dividend, divisor)
	}
	return dividend / divisor, nil
}

// humanValue walks from root down to humn, inverting every operation on the
// way: the branch without humn evaluates to a known number, which pins down
// what the branch holding humn has to yell.
func (m monkeys) humanValue() (int64, error) {
	path, err := m.pathToHuman()
	if err != nil {
		return 0, err
	}

	var target int64
	for index, parent := range path[:len(path)-1] {
		current := m[parent]
		if current.isNumber() {
			return 0, fmt.Errorf("%w: monkey %q must be an operation", errMonkeyMath, parent)
		}
		humanOnLeft := current.left == path[index+1]
		other := current.left
		if humanOnLeft {
			other = current.right
		}
		known, err := m.evaluate(other, len(m))
		if err != nil {
			return 0, err
		}
		// root compares its operands instead of combining them.
		if parent == rootName {
			target = known
			continue
		}
		if target, err = invert(current.operator, target, known, humanOnLeft); err != nil {
			return 0, err
		}
	}
	return target, nil
}

func solve(input string) (part1, part2 int64, err error) {
	parsed, err := parseMonkeys(input)
	if err != nil {
		return 0, 0, err
	}
	if part1, err = parsed.evaluate(rootName, len(parsed)); err != nil {
		return 0, 0, err
	}
	if part2, err = parsed.humanValue(); err != nil {
		return 0, 0, err
	}
	return part1, part2, nil
}
