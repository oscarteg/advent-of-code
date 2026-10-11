package main

import (
	"fmt"
	"regexp"
	"strconv"
)

type Machine struct {
	ax, ay, bx, by, px, py int64
}

func parseMachines(input string) ([]Machine, error) {
	var machines []Machine

	machineRe := regexp.MustCompile(`Button A: X\+(\d+), Y\+(\d+)\nButton B: X\+(\d+), Y\+(\d+)\nPrize: X=(\d+), Y=(\d+)`)

	matches := machineRe.FindAllStringSubmatch(input, -1)
	if len(matches) == 0 {
		return nil, fmt.Errorf("parse: no valid machines found")
	}

	for _, match := range matches {
		if len(match) != 7 {
			return nil, fmt.Errorf("parse: invalid machine format")
		}

		ax, _ := strconv.ParseInt(match[1], 10, 64)
		ay, _ := strconv.ParseInt(match[2], 10, 64)
		bx, _ := strconv.ParseInt(match[3], 10, 64)
		by, _ := strconv.ParseInt(match[4], 10, 64)
		px, _ := strconv.ParseInt(match[5], 10, 64)
		py, _ := strconv.ParseInt(match[6], 10, 64)

		machines = append(machines, Machine{ax, ay, bx, by, px, py})
	}

	return machines, nil
}

func solveMachinePart1(m Machine) int64 {
	for a := int64(0); a <= 100; a++ {
		for b := int64(0); b <= 100; b++ {
			if a*m.ax+b*m.bx == m.px && a*m.ay+b*m.by == m.py {
				return a*3 + b*1
			}
		}
	}
	return 0
}

func solveMachinePart2(m Machine) int64 {
	det := m.ax*m.by - m.ay*m.bx
	if det == 0 {
		return 0
	}

	numA := m.px*m.by - m.py*m.bx
	numB := m.ax*m.py - m.ay*m.px

	if numA%det != 0 || numB%det != 0 {
		return 0
	}

	a := numA / det
	b := numB / det

	if a >= 0 && b >= 0 && a*m.ax+b*m.bx == m.px && a*m.ay+b*m.by == m.py {
		return a*3 + b*1
	}

	return 0
}

func solve(input string) (int64, int64, error) {
	machines, err := parseMachines(input)
	if err != nil {
		return 0, 0, err
	}

	var part1, part2 int64

	for _, m := range machines {
		part1 += solveMachinePart1(m)

		m.px += 10000000000000
		m.py += 10000000000000
		part2 += solveMachinePart2(m)
	}

	return part1, part2, nil
}
