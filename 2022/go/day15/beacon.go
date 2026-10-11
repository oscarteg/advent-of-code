package main

import (
	"cmp"
	"errors"
	"fmt"
	"slices"
	"strings"
)

const (
	realRow          = 2_000_000
	realSearchLimit  = 4_000_000
	tuningMultiplier = 4_000_000
)

var (
	errNoReadings          = errors.New("input contains no sensor readings")
	errNoUncoveredPosition = errors.New("no position in the search area is free of sensor coverage")
)

// reading is a sensor, its closest beacon and the Manhattan distance between the two.
type reading struct {
	sensorX, sensorY int
	beaconX, beaconY int
	radius           int
}

// span is an inclusive range of x positions.
type span struct {
	start, end int
}

func parseReading(line string) (reading, error) {
	var r reading
	count, err := fmt.Sscanf(line, "Sensor at x=%d, y=%d: closest beacon is at x=%d, y=%d",
		&r.sensorX, &r.sensorY, &r.beaconX, &r.beaconY)
	if err != nil {
		return reading{}, fmt.Errorf("malformed reading %q: %w", line, err)
	}
	if count != 4 {
		return reading{}, fmt.Errorf("malformed reading %q: read %d of 4 coordinates", line, count)
	}
	r.radius = abs(r.sensorX-r.beaconX) + abs(r.sensorY-r.beaconY)
	return r, nil
}

func parse(input string) ([]reading, error) {
	var readings []reading
	for _, line := range strings.Split(input, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		r, err := parseReading(line)
		if err != nil {
			return nil, err
		}
		readings = append(readings, r)
	}
	if len(readings) == 0 {
		return nil, errNoReadings
	}
	return readings, nil
}

// rowSpan returns the inclusive x-range the sensor rules out on row, if it reaches it.
func (r reading) rowSpan(row int) (span, bool) {
	spread := r.radius - abs(r.sensorY-row)
	if spread < 0 {
		return span{}, false
	}
	return span{start: r.sensorX - spread, end: r.sensorX + spread}, true
}

// mergedRowSpans fills scratch with the disjoint x-ranges covering row, sorted by
// start and with touching ranges joined, so that any boundary left between them is
// a real gap. Scratch must have room for one span per reading.
func mergedRowSpans(scratch []span, readings []reading, row int) []span {
	spans := scratch[:0]
	for _, r := range readings {
		if s, ok := r.rowSpan(row); ok {
			spans = append(spans, s)
		}
	}
	if len(spans) == 0 {
		return spans
	}
	slices.SortFunc(spans, func(a, b span) int { return cmp.Compare(a.start, b.start) })

	merged := spans[:1]
	for _, s := range spans[1:] {
		last := &merged[len(merged)-1]
		if s.start <= last.end+1 {
			last.end = max(last.end, s.end)
			continue
		}
		merged = append(merged, s)
	}
	return merged
}

// excludedPositions returns the positions on row that cannot hold a beacon: every
// covered position minus the distinct known beacons there, which each sit inside
// their own sensor's range.
func excludedPositions(scratch []span, readings []reading, row int) int {
	covered := 0
	for _, s := range mergedRowSpans(scratch, readings, row) {
		covered += s.end - s.start + 1
	}
	beacons := make(map[int]struct{})
	for _, r := range readings {
		if r.beaconY == row {
			beacons[r.beaconX] = struct{}{}
		}
	}
	return covered - len(beacons)
}

// uncoveredX returns the smallest x in [0, limit] on row that no sensor covers.
func uncoveredX(scratch []span, readings []reading, row, limit int) (int, bool) {
	x := 0
	for _, s := range mergedRowSpans(scratch, readings, row) {
		if s.start > x {
			return x, true
		}
		x = max(x, s.end+1)
		if x > limit {
			return 0, false
		}
	}
	return x, x <= limit
}

// tuningFrequency returns the tuning frequency of the one position in the search
// square that no sensor covers.
func tuningFrequency(scratch []span, readings []reading, limit int) (int, error) {
	for row := 0; row <= limit; row++ {
		if x, ok := uncoveredX(scratch, readings, row, limit); ok {
			return x*tuningMultiplier + row, nil
		}
	}
	return 0, errNoUncoveredPosition
}

// solveWith counts the excluded positions on row and locates the distress beacon
// inside the square [0, searchLimit].
func solveWith(input string, row, searchLimit int) (part1, part2 int, err error) {
	readings, err := parse(input)
	if err != nil {
		return 0, 0, err
	}
	scratch := make([]span, len(readings))
	part2, err = tuningFrequency(scratch, readings, searchLimit)
	if err != nil {
		return 0, 0, err
	}
	return excludedPositions(scratch, readings, row), part2, nil
}

func solve(input string) (part1, part2 int, err error) {
	return solveWith(input, realRow, realSearchLimit)
}

func abs(n int) int {
	if n < 0 {
		return -n
	}
	return n
}
