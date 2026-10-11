package day23

import (
	"github.com/oscarteg/advent-of-code-2020/challenge"
)

func B(input *challenge.Input) (int, error) {
	cups, err := parseCups(input)
	if err != nil {
		return 0, err
	}
	next, err := play(cups, part2Cups, part2Moves)
	if err != nil {
		return 0, err
	}
	first := next[1]
	return first * next[first], nil
}
