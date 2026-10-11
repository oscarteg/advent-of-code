package main

import (
	"errors"
	"fmt"
	"strings"
)

var errInvalidGarden = errors.New("invalid garden map")

type plot struct {
	row, col int
}

// region is what one region costs to fence: part 1 prices area * perimeter,
// part 2 area * sides.
type region struct {
	area      int
	perimeter int
	sides     int
}

// garden is a fixed-width grid of plant letters, stored row-major.
type garden struct {
	plants []byte
	width  int
	height int
}

var (
	neighbourSteps = [4]plot{{-1, 0}, {1, 0}, {0, -1}, {0, 1}}
	diagonalSteps  = [2]int{-1, 1}
)

func parseGarden(input string) (garden, error) {
	var parsed garden
	for _, raw := range strings.Split(strings.ReplaceAll(input, "\r\n", "\n"), "\n") {
		line := strings.TrimSpace(raw)
		if line == "" {
			continue
		}
		if parsed.height == 0 {
			parsed.width = len(line)
		} else if len(line) != parsed.width {
			return garden{}, fmt.Errorf("%w: row %d is %d plots wide, want %d", errInvalidGarden, parsed.height, len(line), parsed.width)
		}
		for i := range line {
			if !isPlantLetter(line[i]) {
				return garden{}, fmt.Errorf("%w: plot %q is not a plant letter", errInvalidGarden, line[i])
			}
		}
		parsed.plants = append(parsed.plants, line...)
		parsed.height++
	}
	if parsed.height == 0 {
		return garden{}, fmt.Errorf("%w: no garden rows", errInvalidGarden)
	}
	return parsed, nil
}

func isPlantLetter(plant byte) bool {
	return plant >= 'a' && plant <= 'z' || plant >= 'A' && plant <= 'Z'
}

func (g garden) indexOf(at plot) int {
	return at.row*g.width + at.col
}

func (g garden) hasPlant(at plot, plant byte) bool {
	if at.row < 0 || at.col < 0 || at.row >= g.height || at.col >= g.width {
		return false
	}
	return g.plants[g.indexOf(at)] == plant
}

// cornersAt counts the corners the plot contributes to its region: a corner is
// convex when both plots adjacent to a diagonal lie outside the region, and
// concave when both lie inside while the diagonal plot itself does not.
func (g garden) cornersAt(at plot, plant byte) int {
	corners := 0
	for _, rowStep := range diagonalSteps {
		for _, colStep := range diagonalSteps {
			vertical := g.hasPlant(plot{at.row + rowStep, at.col}, plant)
			horizontal := g.hasPlant(plot{at.row, at.col + colStep}, plant)
			diagonal := g.hasPlant(plot{at.row + rowStep, at.col + colStep}, plant)
			if !vertical && !horizontal {
				corners++
			}
			if vertical && horizontal && !diagonal {
				corners++
			}
		}
	}
	return corners
}

// measureRegion flood-fills the region containing start, marking every plot it
// reaches visited. A region has as many sides as corners.
func (g garden) measureRegion(start plot, visited []bool) region {
	plant := g.plants[g.indexOf(start)]
	visited[g.indexOf(start)] = true
	queue := []plot{start}

	var measured region
	for len(queue) > 0 {
		current := queue[len(queue)-1]
		queue = queue[:len(queue)-1]
		measured.area++
		measured.sides += g.cornersAt(current, plant)
		for _, step := range neighbourSteps {
			next := plot{current.row + step.row, current.col + step.col}
			if !g.hasPlant(next, plant) {
				measured.perimeter++
				continue
			}
			if visited[g.indexOf(next)] {
				continue
			}
			visited[g.indexOf(next)] = true
			queue = append(queue, next)
		}
	}
	return measured
}

func solve(input string) (part1, part2 int, err error) {
	parsed, err := parseGarden(input)
	if err != nil {
		return 0, 0, err
	}

	visited := make([]bool, len(parsed.plants))
	for row := range parsed.height {
		for col := range parsed.width {
			at := plot{row, col}
			if visited[parsed.indexOf(at)] {
				continue
			}
			measured := parsed.measureRegion(at, visited)
			part1 += measured.area * measured.perimeter
			part2 += measured.area * measured.sides
		}
	}
	return part1, part2, nil
}
