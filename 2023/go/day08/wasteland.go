package main

import (
	"errors"
	"fmt"
	"math"
	"strings"
)

var errInvalidDocuments = errors.New("invalid documents")

const (
	camelStart       = "AAA"
	camelDestination = "ZZZ"
	ghostStart       = "A"
	ghostDestination = "Z"
)

type instruction byte

const (
	goLeft  instruction = 'L'
	goRight instruction = 'R'
)

type node struct {
	left, right string
}

type documents struct {
	instructions []instruction
	nodes        map[string]node
}

func parse(input string) (documents, error) {
	var docs documents
	var lines []string
	for _, line := range strings.Split(input, "\n") {
		if trimmed := strings.TrimSpace(line); trimmed != "" {
			lines = append(lines, trimmed)
		}
	}
	if len(lines) == 0 {
		return docs, fmt.Errorf("%w: missing instructions", errInvalidDocuments)
	}
	for _, step := range []byte(lines[0]) {
		if instruction(step) != goLeft && instruction(step) != goRight {
			return docs, fmt.Errorf("%w: instruction %q is not L or R", errInvalidDocuments, step)
		}
		docs.instructions = append(docs.instructions, instruction(step))
	}
	docs.nodes = make(map[string]node, len(lines)-1)
	for _, line := range lines[1:] {
		name, targets, assigned := strings.Cut(line, " = ")
		if !assigned || !strings.HasPrefix(targets, "(") || !strings.HasSuffix(targets, ")") {
			return docs, fmt.Errorf("%w: %q is not `name = (left, right)`", errInvalidDocuments, line)
		}
		left, right, paired := strings.Cut(strings.Trim(targets, "()"), ", ")
		if !paired {
			return docs, fmt.Errorf("%w: %q needs a left and a right target", errInvalidDocuments, line)
		}
		if _, duplicate := docs.nodes[name]; duplicate {
			return docs, fmt.Errorf("%w: node %q is listed twice", errInvalidDocuments, name)
		}
		docs.nodes[name] = node{left: left, right: right}
	}
	if len(docs.nodes) == 0 {
		return docs, fmt.Errorf("%w: missing nodes", errInvalidDocuments)
	}
	for name, current := range docs.nodes {
		for _, target := range []string{current.left, current.right} {
			if _, known := docs.nodes[target]; !known {
				return docs, fmt.Errorf("%w: node %q points at unknown node %q", errInvalidDocuments, name, target)
			}
		}
	}
	return docs, nil
}

func (d documents) steps(start string, arrived func(string) bool) (uint64, error) {
	if _, known := d.nodes[start]; !known {
		return 0, fmt.Errorf("%w: unknown start node %q", errInvalidDocuments, start)
	}
	current := start
	// Within len(nodes)*len(instructions) steps the (node, instruction) state
	// repeats, so a longer walk can never arrive.
	for steps := 0; steps < len(d.nodes)*len(d.instructions); steps++ {
		if arrived(current) {
			return uint64(steps), nil
		}
		if d.instructions[steps%len(d.instructions)] == goLeft {
			current = d.nodes[current].left
		} else {
			current = d.nodes[current].right
		}
	}
	return 0, fmt.Errorf("%w: no destination is reachable from %q", errInvalidDocuments, start)
}

func (d documents) camelSteps() (uint64, error) {
	return d.steps(camelStart, func(name string) bool { return name == camelDestination })
}

// ghostSteps relies on every ghost revisiting its destination once per walk
// length, so the first step they share is the least common multiple.
func (d documents) ghostSteps() (uint64, error) {
	starts := make([]string, 0, len(d.nodes))
	for name := range d.nodes {
		if strings.HasSuffix(name, ghostStart) {
			starts = append(starts, name)
		}
	}
	if len(starts) == 0 {
		return 0, fmt.Errorf("%w: no node ends in %q", errInvalidDocuments, ghostStart)
	}
	combined := uint64(1)
	for _, start := range starts {
		walk, err := d.steps(start, func(name string) bool {
			return strings.HasSuffix(name, ghostDestination)
		})
		if err != nil {
			return 0, err
		}
		if combined, err = leastCommonMultiple(combined, walk); err != nil {
			return 0, err
		}
	}
	return combined, nil
}

func leastCommonMultiple(a, b uint64) (uint64, error) {
	reduced := a / greatestCommonDivisor(a, b)
	if b > 0 && reduced > math.MaxUint64/b {
		return 0, fmt.Errorf("%w: step count overflows", errInvalidDocuments)
	}
	return reduced * b, nil
}

func greatestCommonDivisor(a, b uint64) uint64 {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

func solve(input string) (uint64, uint64, error) {
	docs, err := parse(input)
	if err != nil {
		return 0, 0, err
	}
	camel, err := docs.camelSteps()
	if err != nil {
		return 0, 0, err
	}
	ghosts, err := docs.ghostSteps()
	if err != nil {
		return 0, 0, err
	}
	return camel, ghosts, nil
}
