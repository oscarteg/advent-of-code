package main

import "fmt"

func solve(input string) (int, int, error) {
	part1, err := findMarker(input, 4)
	if err != nil {
		return 0, 0, fmt.Errorf("part 1: %w", err)
	}
	part2, err := findMarker(input, 14)
	if err != nil {
		return 0, 0, fmt.Errorf("part 2: %w", err)
	}
	return part1, part2, nil
}

func findMarker(signal string, windowSize int) (int, error) {
	if len(signal) < windowSize {
		return 0, fmt.Errorf("signal length %d is less than window size %d", len(signal), windowSize)
	}

	for i := windowSize; i <= len(signal); i++ {
		window := signal[i-windowSize : i]
		if hasUniqueChars(window) {
			return i, nil
		}
	}

	return 0, fmt.Errorf("no marker found")
}

func hasUniqueChars(s string) bool {
	seen := make(map[rune]bool)
	for _, ch := range s {
		if seen[ch] {
			return false
		}
		seen[ch] = true
	}
	return true
}
