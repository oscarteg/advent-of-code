package day8

import (
	"github.com/oscarteg/advent-of-code-2020/challenge"
)

func B(input *challenge.Input) (int, error) {
	program, err := parseProgram(input.LinesSlice())
	if err != nil {
		return 0, err
	}
	return fixAndRun(program)
}
