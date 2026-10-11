package day15

import (
	"github.com/oscarteg/advent-of-code-2020/challenge"
)

func A(input *challenge.Input) (int, error) {
	starting, err := parseStartingNumbers(input)
	if err != nil {
		return 0, err
	}
	return play(starting, 2020)
}
