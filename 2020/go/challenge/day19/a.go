package day19

import (
	"github.com/oscarteg/advent-of-code-2020/challenge"
)

func A(input *challenge.Input) (int, error) {
	rules, messages, err := parse(input)
	if err != nil {
		return 0, err
	}

	return rules.countMatching(messages), nil
}
