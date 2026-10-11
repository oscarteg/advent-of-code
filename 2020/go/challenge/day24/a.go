package day24

import (
	"fmt"

	"github.com/oscarteg/advent-of-code-2020/challenge"
)

type pos struct {
	x, y int
}

var directions = map[string]pos{
	"e":  {1, 0},
	"w":  {-1, 0},
	"ne": {0, 1},
	"sw": {0, -1},
	"se": {1, -1},
	"nw": {-1, 1},
}

func parsePath(line string) (pos, error) {
	p := pos{0, 0}
	i := 0
	for i < len(line) {
		ch := line[i]
		if ch == 'n' || ch == 's' {
			if i+1 >= len(line) {
				return pos{}, fmt.Errorf("incomplete direction after '%c'", ch)
			}
			direction := string([]byte{ch, line[i+1]})
			if d, ok := directions[direction]; ok {
				p.x += d.x
				p.y += d.y
				i += 2
			} else {
				return pos{}, fmt.Errorf("invalid direction: %s", direction)
			}
		} else if ch == 'e' || ch == 'w' {
			direction := string([]byte{ch})
			if d, ok := directions[direction]; ok {
				p.x += d.x
				p.y += d.y
				i++
			} else {
				return pos{}, fmt.Errorf("invalid direction: %s", direction)
			}
		} else {
			return pos{}, fmt.Errorf("invalid character: %c", ch)
		}
	}
	return p, nil
}

func A(input *challenge.Input) (int, error) {
	tiles := make(map[pos]int)

	for line := range input.Lines() {
		if line == "" {
			continue
		}

		p, err := parsePath(line)
		if err != nil {
			return 0, err
		}

		tiles[p]++
	}

	count := 0
	for _, flips := range tiles {
		if flips%2 == 1 {
			count++
		}
	}

	return count, nil
}
