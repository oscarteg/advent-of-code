package day23

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/oscarteg/advent-of-code-2020/challenge"
)

const (
	part1Moves = 100
	part2Cups  = 1_000_000
	part2Moves = 10_000_000
	// With fewer cups the three picked-up cups would overlap the current cup,
	// which the rules do not define.
	minCups = 5
)

// parseCups reads the single line of cup labels into a permutation of 1..n, the
// invariant the successor array relies on.
func parseCups(input *challenge.Input) ([]int, error) {
	labels := ""
	for _, line := range input.LinesSlice() {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		if labels != "" {
			return nil, errors.New("expected a single line of cup labels")
		}
		labels = line
	}
	if labels == "" {
		return nil, errors.New("expected a line of cup labels")
	}

	cups := make([]int, 0, len(labels))
	seen := make([]bool, len(labels)+1)
	for _, character := range labels {
		if character < '0' || character > '9' {
			return nil, fmt.Errorf("invalid cup label %q", character)
		}
		label := int(character - '0')
		if label < 1 || label > len(labels) || seen[label] {
			return nil, fmt.Errorf("cup labels must be a permutation of 1..%d, got %q", len(labels), labels)
		}
		seen[label] = true
		cups = append(cups, label)
	}
	if len(cups) < minCups {
		return nil, fmt.Errorf("need at least %d cups, got %d", minCups, len(cups))
	}
	return cups, nil
}

// play plays the given number of moves on a circle of total cups: the labels of
// cups clockwise, followed by len(cups)+1..total. It returns the successor
// array, in which next[label] is the label immediately clockwise of it.
func play(cups []int, total, moves int) ([]int, error) {
	if len(cups) == 0 {
		return nil, errors.New("circle without cups")
	}
	if total < len(cups) {
		return nil, fmt.Errorf("circle of %d cups cannot hold %d labels", total, len(cups))
	}

	// Every label starts out followed by the next one, which already wires the
	// padding cups; the input labels and the closing link are patched in below.
	next := make([]int, total+1)
	for label := range next {
		next[label] = label + 1
	}
	start := cups[0]
	for i, cup := range cups[:len(cups)-1] {
		next[cup] = cups[i+1]
	}
	if total > len(cups) {
		next[cups[len(cups)-1]] = len(cups) + 1
		next[total] = start
	} else {
		next[cups[len(cups)-1]] = start
	}

	current := start
	for range moves {
		first := next[current]
		second := next[first]
		third := next[second]
		next[current] = next[third]

		destination := counterClockwiseLabel(current, total)
		for destination == first || destination == second || destination == third {
			destination = counterClockwiseLabel(destination, total)
		}

		next[third] = next[destination]
		next[destination] = first
		current = next[current]
	}
	return next, nil
}

func counterClockwiseLabel(label, total int) int {
	if label == 1 {
		return total
	}
	return label - 1
}

func labelsAfterOne(next []int) string {
	var labels strings.Builder
	for cup := next[1]; cup != 1; cup = next[cup] {
		labels.WriteString(strconv.Itoa(cup))
	}
	return labels.String()
}
