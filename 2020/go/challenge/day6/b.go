package day6

import (
	"github.com/oscarteg/advent-of-code-2020/challenge"
)

func B(input *challenge.Input) (int, error) {
	groups, err := parseGroups(input.LinesSlice())
	if err != nil {
		return 0, err
	}

	total := 0
	for _, group := range groups {
		total += len(group.allAnswered())
	}

	return total, nil
}
