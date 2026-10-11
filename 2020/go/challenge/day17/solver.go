package day17

import (
	"fmt"

	"github.com/oscarteg/advent-of-code-2020/challenge"
)

type pos3 struct {
	x, y, z int
}

type pos4 struct {
	x, y, z, w int
}

func activeCells(input *challenge.Input) ([][2]int, error) {
	var cells [][2]int
	for y, line := range input.LinesSlice() {
		for x, ch := range line {
			switch ch {
			case '#':
				cells = append(cells, [2]int{x, y})
			case '.':
			default:
				return nil, fmt.Errorf("line %d: unexpected %q, want '#' or '.'", y+1, ch)
			}
		}
	}
	return cells, nil
}

func parseInput3D(input *challenge.Input) (map[pos3]bool, error) {
	cells, err := activeCells(input)
	if err != nil {
		return nil, err
	}
	active := make(map[pos3]bool, len(cells))
	for _, c := range cells {
		active[pos3{c[0], c[1], 0}] = true
	}
	return active, nil
}

func parseInput4D(input *challenge.Input) (map[pos4]bool, error) {
	cells, err := activeCells(input)
	if err != nil {
		return nil, err
	}
	active := make(map[pos4]bool, len(cells))
	for _, c := range cells {
		active[pos4{c[0], c[1], 0, 0}] = true
	}
	return active, nil
}

func neighbors3D(p pos3) []pos3 {
	var result []pos3
	for dx := -1; dx <= 1; dx++ {
		for dy := -1; dy <= 1; dy++ {
			for dz := -1; dz <= 1; dz++ {
				if dx == 0 && dy == 0 && dz == 0 {
					continue
				}
				result = append(result, pos3{p.x + dx, p.y + dy, p.z + dz})
			}
		}
	}
	return result
}

func neighbors4D(p pos4) []pos4 {
	var result []pos4
	for dx := -1; dx <= 1; dx++ {
		for dy := -1; dy <= 1; dy++ {
			for dz := -1; dz <= 1; dz++ {
				for dw := -1; dw <= 1; dw++ {
					if dx == 0 && dy == 0 && dz == 0 && dw == 0 {
						continue
					}
					result = append(result, pos4{p.x + dx, p.y + dy, p.z + dz, p.w + dw})
				}
			}
		}
	}
	return result
}

func countNeighbors3D(p pos3, active map[pos3]bool) int {
	count := 0
	for _, n := range neighbors3D(p) {
		if active[n] {
			count++
		}
	}
	return count
}

func countNeighbors4D(p pos4, active map[pos4]bool) int {
	count := 0
	for _, n := range neighbors4D(p) {
		if active[n] {
			count++
		}
	}
	return count
}

func step3D(active map[pos3]bool) map[pos3]bool {
	toCheck := make(map[pos3]bool)
	for p := range active {
		toCheck[p] = true
		for _, n := range neighbors3D(p) {
			toCheck[n] = true
		}
	}

	newActive := make(map[pos3]bool)
	for p := range toCheck {
		count := countNeighbors3D(p, active)
		if active[p] && (count == 2 || count == 3) {
			newActive[p] = true
		} else if !active[p] && count == 3 {
			newActive[p] = true
		}
	}
	return newActive
}

func step4D(active map[pos4]bool) map[pos4]bool {
	toCheck := make(map[pos4]bool)
	for p := range active {
		toCheck[p] = true
		for _, n := range neighbors4D(p) {
			toCheck[n] = true
		}
	}

	newActive := make(map[pos4]bool)
	for p := range toCheck {
		count := countNeighbors4D(p, active)
		if active[p] && (count == 2 || count == 3) {
			newActive[p] = true
		} else if !active[p] && count == 3 {
			newActive[p] = true
		}
	}
	return newActive
}
