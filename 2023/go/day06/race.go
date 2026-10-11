package main

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

var errInvalidRaceSheet = errors.New("invalid race sheet")

type race struct {
	time     uint64
	distance uint64
}

func (r race) winningHoldTimes() uint64 {
	var count uint64
	for hold := uint64(0); hold <= r.time; hold++ {
		if hold*(r.time-hold) > r.distance {
			count++
		}
	}
	return count
}

func splitLines(input string) (string, string, error) {
	lines := strings.Split(strings.TrimSpace(strings.ReplaceAll(input, "\r\n", "\n")), "\n")
	if len(lines) != 2 {
		return "", "", fmt.Errorf("%w: expected a time and a distance line", errInvalidRaceSheet)
	}
	return lines[0], lines[1], nil
}

func numbersAfter(line, prefix string) ([]uint64, error) {
	rest, ok := strings.CutPrefix(line, prefix)
	if !ok {
		return nil, fmt.Errorf("%w: missing prefix %q", errInvalidRaceSheet, prefix)
	}
	fields := strings.Fields(rest)
	numbers := make([]uint64, 0, len(fields))
	for _, field := range fields {
		number, err := strconv.ParseUint(field, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("number: %w", err)
		}
		numbers = append(numbers, number)
	}
	return numbers, nil
}

func parseRaces(input string) ([]race, error) {
	timeLine, distanceLine, err := splitLines(input)
	if err != nil {
		return nil, err
	}
	times, err := numbersAfter(timeLine, "Time:")
	if err != nil {
		return nil, err
	}
	distances, err := numbersAfter(distanceLine, "Distance:")
	if err != nil {
		return nil, err
	}
	if len(times) == 0 || len(times) != len(distances) {
		return nil, fmt.Errorf("%w: time and distance columns must match and be non-empty", errInvalidRaceSheet)
	}
	races := make([]race, len(times))
	for i := range times {
		races[i] = race{time: times[i], distance: distances[i]}
	}
	return races, nil
}

func joinedDigits(line, prefix string) (uint64, error) {
	rest, ok := strings.CutPrefix(line, prefix)
	if !ok {
		return 0, fmt.Errorf("%w: missing prefix %q", errInvalidRaceSheet, prefix)
	}
	digits := strings.ReplaceAll(rest, " ", "")
	if digits == "" {
		return 0, fmt.Errorf("%w: no digits on line", errInvalidRaceSheet)
	}
	number, err := strconv.ParseUint(digits, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("number: %w", err)
	}
	return number, nil
}

func parseSingleRace(input string) (race, error) {
	timeLine, distanceLine, err := splitLines(input)
	if err != nil {
		return race{}, err
	}
	time, err := joinedDigits(timeLine, "Time:")
	if err != nil {
		return race{}, err
	}
	distance, err := joinedDigits(distanceLine, "Distance:")
	if err != nil {
		return race{}, err
	}
	return race{time: time, distance: distance}, nil
}

func solve(input string) (uint64, uint64, error) {
	races, err := parseRaces(input)
	if err != nil {
		return 0, 0, err
	}
	part1 := uint64(1)
	for _, r := range races {
		part1 *= r.winningHoldTimes()
	}

	joined, err := parseSingleRace(input)
	if err != nil {
		return 0, 0, err
	}
	part2 := joined.winningHoldTimes()

	return part1, part2, nil
}
