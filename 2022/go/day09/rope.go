package main

import (
	"fmt"
	"strconv"
	"strings"
)

type point struct {
	x, y int
}

type motion struct {
	dx, dy int
	steps  int
}

func parseMotions(input string) ([]motion, error) {
	lines := strings.Split(strings.TrimSpace(input), "\n")
	motions := make([]motion, 0, len(lines))
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) != 2 {
			return nil, fmt.Errorf("malformed motion line %q: expected direction and steps", line)
		}
		steps, err := strconv.Atoi(fields[1])
		if err != nil {
			return nil, fmt.Errorf("invalid step count in line %q: %w", line, err)
		}
		var dx, dy int
		switch fields[0] {
		case "R":
			dx, dy = 1, 0
		case "L":
			dx, dy = -1, 0
		case "U":
			dx, dy = 0, 1
		case "D":
			dx, dy = 0, -1
		default:
			return nil, fmt.Errorf("unknown direction %q in line %q", fields[0], line)
		}
		motions = append(motions, motion{dx: dx, dy: dy, steps: steps})
	}
	if len(motions) == 0 {
		return nil, fmt.Errorf("input contains no motions")
	}
	return motions, nil
}

func sign(n int) int {
	switch {
	case n > 0:
		return 1
	case n < 0:
		return -1
	default:
		return 0
	}
}

func follow(head, tail point) point {
	dx, dy := head.x-tail.x, head.y-tail.y
	if abs(dx) <= 1 && abs(dy) <= 1 {
		return tail
	}
	return point{x: tail.x + sign(dx), y: tail.y + sign(dy)}
}

func abs(n int) int {
	if n < 0 {
		return -n
	}
	return n
}

func simulate(motions []motion, ropeLen int) int {
	knots := make([]point, ropeLen)
	visited := map[point]bool{knots[ropeLen-1]: true}

	for _, m := range motions {
		for i := 0; i < m.steps; i++ {
			knots[0].x += m.dx
			knots[0].y += m.dy
			for k := 1; k < ropeLen; k++ {
				knots[k] = follow(knots[k-1], knots[k])
			}
			visited[knots[ropeLen-1]] = true
		}
	}

	return len(visited)
}

func solve(input string) (part1, part2 int, err error) {
	motions, err := parseMotions(input)
	if err != nil {
		return 0, 0, err
	}
	return simulate(motions, 2), simulate(motions, 10), nil
}
