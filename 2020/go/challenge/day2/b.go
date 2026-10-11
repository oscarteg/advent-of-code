package day2

import (
	"github.com/oscarteg/advent-of-code-2020/challenge"
)

func B(input *challenge.Input) (int, error) {
	policies, err := parsePolicies(input)
	if err != nil {
		return 0, err
	}

	var valid int
	for _, p := range policies {
		ok, err := p.isValidByPosition()
		if err != nil {
			return 0, err
		}
		if ok {
			valid++
		}
	}

	return valid, nil
}
