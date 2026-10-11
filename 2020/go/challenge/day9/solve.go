package day9

import (
	"errors"
	"fmt"
	"slices"
	"strconv"

	"github.com/oscarteg/advent-of-code-2020/challenge"
)

const RealPreamble = 25

var errNotFound = errors.New("no answer in input")

func parseNumbers(input *challenge.Input) ([]int64, error) {
	lines := input.LinesSlice()
	numbers := make([]int64, len(lines))
	for i, line := range lines {
		n, err := strconv.ParseInt(line, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("line %d: %w", i+1, err)
		}
		numbers[i] = n
	}
	return numbers, nil
}

func isSumOfPair(number int64, window []int64) bool {
	for i := range window {
		for j := i + 1; j < len(window); j++ {
			if window[i]+window[j] == number {
				return true
			}
		}
	}
	return false
}

func findInvalid(numbers []int64, preamble int) (int64, error) {
	if preamble < 2 || len(numbers) <= preamble {
		return 0, fmt.Errorf("need a preamble of at least 2 and more than %d numbers, got %d", preamble, len(numbers))
	}
	for i := preamble; i < len(numbers); i++ {
		if !isSumOfPair(numbers[i], numbers[i-preamble:i]) {
			return numbers[i], nil
		}
	}
	return 0, fmt.Errorf("every number is a sum of its preamble: %w", errNotFound)
}

func findWeakness(numbers []int64, target int64) (int64, error) {
	for start := range numbers {
		sum := numbers[start]
		for end := start + 1; end < len(numbers) && sum < target; end++ {
			sum += numbers[end]
			if sum == target {
				run := numbers[start : end+1]
				return slices.Min(run) + slices.Max(run), nil
			}
		}
	}
	return 0, fmt.Errorf("no contiguous run sums to %d: %w", target, errNotFound)
}
