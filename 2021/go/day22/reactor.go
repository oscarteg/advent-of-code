package main

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

const initializationLimit = 50

var errInvalidStep = errors.New("invalid reboot step")

type interval struct {
	start, end int
}

func parseInterval(field, axis string) (interval, error) {
	rest, ok := strings.CutPrefix(strings.TrimSpace(field), axis+"=")
	if !ok {
		return interval{}, fmt.Errorf("%w: expected %s=start..end, got %q", errInvalidStep, axis, field)
	}
	start, end, ok := strings.Cut(rest, "..")
	if !ok {
		return interval{}, fmt.Errorf("%w: range %q must contain ..", errInvalidStep, rest)
	}
	lower, err := strconv.Atoi(start)
	if err != nil {
		return interval{}, fmt.Errorf("range start: %w", err)
	}
	upper, err := strconv.Atoi(end)
	if err != nil {
		return interval{}, fmt.Errorf("range end: %w", err)
	}
	if lower > upper {
		return interval{}, fmt.Errorf("%w: range %d..%d is empty", errInvalidStep, lower, upper)
	}
	return interval{start: lower, end: upper}, nil
}

func (i interval) intersect(other interval) (interval, bool) {
	start, end := max(i.start, other.start), min(i.end, other.end)
	return interval{start: start, end: end}, start <= end
}

func (i interval) length() int {
	return i.end - i.start + 1
}

type cuboid struct {
	x, y, z interval
}

func parseCuboid(ranges string) (cuboid, error) {
	fields := strings.Split(ranges, ",")
	if len(fields) != 3 {
		return cuboid{}, fmt.Errorf("%w: expected x, y and z ranges, got %q", errInvalidStep, ranges)
	}
	x, err := parseInterval(fields[0], "x")
	if err != nil {
		return cuboid{}, err
	}
	y, err := parseInterval(fields[1], "y")
	if err != nil {
		return cuboid{}, err
	}
	z, err := parseInterval(fields[2], "z")
	if err != nil {
		return cuboid{}, err
	}
	return cuboid{x: x, y: y, z: z}, nil
}

func (c cuboid) intersect(other cuboid) (cuboid, bool) {
	x, ok := c.x.intersect(other.x)
	if !ok {
		return cuboid{}, false
	}
	y, ok := c.y.intersect(other.y)
	if !ok {
		return cuboid{}, false
	}
	z, ok := c.z.intersect(other.z)
	if !ok {
		return cuboid{}, false
	}
	return cuboid{x: x, y: y, z: z}, true
}

func (c cuboid) volume() int {
	return c.x.length() * c.y.length() * c.z.length()
}

type step struct {
	on     bool
	cuboid cuboid
}

func parseStep(line string) (step, error) {
	state, ranges, ok := strings.Cut(strings.TrimSpace(line), " ")
	if !ok {
		return step{}, fmt.Errorf("%w: %q must look like on x=0..1,y=0..1,z=0..1", errInvalidStep, line)
	}
	var on bool
	switch state {
	case "on":
		on = true
	case "off":
		on = false
	default:
		return step{}, fmt.Errorf("%w: unknown state %q", errInvalidStep, state)
	}
	parsed, err := parseCuboid(ranges)
	if err != nil {
		return step{}, err
	}
	return step{on: on, cuboid: parsed}, nil
}

func parseSteps(input string) ([]step, error) {
	var steps []step
	for _, line := range strings.Split(strings.ReplaceAll(input, "\r\n", "\n"), "\n") {
		if strings.TrimSpace(line) == "" {
			continue
		}
		parsed, err := parseStep(line)
		if err != nil {
			return nil, err
		}
		steps = append(steps, parsed)
	}
	return steps, nil
}

type signedCuboid struct {
	cuboid cuboid
	sign   int
}

// Each counted cuboid carries a sign: a new cuboid cancels its overlap with
// every already counted cuboid by adding that overlap with the opposite sign,
// so the signed volumes always sum to the number of cubes that are on.
// Steps are clipped to region; a nil region counts the whole reactor.
func countOn(steps []step, region *cuboid) int {
	var counted []signedCuboid
	for _, current := range steps {
		box := current.cuboid
		if region != nil {
			clipped, ok := region.intersect(box)
			if !ok {
				continue
			}
			box = clipped
		}
		corrections := make([]signedCuboid, 0, len(counted))
		for _, previous := range counted {
			if overlap, ok := previous.cuboid.intersect(box); ok {
				corrections = append(corrections, signedCuboid{cuboid: overlap, sign: -previous.sign})
			}
		}
		counted = append(counted, corrections...)
		if current.on {
			counted = append(counted, signedCuboid{cuboid: box, sign: 1})
		}
	}
	total := 0
	for _, entry := range counted {
		total += entry.sign * entry.cuboid.volume()
	}
	return total
}

func solve(input string) (part1, part2 int, err error) {
	steps, err := parseSteps(input)
	if err != nil {
		return 0, 0, err
	}
	limit := interval{start: -initializationLimit, end: initializationLimit}
	initialization := cuboid{x: limit, y: limit, z: limit}
	return countOn(steps, &initialization), countOn(steps, nil), nil
}
