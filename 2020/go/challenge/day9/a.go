package day9

import (
	"github.com/oscarteg/advent-of-code-2020/challenge"
)

func A(input *challenge.Input, preamble int) (int64, error) {
	numbers, err := parseNumbers(input)
	if err != nil {
		return 0, err
	}
	return findInvalid(numbers, preamble)
}
