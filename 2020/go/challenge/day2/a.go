package day2

import (
	"github.com/oscarteg/advent-of-code-2020/challenge"
)

func A(input *challenge.Input) (int, error) {
	policies, err := parsePolicies(input)
	if err != nil {
		return 0, err
	}

	var valid int
	for _, p := range policies {
		if p.isValidByCount() {
			valid++
		}
	}

	return valid, nil
}
