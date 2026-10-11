package day1

import (
	"fmt"
	"strconv"

	"github.com/oscarteg/advent-of-code-2020/challenge"
)

func A(input *challenge.Input) (int, error) {
	entries := make(map[int]bool)

	for line := range input.Lines() {
		num, err := strconv.Atoi(line)
		if err != nil {
			return 0, err
		}
		entries[num] = true
	}

	for x := range entries {
		y := 2020 - x
		if entries[y] && x != y {
			return x * y, nil
		}
	}

	return 0, fmt.Errorf("no solution found for part a")
}
