package day16

import (
	"github.com/oscarteg/advent-of-code-2020/challenge"
)

func A(input *challenge.Input) (int, error) {
	parsed, err := parseNotes(input.LinesSlice())
	if err != nil {
		return 0, err
	}
	return parsed.errorRate(), nil
}
