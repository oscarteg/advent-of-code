package day1

import (
	"fmt"
	"strconv"

	"github.com/oscarteg/advent-of-code-2020/challenge"
)

func B(input *challenge.Input) (int, error) {
	entries := make(map[int]bool)
	var entryList []int

	for line := range input.Lines() {
		num, err := strconv.Atoi(line)
		if err != nil {
			return 0, err
		}
		entries[num] = true
		entryList = append(entryList, num)
	}

	for _, x := range entryList {
		for _, y := range entryList {
			sum := x + y
			if sum >= 2020 {
				continue
			}
			z := 2020 - sum
			if entries[z] && z != x && z != y {
				return x * y * z, nil
			}
		}
	}

	return 0, fmt.Errorf("no solution found for part b")
}
