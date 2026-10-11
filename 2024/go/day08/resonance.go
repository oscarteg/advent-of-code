package main

import (
	"fmt"
	"strings"
)

func gcd(a, b int) int {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func solve(input string) (int, int, error) {
	lines := strings.Split(strings.TrimSpace(input), "\n")
	if len(lines) == 0 || len(lines[0]) == 0 {
		return 0, 0, fmt.Errorf("empty grid")
	}

	rows := len(lines)
	cols := len(lines[0])

	for _, line := range lines {
		if len(line) != cols {
			return 0, 0, fmt.Errorf("irregular grid dimensions")
		}
	}

	antennasByFreq := make(map[rune][]struct{ r, c int })
	for r, line := range lines {
		for c, ch := range line {
			if ch != '.' {
				antennasByFreq[ch] = append(antennasByFreq[ch], struct{ r, c int }{r, c})
			}
		}
	}

	antinodesPart1 := make(map[struct{ r, c int }]bool)
	antinodesPart2 := make(map[struct{ r, c int }]bool)

	for _, positions := range antennasByFreq {
		n := len(positions)
		for i := 0; i < n; i++ {
			for j := i + 1; j < n; j++ {
				r1, c1 := positions[i].r, positions[i].c
				r2, c2 := positions[j].r, positions[j].c

				dr := r2 - r1
				dc := c2 - c1

				antinode_r1 := r1 - dr
				antinode_c1 := c1 - dc
				if antinode_r1 >= 0 && antinode_r1 < rows && antinode_c1 >= 0 && antinode_c1 < cols {
					antinodesPart1[struct{ r, c int }{antinode_r1, antinode_c1}] = true
				}

				antinode_r2 := r2 + dr
				antinode_c2 := c2 + dc
				if antinode_r2 >= 0 && antinode_r2 < rows && antinode_c2 >= 0 && antinode_c2 < cols {
					antinodesPart1[struct{ r, c int }{antinode_r2, antinode_c2}] = true
				}

				if dr == 0 && dc == 0 {
					continue
				}

				g := gcd(abs(dr), abs(dc))
				stepR := dr / g
				stepC := dc / g

				for r, c := r1, c1; r >= 0 && r < rows && c >= 0 && c < cols; r, c = r+stepR, c+stepC {
					antinodesPart2[struct{ r, c int }{r, c}] = true
				}

				for r, c := r1-stepR, c1-stepC; r >= 0 && r < rows && c >= 0 && c < cols; r, c = r-stepR, c-stepC {
					antinodesPart2[struct{ r, c int }{r, c}] = true
				}
			}
		}
	}

	return len(antinodesPart1), len(antinodesPart2), nil
}
