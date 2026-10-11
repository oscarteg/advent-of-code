package main

import (
	"cmp"
	"errors"
	"fmt"
	"slices"
	"strconv"
	"strings"
)

var errInvalidSnapshot = errors.New("invalid snapshot")

// brick holds inclusive bounds per axis.
type brick struct {
	min, max [3]int
}

type column struct {
	x, y int
}

type stack struct {
	// supports[i] holds the bricks resting directly on brick i.
	supports [][]int
	// supportedBy[i] holds the bricks directly underneath brick i.
	supportedBy [][]int
}

func parseCorner(input string) ([3]int, error) {
	var corner [3]int
	fields := strings.Split(strings.TrimSpace(input), ",")
	if len(fields) != len(corner) {
		return corner, fmt.Errorf("%w: corner needs three coordinates", errInvalidSnapshot)
	}
	for axis, field := range fields {
		value, err := strconv.Atoi(field)
		if err != nil {
			return corner, fmt.Errorf("coordinate: %w", err)
		}
		if value < 0 {
			return corner, fmt.Errorf("%w: coordinates are not negative", errInvalidSnapshot)
		}
		corner[axis] = value
	}
	return corner, nil
}

func parseBrick(line string) (brick, error) {
	start, end, found := strings.Cut(line, "~")
	if !found {
		return brick{}, fmt.Errorf("%w: brick needs two ends separated by ~", errInvalidSnapshot)
	}
	from, err := parseCorner(start)
	if err != nil {
		return brick{}, err
	}
	to, err := parseCorner(end)
	if err != nil {
		return brick{}, err
	}
	var result brick
	axes := 0
	for axis := range from {
		result.min[axis] = min(from[axis], to[axis])
		result.max[axis] = max(from[axis], to[axis])
		if result.min[axis] != result.max[axis] {
			axes++
		}
	}
	if result.min[2] < 1 {
		return brick{}, fmt.Errorf("%w: bricks rest at z >= 1", errInvalidSnapshot)
	}
	if axes > 1 {
		return brick{}, fmt.Errorf("%w: brick must be a straight line of cubes", errInvalidSnapshot)
	}
	return result, nil
}

func (b brick) footprint() []column {
	columns := make([]column, 0, (b.max[0]-b.min[0]+1)*(b.max[1]-b.min[1]+1))
	for x := b.min[0]; x <= b.max[0]; x++ {
		for y := b.min[1]; y <= b.max[1]; y++ {
			columns = append(columns, column{x: x, y: y})
		}
	}
	return columns
}

func (b brick) height() int {
	return b.max[2] - b.min[2] + 1
}

func parseBricks(input string) ([]brick, error) {
	var bricks []brick
	for _, line := range strings.Split(strings.ReplaceAll(input, "\r\n", "\n"), "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		b, err := parseBrick(line)
		if err != nil {
			return nil, err
		}
		bricks = append(bricks, b)
	}
	if len(bricks) == 0 {
		return nil, fmt.Errorf("%w: empty snapshot", errInvalidSnapshot)
	}
	return bricks, nil
}

// settle drops every brick as far as it goes, lowest first, and records the
// support relations of the resulting stack.
func settle(input string) (stack, error) {
	bricks, err := parseBricks(input)
	if err != nil {
		return stack{}, err
	}
	slices.SortFunc(bricks, func(a, b brick) int { return cmp.Compare(a.min[2], b.min[2]) })

	// Per column: the highest occupied z and the brick owning that cube.
	type top struct {
		z     int
		brick int
	}
	tops := make(map[column]top)
	settled := stack{
		supports:    make([][]int, len(bricks)),
		supportedBy: make([][]int, len(bricks)),
	}
	for index, b := range bricks {
		columns := b.footprint()
		resting := 0
		for _, c := range columns {
			if occupied, ok := tops[c]; ok && occupied.z > resting {
				resting = occupied.z
			}
		}
		var supporters []int
		for _, c := range columns {
			if occupied, ok := tops[c]; ok && occupied.z == resting {
				supporters = append(supporters, occupied.brick)
			}
		}
		slices.Sort(supporters)
		supporters = slices.Compact(supporters)
		for _, supporter := range supporters {
			settled.supports[supporter] = append(settled.supports[supporter], index)
		}
		settled.supportedBy[index] = supporters
		for _, c := range columns {
			tops[c] = top{z: resting + b.height(), brick: index}
		}
	}
	return settled, nil
}

func (s stack) safeToDisintegrate() int {
	safe := 0
	for _, resting := range s.supports {
		if !slices.ContainsFunc(resting, func(above int) bool { return len(s.supportedBy[above]) == 1 }) {
			safe++
		}
	}
	return safe
}

func (s stack) totalChainReaction() int {
	total := 0
	for index := range s.supports {
		total += s.chainReaction(index)
	}
	return total
}

// chainReaction walks up the support graph: a brick falls once all of its
// supporters do.
func (s stack) chainReaction(removed int) int {
	falling := make([]bool, len(s.supports))
	falling[removed] = true
	queue := []int{removed}
	fallen := 0
	for len(queue) > 0 {
		index := queue[0]
		queue = queue[1:]
		for _, above := range s.supports[index] {
			if falling[above] || !allFalling(s.supportedBy[above], falling) {
				continue
			}
			falling[above] = true
			fallen++
			queue = append(queue, above)
		}
	}
	return fallen
}

func allFalling(supporters []int, falling []bool) bool {
	for _, supporter := range supporters {
		if !falling[supporter] {
			return false
		}
	}
	return true
}

func solve(input string) (int, int, error) {
	settled, err := settle(input)
	if err != nil {
		return 0, 0, err
	}
	return settled.safeToDisintegrate(), settled.totalChainReaction(), nil
}
