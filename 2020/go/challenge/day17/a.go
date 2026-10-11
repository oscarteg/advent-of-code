package day17

import (
	"github.com/oscarteg/advent-of-code-2020/challenge"
)

func A(input *challenge.Input) (int, error) {
	active, err := parseInput3D(input)
	if err != nil {
		return 0, err
	}
	for i := 0; i < 6; i++ {
		active = step3D(active)
	}
	return len(active), nil
}
