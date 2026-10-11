package main

import (
	"fmt"
	"strings"
	"unicode"
)

type Value struct {
	kind  string // "number", "symbol", "empty"
	digit uint8
	char  rune
}

func solve(input string) (int, int, error) {
	lines := strings.Split(strings.TrimSpace(input), "\n")
	if len(lines) == 0 {
		return 0, 0, nil
	}

	// Build grid
	grid := make(map[Pos]Value)
	for y, line := range lines {
		if len(line) != len(lines[0]) {
			return 0, 0, fmt.Errorf("line %d has length %d, want %d", y+1, len(line), len(lines[0]))
		}
		for x, ch := range line {
			pos := Pos{x, y}
			if ch == '.' {
				grid[pos] = Value{kind: "empty"}
			} else if unicode.IsDigit(ch) {
				grid[pos] = Value{kind: "number", digit: uint8(ch - '0')}
			} else {
				grid[pos] = Value{kind: "symbol", char: ch}
			}
		}
	}

	part1 := solvePart1(grid)
	part2 := solvePart2(grid)
	return part1, part2, nil
}

type Pos struct {
	x, y int
}

func (p Pos) adjacent() []Pos {
	return []Pos{
		{p.x - 1, p.y - 1}, {p.x, p.y - 1}, {p.x + 1, p.y - 1},
		{p.x - 1, p.y}, {p.x + 1, p.y},
		{p.x - 1, p.y + 1}, {p.x, p.y + 1}, {p.x + 1, p.y + 1},
	}
}

func solvePart1(grid map[Pos]Value) int {
	sum := 0
	processed := make(map[Pos]bool)

	for pos, val := range grid {
		if val.kind != "symbol" {
			continue
		}

		for _, adj := range pos.adjacent() {
			if grid[adj].kind == "number" && !processed[adj] {
				// Find the start of the number
				start := adj
				for start.x > 0 && grid[Pos{start.x - 1, start.y}].kind == "number" {
					start.x--
				}

				// Get the full number
				num := 0
				for p := start; grid[p].kind == "number"; p.x++ {
					num = num*10 + int(grid[p].digit)
					processed[p] = true
				}

				sum += num
				processed[adj] = true
			}
		}
	}

	return sum
}

func solvePart2(grid map[Pos]Value) int {
	sum := 0

	for pos, val := range grid {
		if val.kind != "symbol" || val.char != '*' {
			continue
		}

		var numbers []int
		processed := make(map[Pos]bool)

		for _, adj := range pos.adjacent() {
			if grid[adj].kind == "number" && !processed[adj] {
				// Find the start of the number
				start := adj
				for start.x > 0 && grid[Pos{start.x - 1, start.y}].kind == "number" {
					start.x--
				}

				// Get the full number
				num := 0
				for p := start; grid[p].kind == "number"; p.x++ {
					num = num*10 + int(grid[p].digit)
					processed[p] = true
				}

				numbers = append(numbers, num)
				processed[adj] = true
			}
		}

		if len(numbers) == 2 {
			sum += numbers[0] * numbers[1]
		}
	}

	return sum
}
