package day6

import (
	"github.com/oscarteg/advent-of-code-2020/challenge"
)

func A(input *challenge.Input) (int, error) {
	groups, err := parseGroups(input.LinesSlice())
	if err != nil {
		return 0, err
	}

	total := 0
	for _, group := range groups {
		total += len(group.anyAnswered())
	}

	return total, nil
}
