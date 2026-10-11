package day15

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/oscarteg/advent-of-code-2020/challenge"
)

func B(input *challenge.Input) (int, error) {
	starting, err := parseStartingNumbers(input)
	if err != nil {
		return 0, err
	}
	return play(starting, 30_000_000)
}

func parseStartingNumbers(input *challenge.Input) ([]int, error) {
	line := strings.TrimSpace(strings.Join(input.LinesSlice(), ""))
	if line == "" {
		return nil, fmt.Errorf("missing starting numbers")
	}
	fields := strings.Split(line, ",")
	numbers := make([]int, len(fields))
	for i, field := range fields {
		value, err := strconv.Atoi(strings.TrimSpace(field))
		if err != nil {
			return nil, fmt.Errorf("invalid starting number %q: %w", field, err)
		}
		numbers[i] = value
	}
	return numbers, nil
}

// play simulates the memory game up to turn using a flat slice of last-seen
// turns per number, so each step is O(1) instead of scanning history.
func play(starting []int, turn int) (int, error) {
	if turn < len(starting) {
		return 0, fmt.Errorf("turn must reach past the starting numbers")
	}
	maxStarting := 0
	for _, number := range starting {
		if number > maxStarting {
			maxStarting = number
		}
	}
	lastSeenSize := turn
	if maxStarting+1 > lastSeenSize {
		lastSeenSize = maxStarting + 1
	}
	lastSeen := make([]int, lastSeenSize)
	for index, number := range starting[:len(starting)-1] {
		lastSeen[number] = index + 1
	}
	last := starting[len(starting)-1]
	for currentTurn := len(starting); currentTurn < turn; currentTurn++ {
		seenAt := lastSeen[last]
		next := 0
		if seenAt != 0 {
			next = currentTurn - seenAt
		}
		lastSeen[last] = currentTurn
		last = next
	}
	return last, nil
}
