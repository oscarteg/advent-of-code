package main

import (
	"fmt"
	"strconv"
	"strings"
)

// pickKDigits greedily picks k digits from a bank string to form the largest k-digit number.
// For each position, picks the largest digit such that enough digits remain.
func pickKDigits(bank string, k int) (uint64, error) {
	if k > len(bank) {
		return 0, fmt.Errorf("not enough digits: need %d, have %d", k, len(bank))
	}

	result := make([]byte, 0, k)
	startIdx := 0

	for pos := 0; pos < k; pos++ {
		remainingNeeded := k - pos - 1
		endIdx := len(bank) - remainingNeeded - 1

		var bestDigit byte
		var bestIdx int = startIdx

		for i := startIdx; i <= endIdx; i++ {
			digit := bank[i]
			if digit < '0' || digit > '9' {
				return 0, fmt.Errorf("invalid character: %c", digit)
			}
			if i == startIdx || digit > bestDigit {
				bestDigit = digit
				bestIdx = i
			}
		}

		result = append(result, bestDigit)
		startIdx = bestIdx + 1
	}

	numStr := string(result)
	num, err := strconv.ParseUint(numStr, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("failed to parse number %s: %w", numStr, err)
	}

	return num, nil
}

func solve(input string) (uint64, uint64, error) {
	lines := strings.Split(strings.TrimSpace(input), "\n")
	var part1Sum, part2Sum uint64

	for _, line := range lines {
		if line == "" {
			continue
		}

		for _, c := range line {
			if c < '0' || c > '9' {
				return 0, 0, fmt.Errorf("invalid character in line: %c", c)
			}
		}

		p1, err := pickKDigits(line, 2)
		if err != nil {
			return 0, 0, err
		}
		part1Sum += p1

		p2, err := pickKDigits(line, 12)
		if err != nil {
			return 0, 0, err
		}
		part2Sum += p2
	}

	return part1Sum, part2Sum, nil
}
