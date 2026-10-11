package day4

import (
	"github.com/oscarteg/advent-of-code-2020/challenge"
)

func B(input *challenge.Input) (int, error) {
	passports, err := parsePassports(input)
	if err != nil {
		return 0, err
	}

	count := 0
	for _, p := range passports {
		if p.isValid() {
			count++
		}
	}
	return count, nil
}
