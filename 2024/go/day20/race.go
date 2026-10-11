package main

import (
	"errors"
	"fmt"
	"strings"
)

const realThreshold = 100

var errInvalidTrack = errors.New("invalid racetrack")

type point struct {
	x, y int
}

type track struct {
	width, height int
	wall          []bool
	start, end    point
}

func (t track) index(p point) int { return p.y*t.width + p.x }

func (t track) contains(p point) bool {
	return p.x >= 0 && p.x < t.width && p.y >= 0 && p.y < t.height
}

type race struct {
	track         track
	fromStart     []int
	toEnd         []int
	uncheatedTime int
}

func parseTrack(input string) (track, error) {
	var parsed track
	foundStart, foundEnd := false, false

	for _, raw := range strings.Split(strings.ReplaceAll(input, "\r\n", "\n"), "\n") {
		line := strings.TrimSpace(raw)
		if line == "" {
			continue
		}
		if parsed.width == 0 {
			parsed.width = len(line)
		} else if len(line) != parsed.width {
			return track{}, fmt.Errorf("%w: rows have different widths", errInvalidTrack)
		}

		for column := range len(line) {
			cell := line[column]
			here := point{x: column, y: parsed.height}
			switch cell {
			case '#', '.':
			case 'S':
				if foundStart {
					return track{}, fmt.Errorf("%w: more than one start position", errInvalidTrack)
				}
				parsed.start, foundStart = here, true
			case 'E':
				if foundEnd {
					return track{}, fmt.Errorf("%w: more than one end position", errInvalidTrack)
				}
				parsed.end, foundEnd = here, true
			default:
				return track{}, fmt.Errorf("%w: invalid map character %q", errInvalidTrack, cell)
			}
			parsed.wall = append(parsed.wall, cell == '#')
		}
		parsed.height++
	}

	switch {
	case parsed.height == 0:
		return track{}, fmt.Errorf("%w: no map rows", errInvalidTrack)
	case !foundStart:
		return track{}, fmt.Errorf("%w: no start position", errInvalidTrack)
	case !foundEnd:
		return track{}, fmt.Errorf("%w: no end position", errInvalidTrack)
	}
	return parsed, nil
}

// distancesFrom returns the steps from origin to every track cell reachable from it,
// and -1 for every cell it cannot reach.
func (t track) distancesFrom(origin point) []int {
	distance := make([]int, len(t.wall))
	for i := range distance {
		distance[i] = -1
	}

	deltas := [4]point{{1, 0}, {-1, 0}, {0, 1}, {0, -1}}
	distance[t.index(origin)] = 0
	frontier := []point{origin}

	for head := 0; head < len(frontier); head++ {
		current := frontier[head]
		steps := distance[t.index(current)]
		for _, delta := range deltas {
			next := point{x: current.x + delta.x, y: current.y + delta.y}
			if !t.contains(next) {
				continue
			}
			at := t.index(next)
			if t.wall[at] || distance[at] != -1 {
				continue
			}
			distance[at] = steps + 1
			frontier = append(frontier, next)
		}
	}
	return distance
}

func prepare(input string) (race, error) {
	parsed, err := parseTrack(input)
	if err != nil {
		return race{}, err
	}

	fromStart := parsed.distancesFrom(parsed.start)
	uncheatedTime := fromStart[parsed.index(parsed.end)]
	if uncheatedTime < 0 {
		return race{}, fmt.Errorf("%w: the end is unreachable from the start", errInvalidTrack)
	}

	return race{
		track:         parsed,
		fromStart:     fromStart,
		toEnd:         parsed.distancesFrom(parsed.end),
		uncheatedTime: uncheatedTime,
	}, nil
}

func abs(value int) int {
	if value < 0 {
		return -value
	}
	return value
}

// countCheats counts the cheats saving at least minSaving picoseconds. A cheat ignores
// walls for at most maxCheatTime picoseconds, so landing on a track cell at Manhattan
// distance d costs exactly d picoseconds: the cheated race then takes
// fromStart[begin] + d + toEnd[landing]. Cheats are identified by those two cells.
func countCheats(r race, maxCheatTime, minSaving int) int {
	cheats := 0
	for y := range r.track.height {
		for x := range r.track.width {
			toCheat := r.fromStart[r.track.index(point{x: x, y: y})]
			if toCheat < 0 {
				continue
			}
			for dy := -maxCheatTime; dy <= maxCheatTime; dy++ {
				room := maxCheatTime - abs(dy)
				for dx := -room; dx <= room; dx++ {
					landing := point{x: x + dx, y: y + dy}
					if !r.track.contains(landing) {
						continue
					}
					fromCheat := r.toEnd[r.track.index(landing)]
					if fromCheat < 0 {
						continue
					}
					jump := abs(dx) + abs(dy)
					if r.uncheatedTime-(toCheat+jump+fromCheat) >= minSaving {
						cheats++
					}
				}
			}
		}
	}
	return cheats
}

func solveWith(input string, minSaving int) (part1, part2 int, err error) {
	r, err := prepare(input)
	if err != nil {
		return 0, 0, err
	}
	return countCheats(r, 2, minSaving), countCheats(r, 20, minSaving), nil
}

func solve(input string) (part1, part2 int, err error) {
	return solveWith(input, realThreshold)
}
