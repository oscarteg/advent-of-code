package day5

import (
	"fmt"
	"sort"

	"github.com/oscarteg/advent-of-code-2020/challenge"
)

func B(input *challenge.Input) (int, error) {
	lines := input.LinesSlice()
	var seats []int

	for _, line := range lines {
		if line == "" {
			continue
		}

		seat, err := parseSeat(line)
		if err != nil {
			return 0, err
		}
		seats = append(seats, seat)
	}

	sort.Ints(seats)

	for i := 0; i < len(seats)-1; i++ {
		if seats[i+1]-seats[i] == 2 {
			return seats[i] + 1, nil
		}
	}

	return 0, fmt.Errorf("missing seat not found")
}
