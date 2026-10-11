package main

import (
	"fmt"
	"strconv"
	"strings"
)

type Instruction struct {
	op    string
	value int
}

func parseInput(input string) ([]Instruction, error) {
	var instructions []Instruction
	for _, line := range strings.Split(strings.TrimSpace(input), "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		if line == "noop" {
			instructions = append(instructions, Instruction{op: "noop"})
		} else if strings.HasPrefix(line, "addx ") {
			valueStr := line[5:]
			value, err := strconv.Atoi(valueStr)
			if err != nil {
				return nil, fmt.Errorf("invalid addx value: %w", err)
			}
			instructions = append(instructions, Instruction{op: "addx", value: value})
		} else {
			return nil, fmt.Errorf("invalid instruction: %s", line)
		}
	}
	return instructions, nil
}

func solve(input string) (int, string, error) {
	instructions, err := parseInput(input)
	if err != nil {
		return 0, "", err
	}

	xValues := []int{1}
	x := 1

	for _, instr := range instructions {
		switch instr.op {
		case "noop":
			xValues = append(xValues, x)
		case "addx":
			xValues = append(xValues, x)
			xValues = append(xValues, x)
			x += instr.value
		}
	}

	part1 := calculateSignalStrength(xValues)
	part2 := renderCRT(xValues)

	return part1, part2, nil
}

func calculateSignalStrength(xValues []int) int {
	cycles := []int{20, 60, 100, 140, 180, 220}
	strength := 0
	for _, cycle := range cycles {
		if cycle < len(xValues) {
			strength += cycle * xValues[cycle]
		}
	}
	return strength
}

func renderCRT(xValues []int) string {
	var result strings.Builder
	for row := 0; row < 6; row++ {
		if row > 0 {
			result.WriteByte('\n')
		}
		for col := 0; col < 40; col++ {
			cycle := row*40 + col + 1
			pixelX := col
			if cycle < len(xValues) {
				spriteX := xValues[cycle]
				if spriteX-1 <= pixelX && pixelX <= spriteX+1 {
					result.WriteByte('#')
				} else {
					result.WriteByte('.')
				}
			} else {
				result.WriteByte('.')
			}
		}
	}
	return result.String()
}
