package day19

import (
	"github.com/oscarteg/advent-of-code-2020/challenge"
)

func B(input *challenge.Input) (int, error) {
	rules, messages, err := parse(input)
	if err != nil {
		return 0, err
	}
	looped, err := rules.withLoops()
	if err != nil {
		return 0, err
	}

	return looped.countMatching(messages), nil
}
