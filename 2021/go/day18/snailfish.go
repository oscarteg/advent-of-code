package main

import (
	"errors"
	"fmt"
	"slices"
	"strconv"
	"strings"
)

// explodeDepth is the depth of the regular numbers in a pair that is nested
// inside four pairs.
const explodeDepth = 5

var errInvalidNumber = errors.New("invalid snailfish number")

// element is a regular number together with the number of pairs it is nested in.
type element struct {
	value int
	depth int
}

// number is a snailfish number flattened to its regular numbers in reading
// order; the bracket structure follows from the depths alone.
type number []element

type parser struct {
	input string
	pos   int
}

func parseNumber(line string) (number, error) {
	p := parser{input: line}
	parsed, err := p.pair(1)
	if err != nil {
		return nil, fmt.Errorf("%q: %w", line, err)
	}
	if p.pos != len(p.input) {
		return nil, fmt.Errorf("%w: trailing input in %q", errInvalidNumber, line)
	}
	return parsed, nil
}

func (p *parser) pair(depth int) (number, error) {
	if err := p.expect('['); err != nil {
		return nil, err
	}
	left, err := p.element(depth)
	if err != nil {
		return nil, err
	}
	if err := p.expect(','); err != nil {
		return nil, err
	}
	right, err := p.element(depth)
	if err != nil {
		return nil, err
	}
	if err := p.expect(']'); err != nil {
		return nil, err
	}
	return append(left, right...), nil
}

func (p *parser) element(depth int) (number, error) {
	if p.pos < len(p.input) && p.input[p.pos] == '[' {
		return p.pair(depth + 1)
	}
	start := p.pos
	for p.pos < len(p.input) && isDigit(p.input[p.pos]) {
		p.pos++
	}
	if p.pos == start {
		return nil, fmt.Errorf("%w: expected a pair or a regular number at offset %d", errInvalidNumber, start)
	}
	value, err := strconv.Atoi(p.input[start:p.pos])
	if err != nil {
		return nil, fmt.Errorf("regular number: %w", err)
	}
	return number{{value: value, depth: depth}}, nil
}

func (p *parser) expect(want byte) error {
	if p.pos >= len(p.input) || p.input[p.pos] != want {
		return fmt.Errorf("%w: expected %q at offset %d", errInvalidNumber, want, p.pos)
	}
	p.pos++
	return nil
}

func isDigit(b byte) bool {
	return '0' <= b && b <= '9'
}

// add pairs both numbers and reduces the sum.
func (n number) add(other number) number {
	sum := make(number, 0, len(n)+len(other))
	for _, operand := range []number{n, other} {
		for _, e := range operand {
			sum = append(sum, element{value: e.value, depth: e.depth + 1})
		}
	}
	return sum.reduce()
}

// reduce explodes and splits until neither action applies. It rewrites the
// slice in place, so the caller must own the backing array.
func (n number) reduce() number {
	for {
		if exploded, ok := n.explode(); ok {
			n = exploded
			continue
		}
		split, ok := n.split()
		if !ok {
			return n
		}
		n = split
	}
}

// explode replaces the leftmost pair nested inside four pairs - the leftmost
// neighbours sharing a depth of five or more - with a zero one level up, after
// handing its values to the neighbouring regular numbers.
func (n number) explode() (number, bool) {
	for i := 0; i+1 < len(n); i++ {
		if n[i].depth < explodeDepth || n[i].depth != n[i+1].depth {
			continue
		}
		if i > 0 {
			n[i-1].value += n[i].value
		}
		if i+2 < len(n) {
			n[i+2].value += n[i+1].value
		}
		n[i] = element{value: 0, depth: n[i].depth - 1}
		return slices.Delete(n, i+1, i+2), true
	}
	return n, false
}

// split replaces the leftmost regular number of ten or more with a pair that
// rounds the left half down and the right half up.
func (n number) split() (number, bool) {
	for i, e := range n {
		if e.value < 10 {
			continue
		}
		n[i] = element{value: e.value / 2, depth: e.depth + 1}
		return slices.Insert(n, i+1, element{value: (e.value + 1) / 2, depth: e.depth + 1}), true
	}
	return n, false
}

// magnitude folds neighbours of equal depth into 3*left + 2*right until the
// whole number collapses into a single value.
func (n number) magnitude() (int, error) {
	stack := make([]element, 0, len(n))
	for _, e := range n {
		folded := e
		for len(stack) > 0 {
			top := stack[len(stack)-1]
			if top.depth != folded.depth {
				break
			}
			stack = stack[:len(stack)-1]
			folded = element{value: 3*top.value + 2*folded.value, depth: folded.depth - 1}
		}
		stack = append(stack, folded)
	}
	if len(stack) != 1 || stack[0].depth != 0 {
		return 0, fmt.Errorf("%w: not a single pair", errInvalidNumber)
	}
	return stack[0].value, nil
}

func parseNumbers(input string) ([]number, error) {
	var numbers []number
	for _, line := range strings.Split(input, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		parsed, err := parseNumber(line)
		if err != nil {
			return nil, err
		}
		numbers = append(numbers, parsed)
	}
	if len(numbers) == 0 {
		return nil, fmt.Errorf("%w: no snailfish numbers in the input", errInvalidNumber)
	}
	return numbers, nil
}

// solve returns the magnitude of the homework sum and the largest magnitude of
// any sum of two different numbers.
func solve(input string) (int, int, error) {
	numbers, err := parseNumbers(input)
	if err != nil {
		return 0, 0, err
	}
	if len(numbers) < 2 {
		return 0, 0, fmt.Errorf("%w: part two needs at least two snailfish numbers", errInvalidNumber)
	}
	sum := numbers[0]
	for _, next := range numbers[1:] {
		sum = sum.add(next)
	}
	homework, err := sum.magnitude()
	if err != nil {
		return 0, 0, err
	}
	best := 0
	for i, left := range numbers {
		for _, right := range numbers[i+1:] {
			for _, candidate := range []number{left.add(right), right.add(left)} {
				magnitude, err := candidate.magnitude()
				if err != nil {
					return 0, 0, err
				}
				best = max(best, magnitude)
			}
		}
	}
	return homework, best, nil
}
