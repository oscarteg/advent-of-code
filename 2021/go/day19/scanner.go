package main

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// minOverlap is the number of beacons two scanners must share before their
// reports can be merged.
const minOverlap = 12

var errInvalidReport = errors.New("invalid scanner report")

// point is a beacon or scanner position in three dimensions.
type point [3]int

// report holds the beacons one scanner sees, in that scanner's own frame.
type report struct {
	beacons []point
}

// rotation is a signed axis permutation with determinant +1: one of the 24
// orientations a scanner can have.
type rotation struct {
	axes  [3]int
	signs [3]int
}

func (r rotation) apply(p point) point {
	return point{
		r.signs[0] * p[r.axes[0]],
		r.signs[1] * p[r.axes[1]],
		r.signs[2] * p[r.axes[2]],
	}
}

// allRotations combines the six axis permutations with the sign flips that keep
// the determinant positive, which drops the 24 mirror images.
func allRotations() []rotation {
	permutations := []struct {
		axes   [3]int
		parity int
	}{
		{[3]int{0, 1, 2}, 1},
		{[3]int{1, 2, 0}, 1},
		{[3]int{2, 0, 1}, 1},
		{[3]int{0, 2, 1}, -1},
		{[3]int{2, 1, 0}, -1},
		{[3]int{1, 0, 2}, -1},
	}
	rotations := make([]rotation, 0, 24)
	for _, permutation := range permutations {
		for flips := 0; flips < 8; flips++ {
			signs := [3]int{1, 1, 1}
			for axis := range signs {
				if flips&(1<<axis) != 0 {
					signs[axis] = -1
				}
			}
			if permutation.parity*signs[0]*signs[1]*signs[2] == 1 {
				rotations = append(rotations, rotation{axes: permutation.axes, signs: signs})
			}
		}
	}
	return rotations
}

func scannerHeader(line string) (string, bool) {
	rest, ok := strings.CutPrefix(line, "--- scanner ")
	if !ok {
		return "", false
	}
	return strings.CutSuffix(rest, " ---")
}

func parsePoint(line string) (point, error) {
	fields := strings.Split(line, ",")
	if len(fields) != 3 {
		return point{}, fmt.Errorf("%w: beacon %q needs three coordinates", errInvalidReport, line)
	}
	var parsed point
	for axis, field := range fields {
		value, err := strconv.Atoi(strings.TrimSpace(field))
		if err != nil {
			return point{}, fmt.Errorf("beacon %q: %w", line, err)
		}
		parsed[axis] = value
	}
	return parsed, nil
}

func parseReports(input string) ([]report, error) {
	var reports []report
	for _, line := range strings.Split(strings.ReplaceAll(input, "\r\n", "\n"), "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		if id, isHeader := scannerHeader(line); isHeader {
			if _, err := strconv.Atoi(id); err != nil {
				return nil, fmt.Errorf("scanner id %q: %w", id, err)
			}
			reports = append(reports, report{})
			continue
		}
		if len(reports) == 0 {
			return nil, fmt.Errorf("%w: beacon %q appears before any scanner header", errInvalidReport, line)
		}
		beacon, err := parsePoint(line)
		if err != nil {
			return nil, err
		}
		current := &reports[len(reports)-1]
		current.beacons = append(current.beacons, beacon)
	}
	return reports, nil
}

// align searches for the orientation and offset that move a report into the frame
// of the already placed beacons, returning the scanner position and its beacons.
func align(placed map[point]bool, r report, rotations []rotation) (point, []point, bool) {
	for _, rot := range rotations {
		rotated := make([]point, len(r.beacons))
		for i, beacon := range r.beacons {
			rotated[i] = rot.apply(beacon)
		}
		offsets := make(map[point]int)
		for known := range placed {
			for _, beacon := range rotated {
				offset := point{known[0] - beacon[0], known[1] - beacon[1], known[2] - beacon[2]}
				offsets[offset]++
				if offsets[offset] < minOverlap {
					continue
				}
				aligned := make([]point, len(rotated))
				for i, beacon := range rotated {
					aligned[i] = point{beacon[0] + offset[0], beacon[1] + offset[1], beacon[2] + offset[2]}
				}
				return offset, aligned, true
			}
		}
	}
	return point{}, nil, false
}

// assemble places every report in scanner 0's frame, repeatedly matching the
// reports that are still unplaced against the beacons of a freshly placed one.
func assemble(reports []report) (map[point]bool, []point, error) {
	if len(reports) == 0 {
		return nil, nil, fmt.Errorf("%w: input contains no scanner reports", errInvalidReport)
	}
	rotations := allRotations()
	beacons := make(map[point]bool)
	for _, beacon := range reports[0].beacons {
		beacons[beacon] = true
	}
	scanners := []point{{0, 0, 0}}
	frontier := [][]point{reports[0].beacons}
	pending := reports[1:]

	for len(frontier) > 0 {
		last := len(frontier) - 1
		placed := make(map[point]bool, len(frontier[last]))
		for _, beacon := range frontier[last] {
			placed[beacon] = true
		}
		frontier = frontier[:last]

		unplaced := make([]report, 0, len(pending))
		for _, r := range pending {
			scanner, aligned, ok := align(placed, r, rotations)
			if !ok {
				unplaced = append(unplaced, r)
				continue
			}
			for _, beacon := range aligned {
				beacons[beacon] = true
			}
			scanners = append(scanners, scanner)
			frontier = append(frontier, aligned)
		}
		pending = unplaced
	}

	if len(pending) > 0 {
		return nil, nil, fmt.Errorf("%w: %d reports share fewer than %d beacons with the rest", errInvalidReport, len(pending), minOverlap)
	}
	return beacons, scanners, nil
}

func manhattan(from, to point) int {
	total := 0
	for axis := range from {
		delta := from[axis] - to[axis]
		if delta < 0 {
			delta = -delta
		}
		total += delta
	}
	return total
}

func largestDistance(scanners []point) int {
	largest := 0
	for i, from := range scanners {
		for _, to := range scanners[i+1:] {
			if distance := manhattan(from, to); distance > largest {
				largest = distance
			}
		}
	}
	return largest
}

func solve(input string) (int, int, error) {
	reports, err := parseReports(input)
	if err != nil {
		return 0, 0, err
	}
	beacons, scanners, err := assemble(reports)
	if err != nil {
		return 0, 0, err
	}
	return len(beacons), largestDistance(scanners), nil
}
