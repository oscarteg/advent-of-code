// Package main solves Advent of Code 2022 Day 20: Grove Positioning System.
package main

import (
	"fmt"
	"strconv"
	"strings"
)

const (
	decryptionKey  = 811589153
	mixRoundsPart2 = 10
)

func parseNumbers(input string) ([]int, error) {
	var numbers []int
	for _, line := range strings.Split(input, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		value, err := strconv.Atoi(line)
		if err != nil {
			return nil, fmt.Errorf("invalid number %q: %w", line, err)
		}
		numbers = append(numbers, value)
	}
	if len(numbers) == 0 {
		return nil, fmt.Errorf("input must contain at least one number")
	}
	return numbers, nil
}

// mix performs rounds passes over numbers, moving each value (by its
// original index, since values can repeat) forward or backward in the
// circular list by its own amount, and returns the resulting order.
func mix(numbers []int, rounds int) []int {
	length := len(numbers)
	order := make([]int, length)
	for i := range order {
		order[i] = i
	}
	if length <= 1 {
		mixed := make([]int, length)
		copy(mixed, numbers)
		return mixed
	}

	indexOf := func(originalIndex int) int {
		for pos, idx := range order {
			if idx == originalIndex {
				return pos
			}
		}
		panic("original index must be present in order")
	}

	for r := 0; r < rounds; r++ {
		for originalIndex, value := range numbers {
			currentPos := indexOf(originalIndex)
			order = append(order[:currentPos], order[currentPos+1:]...)
			newPos := ((currentPos+value)%(length-1) + (length - 1)) % (length - 1)
			order = append(order[:newPos], append([]int{originalIndex}, order[newPos:]...)...)
		}
	}

	mixed := make([]int, length)
	for i, originalIndex := range order {
		mixed[i] = numbers[originalIndex]
	}
	return mixed
}

func groveCoordinatesSum(mixed []int) (int, error) {
	zeroPos := -1
	for i, v := range mixed {
		if v == 0 {
			zeroPos = i
			break
		}
	}
	if zeroPos == -1 {
		return 0, fmt.Errorf("mixed list must contain a 0 value")
	}
	length := len(mixed)
	sum := 0
	for _, offset := range [3]int{1000, 2000, 3000} {
		sum += mixed[(zeroPos+offset)%length]
	}
	return sum, nil
}

func solve(input string) (int, int, error) {
	numbers, err := parseNumbers(input)
	if err != nil {
		return 0, 0, err
	}

	mixedOnce := mix(numbers, 1)
	part1, err := groveCoordinatesSum(mixedOnce)
	if err != nil {
		return 0, 0, err
	}

	decrypted := make([]int, len(numbers))
	for i, n := range numbers {
		decrypted[i] = n * decryptionKey
	}
	mixedMany := mix(decrypted, mixRoundsPart2)
	part2, err := groveCoordinatesSum(mixedMany)
	if err != nil {
		return 0, 0, err
	}

	return part1, part2, nil
}
