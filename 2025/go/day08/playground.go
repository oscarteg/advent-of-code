// Package main solves Advent of Code 2025 day 8, "Playground".
package main

import (
	"cmp"
	"errors"
	"fmt"
	"slices"
	"strconv"
	"strings"
)

var (
	errInvalidBox     = errors.New("invalid junction box")
	errNotEnoughBoxes = errors.New("need at least two junction boxes")
	errTooFewCircuits = errors.New("fewer than three circuits left")
)

// The puzzle connects the 1000 closest pairs of the real input and the 10
// closest of its 20-box example, so the count follows the input size.
const (
	realConnections    = 1000
	exampleBoxes       = 20
	exampleConnections = 10
)

type box struct {
	x, y, z int
}

func (b box) distanceSquared(other box) int {
	dx, dy, dz := b.x-other.x, b.y-other.y, b.z-other.z
	return dx*dx + dy*dy + dz*dz
}

// Squared distances order the same way straight-line distances do, so a
// connection never needs a square root.
type connection struct {
	first, second   int
	distanceSquared int
}

func parseBoxes(input string) ([]box, error) {
	var boxes []box
	for _, line := range strings.Split(input, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		fields := strings.Split(line, ",")
		if len(fields) != 3 {
			return nil, fmt.Errorf("%w %q: got %d coordinates, want 3", errInvalidBox, line, len(fields))
		}
		var coordinates [3]int
		for i, field := range fields {
			coordinate, err := strconv.Atoi(strings.TrimSpace(field))
			if err != nil {
				return nil, fmt.Errorf("%w %q: %w", errInvalidBox, line, err)
			}
			coordinates[i] = coordinate
		}

		boxes = append(boxes, box{x: coordinates[0], y: coordinates[1], z: coordinates[2]})
	}
	if len(boxes) < 2 {
		return nil, fmt.Errorf("%w: got %d", errNotEnoughBoxes, len(boxes))
	}

	return boxes, nil
}

// sortedConnections returns every pair of boxes, closest first; equal
// distances fall back to the box order, keeping ties deterministic.
func sortedConnections(boxes []box) []connection {
	connections := make([]connection, 0, len(boxes)*(len(boxes)-1)/2)
	for first, current := range boxes {
		for second := first + 1; second < len(boxes); second++ {
			connections = append(connections, connection{
				first:           first,
				second:          second,
				distanceSquared: current.distanceSquared(boxes[second]),
			})
		}
	}
	slices.SortFunc(connections, func(a, b connection) int {
		return cmp.Or(
			cmp.Compare(a.distanceSquared, b.distanceSquared),
			cmp.Compare(a.first, b.first),
			cmp.Compare(a.second, b.second),
		)
	})
	return connections
}

// circuits tracks which boxes share a circuit, as a union-find forest
// merged by size.
type circuits struct {
	parent []int
	size   []int
	count  int
}

func newCircuits(boxes int) *circuits {
	c := &circuits{parent: make([]int, boxes), size: make([]int, boxes), count: boxes}
	for i := range c.parent {
		c.parent[i] = i
		c.size[i] = 1
	}
	return c
}

func (c *circuits) find(box int) int {
	for c.parent[box] != box {
		c.parent[box] = c.parent[c.parent[box]] // path halving
		box = c.parent[box]
	}
	return box
}

// connect joins two boxes, reporting whether this merged two circuits.
func (c *circuits) connect(first, second int) bool {
	root, merged := c.find(first), c.find(second)
	if root == merged {
		return false
	}
	if c.size[root] < c.size[merged] {
		root, merged = merged, root
	}
	c.parent[merged] = root
	c.size[root] += c.size[merged]
	c.count--
	return true
}

func (c *circuits) threeLargestProduct() (int, error) {
	if c.count < 3 {
		return 0, fmt.Errorf("%w: %d", errTooFewCircuits, c.count)
	}
	var largest [3]int
	for box, parent := range c.parent {
		if parent != box {
			continue
		}
		candidate := c.size[box]
		for i, size := range largest {
			if candidate > size {
				largest[i], candidate = candidate, size
			}
		}
	}
	return largest[0] * largest[1] * largest[2], nil
}

// solve reports the product of the three largest circuit sizes after the
// prescribed connections (part1) and the product of the X coordinates of
// the last pair that merged two circuits (part2).
func solve(input string) (part1, part2 int, err error) {
	return solveWith(input, 0)
}

// solveWith overrides the number of connections part 1 makes; a count of 0
// follows the number the puzzle prescribes for this input size.
func solveWith(input string, connectionCount int) (part1, part2 int, err error) {
	boxes, err := parseBoxes(input)
	if err != nil {
		return 0, 0, err
	}
	if connectionCount == 0 {
		connectionCount = realConnections
		if len(boxes) <= exampleBoxes {
			connectionCount = exampleConnections
		}
	}

	connections := sortedConnections(boxes)
	circuits := newCircuits(len(boxes))

	// The first connection always merges, so it is a safe starting value.
	lastMerge := connections[0]
	havePart1 := false
	for index, candidate := range connections {
		if index == connectionCount {
			if part1, err = circuits.threeLargestProduct(); err != nil {
				return 0, 0, err
			}
			havePart1 = true
		}
		if circuits.connect(candidate.first, candidate.second) {
			lastMerge = candidate
		}
	}
	if !havePart1 {
		if part1, err = circuits.threeLargestProduct(); err != nil {
			return 0, 0, err
		}
	}

	return part1, boxes[lastMerge.first].x * boxes[lastMerge.second].x, nil
}
