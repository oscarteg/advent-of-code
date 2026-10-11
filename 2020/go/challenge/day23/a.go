package day23

import (
	"github.com/oscarteg/advent-of-code-2020/challenge"
)

func A(input *challenge.Input) (string, error) {
	cups, err := parseCups(input)
	if err != nil {
		return "", err
	}
	next, err := play(cups, len(cups), part1Moves)
	if err != nil {
		return "", err
	}
	return labelsAfterOne(next), nil
}
