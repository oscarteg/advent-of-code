package day8

import (
	"github.com/oscarteg/advent-of-code-2020/challenge"
)

func A(input *challenge.Input) (int, error) {
	program, err := parseProgram(input.LinesSlice())
	if err != nil {
		return 0, err
	}
	_, acc := run(program)
	return acc, nil
}
