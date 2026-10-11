package main

import (
	"fmt"
	"strconv"
	"strings"
)

func solve(input string) (part1, part2 int, err error) {
	lines := strings.Split(strings.TrimSpace(input), "\n")
	if len(lines) == 0 {
		return 0, 0, fmt.Errorf("input is empty")
	}

	for i, line := range lines {
		if len(line) == 0 {
			return 0, 0, fmt.Errorf("line %d is empty", i)
		}
		if len(line) != len(lines[0]) {
			return 0, 0, fmt.Errorf("line %d has length %d, expected %d", i, len(line), len(lines[0]))
		}
		for j, ch := range line {
			if ch != '0' && ch != '1' {
				return 0, 0, fmt.Errorf("line %d char %d: invalid bit %q", i, j, ch)
			}
		}
	}

	bits := len(lines[0])

	part1, err = powerConsumption(lines, bits)
	if err != nil {
		return 0, 0, err
	}

	part2, err = lifeSupportRating(lines, bits)
	if err != nil {
		return 0, 0, err
	}

	return part1, part2, nil
}

func powerConsumption(lines []string, bits int) (int, error) {
	max := (len(lines) + 1) / 2
	gamma := 0

	for i := 0; i < bits; i++ {
		ones := 0
		for _, line := range lines {
			if line[i] == '1' {
				ones++
			}
		}
		if ones >= max {
			gamma |= 1 << (bits - i - 1)
		}
	}

	bitmask := (1 << bits) - 1
	epsilon := ^gamma & bitmask

	return gamma * epsilon, nil
}

func lifeSupportRating(readings []string, bits int) (int, error) {
	oxygenStr, err := findRating(readings, bits, true)
	if err != nil {
		return 0, err
	}

	co2Str, err := findRating(readings, bits, false)
	if err != nil {
		return 0, err
	}

	oxygen, err := strconv.ParseInt(oxygenStr, 2, 32)
	if err != nil {
		return 0, fmt.Errorf("failed to parse oxygen rating: %w", err)
	}

	co2, err := strconv.ParseInt(co2Str, 2, 32)
	if err != nil {
		return 0, fmt.Errorf("failed to parse CO2 rating: %w", err)
	}

	return int(oxygen) * int(co2), nil
}

func findRating(readings []string, bits int, keepMostCommon bool) (string, error) {
	candidates := make([]string, len(readings))
	copy(candidates, readings)

	for i := 0; i < bits; i++ {
		if len(candidates) == 1 {
			break
		}

		ones := 0
		for _, s := range candidates {
			if s[i] == '1' {
				ones++
			}
		}
		zeros := len(candidates) - ones

		var desiredBit byte
		if (ones >= zeros) == keepMostCommon {
			desiredBit = '1'
		} else {
			desiredBit = '0'
		}

		filtered := make([]string, 0, len(candidates))
		for _, s := range candidates {
			if s[i] == desiredBit {
				filtered = append(filtered, s)
			}
		}
		candidates = filtered
	}

	if len(candidates) == 0 {
		return "", fmt.Errorf("no candidates found")
	}

	return candidates[0], nil
}
